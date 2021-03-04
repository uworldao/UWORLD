package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jessevdk/go-flags"
	"github.com/uworldao/UWORLD/common/hasharry"
	"github.com/uworldao/UWORLD/common/keystore"
	"github.com/uworldao/UWORLD/common/utils"
	"github.com/uworldao/UWORLD/core/types"
	common "github.com/uworldao/UWORLD/log"
	"github.com/uworldao/UWORLD/log/log15"
	"github.com/uworldao/UWORLD/param"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	DefaultHomeDir     = utils.AppDataDir(param.AppName, false)
	defaultP2pPort     = "33330"
	DefaultRpcPort     = "33331"
	defaultPrivateFile = "default_key.json"
	defaultKey         = "ub_chain"
	defaultExternalIp  = "0.0.0.0"
	DefaultFallBack    = int64(-1)
	defaultCoinHeight  = uint64(1)
)

// Config is the node startup parameter
type Config struct {
	ConfigFile  string `long:"config" description:"Start with a configuration file"`
	HomeDir     string `long:"appdata" description:"Path to application home directory"`
	DataDir     string `long:"data" description:"Path to application data directory"`
	FileLogging bool   `long:"filelogging" description:"Logging switch"`
	ExternalIp  string `long:"externalip" description:"External network IP address"`
	Bootstrap   string `long:"bootstrap" description:"Custom bootstrap"`
	P2pPort     string `long:"p2pport" description:"Add an interface/port to listen for connections"`
	RpcPort     string `long:"rpcport" description:"Add an interface/port to listen for RPC connections"`
	RpcTLS      bool   `long:"rpctls" description:"Open TLS for the RPC server -- NOTE: This is only allowed if the RPC server is bound to localhost"`
	RpcCert     string `long:"rpccert" description:"File containing the certificate file"`
	RpcKey      string `long:"rpckey" description:"File containing the certificate key"`
	RpcPass     string `long:"rpcpass" description:"Password for RPC connections"`
	TestNet     bool   `long:"testnet" description:"Use the test network"`
	KeyFile     string `long:"keyfile" description:"If you participate in mining, you need to configure the mining address key file"`
	KeyPass     string `long:"keypass" description:"The decryption password for key file"`
	FallBackTo  int64  `long:"fallbackto" description:"Force back to a height"`
	Version     bool   `long:"version" description:"View Version number"`
	NodePrivate *NodePrivate
}

// LoadConfig load the parse node startup parameter
func LoadConfig() (*Config, error) {
	cfg := &Config{
		HomeDir:    DefaultHomeDir,
		P2pPort:    defaultP2pPort,
		RpcPort:    DefaultRpcPort,
		FallBackTo: DefaultFallBack,
	}
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	preParser := newConfigParser(cfg, flags.HelpFlag)
	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type != flags.ErrHelp {
			return nil, err
		} else if ok && e.Type == flags.ErrHelp {
			return nil, err
		}
	}

	if cfg.ConfigFile != "" {
		_, err = toml.DecodeFile(cfg.ConfigFile, cfg)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Version {
		fmt.Printf("UWorld version %s\n", param.Version)
		os.Exit(0)
	}

	// Set the default external IP. If the external IP is not set,
	// other nodes can only know you but cannot send messages to you.
	if cfg.ExternalIp == "" {
		cfg.ExternalIp = defaultExternalIp
	}

	// Node data and file storage directory, if not set,
	// use the default directory
	if cfg.HomeDir == "" {
		cfg.HomeDir = DefaultHomeDir
	}

	// p2p service listening port, if not, use the default port
	if cfg.P2pPort == "" {
		cfg.P2pPort = defaultP2pPort
	}

	// rpc service listening port, if not, use the default port
	if cfg.RpcPort == "" {
		cfg.RpcPort = DefaultRpcPort
	}

	if cfg.TestNet {
		param.Net = param.TestNet
	}

	// p2p same network label, the label is different and cannot communicate
	param.UniqueNetWork = param.Net + param.UniqueNetWork

	if !utils.IsExist(cfg.HomeDir) {
		if err := os.Mkdir(cfg.HomeDir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	if cfg.DataDir == "" {
		cfg.DataDir = cfg.HomeDir + "/" + cfg.P2pPort
	}
	if !utils.IsExist(cfg.DataDir) {
		if err := os.Mkdir(cfg.DataDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// Each node requires a secp256k1 private key, which is used as the p2p id
	// generation and signature of the node that generates the block.
	// If this parameter is not configured in the startup parameter,
	// the node will be automatically generated and loaded automatically at startup
	if cfg.KeyFile == "" {
		cfg.KeyFile = defaultPrivateFile
		cfg.NodePrivate, err = LoadNodePrivate(cfg.DataDir+"/"+cfg.KeyFile, defaultKey)
		if err != nil {
			cfg.NodePrivate, err = CreateNewNodePrivate(param.Net)
			if err != nil {
				return nil, fmt.Errorf("create new node priavte failed! %s", err.Error())
			}
		}
		j, err := keystore.PrivateToJson(param.Net, cfg.NodePrivate.PrivateKey, cfg.NodePrivate.Mnemonic, []byte(defaultKey))
		if err != nil {
			return nil, fmt.Errorf("key json creation failed! %s", err.Error())
		}
		bytes, _ := json.Marshal(j)
		err = ioutil.WriteFile(cfg.DataDir+"/"+cfg.KeyFile, bytes, 0644)
		if err != nil {
			return nil, fmt.Errorf("write jsonfile failed! %s", err.Error())
		}
	} else {
		// The private key of the node is encrypted in the key file,
		// and a password is required to unlock the key file
		if cfg.KeyPass == "" {
			fmt.Println("Please enter the password for the keyfile:")
			passWd, err := readPassWd()
			if err != nil {
				return nil, fmt.Errorf("read password failed! %s", err.Error())
			}
			cfg.KeyPass = string(passWd)
		}
		cfg.NodePrivate, err = LoadNodePrivate(cfg.KeyFile, cfg.KeyPass)
		if err != nil {
			return nil, fmt.Errorf("failed to load keyfile %s! %s", cfg.KeyFile, err.Error())
		}
	}

	// If this parameter is true, the log is also written to the file
	if cfg.FileLogging {
		logDir := cfg.DataDir + "/log"
		if !utils.IsExist(logDir) {
			if err := os.Mkdir(logDir, os.ModePerm); err != nil {
				return nil, err
			}
		}
		utils.CleanAndExpandPath(logDir)
		logDir = filepath.Join(logDir, param.Net)
		common.InitLogRotator(filepath.Join(logDir, "blockchain.log"))
	}
	log15.Info("chain data directory", "path", cfg.DataDir)
	return cfg, nil
}

func newConfigParser(cfg *Config, options flags.Options) *flags.Parser {
	parser := flags.NewParser(cfg, options)
	return parser
}

func loadInitialCandidates(file string) ([]*types.Candidate, error) {
	var candidates []*types.Candidate
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, _ := rd.ReadString('\n')
		strs := strings.Split(strings.ReplaceAll(line, "\n", ""), ":")
		if len(strs) != 2 {
			break
		}
		candidates = append(candidates, &types.Candidate{Signer: hasharry.StringToAddress(strs[0]), PeerId: strs[1]})
	}
	return candidates, nil
}

// Read the password entered by stdin
func readPassWd() ([]byte, error) {
	var passWd [33]byte

	n, err := os.Stdin.Read(passWd[:])
	if err != nil {
		return nil, err
	}
	if n <= 1 {
		return nil, errors.New("not read")
	}
	return passWd[:n-1], nil
}
