package config

type RpcConfig struct {
	DataDir    string
	RpcIp      string
	RpcPort    string
	RpcTLS     bool
	RpcCert    string
	RpcCertKey string
	RpcPass    string
}
