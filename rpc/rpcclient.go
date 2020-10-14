package rpc

import (
	"context"
	"fmt"
	"github.com/jhdriver/UWORLD/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const timeout = 30

type Client struct {
	conn *grpc.ClientConn
	Gc   GreeterClient
	cfg  *config.RpcConfig
}

func NewClient(cfg *config.RpcConfig) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Connect() error {
	var conn *grpc.ClientConn
	var err error
	var opts []grpc.DialOption
	if c.cfg.RpcTLS {
		creds, err := credentials.NewClientTLSFromFile(c.cfg.RpcCert, "")
		if err != nil {
			return fmt.Errorf("failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithPerRPCCredentials(&customCredential{Password: c.cfg.RpcPass, OpenTLS: c.cfg.RpcTLS}))

	conn, err = grpc.Dial(c.cfg.RpcIp+":"+c.cfg.RpcPort, opts...)
	if err != nil {
		return err
	}

	c.conn = conn
	c.Gc = NewGreeterClient(c.conn)
	return nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

type customCredential struct {
	OpenTLS  bool
	Password string
}

func (c *customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"password": c.Password,
	}, nil
}

func (c *customCredential) RequireTransportSecurity() bool {
	if c.OpenTLS {
		return true
	}

	return false
}
