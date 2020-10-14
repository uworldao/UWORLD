package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jhdriver/UWORLD/cmd/wallet/command"
	"github.com/jhdriver/UWORLD/cmd/wallet/config"
	"github.com/jhdriver/UWORLD/common/utils"
	config2 "github.com/jhdriver/UWORLD/config"
	"github.com/jhdriver/UWORLD/param"
	"github.com/spf13/cobra"
	"os"
)

var preConfig *config.Config
var (
	defaultFormat      = true
	defaultTestNet     = false
	defaultKeyStoreDir = "keystore"
	defaultRpcCer      = "server.pem"
	defaultRpcIp       = "127.0.0.1"
)

func init() {
	cobra.EnableCommandSorting = true
	bindFlags()
}

func main() {
	command.RootCmd.PersistentPreRunE = LoadConfig
	if err := command.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func bindFlags() {
	preConfig = &config.Config{}
	gFlags := command.RootCmd.PersistentFlags()

	gFlags.StringVarP(&preConfig.ConfigFile, "config", "c", "wallet.toml", "Wallet profile")
}

// LoadConfig config file and flags
func LoadConfig(cmd *cobra.Command, args []string) (err error) {
	fileCfg := &config.Config{}
	_, err = toml.DecodeFile(preConfig.ConfigFile, fileCfg)
	if err != nil {
		if !cmd.Flag("config").Changed {
			if fExit := utils.IsExist(preConfig.ConfigFile); !fExit {
				return fmt.Errorf("config file is not exist")
			}
			_, err = toml.DecodeFile(cmd.Flag("config").Value.String(), fileCfg)
			if err != nil {
				return fmt.Errorf("config file %s is not exist", cmd.Flag("config").Value.String())
			}
		}
	}
	if fileCfg.KeyStoreDir == "" {
		fileCfg.KeyStoreDir = defaultKeyStoreDir
	}

	if fileCfg.RpcPort == "" {
		fileCfg.RpcPort = config2.DefaultRpcPort
	}

	if fileCfg.RpcIp == "" {
		fileCfg.RpcIp = defaultRpcIp
	}

	command.Cfg = fileCfg
	if command.Cfg.TestNet {
		command.Net = param.TestNet
	}
	return nil
}
