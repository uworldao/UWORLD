package config

import "github.com/jhdriver/UWORLD/config"

type Config struct {
	ConfigFile  string
	Format      bool
	TestNet     bool
	KeyStoreDir string
	config.RpcConfig
}
