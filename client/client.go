package client

import (
	"github.com/b-harvest/indep_node_alarm_go/client/grpc"
	"github.com/b-harvest/indep_node_alarm_go/client/rpc"
)

var (
	DefaultRPCTimeout  = int64(5)
	DefaultGRPCTimeout = int64(5)
)

// Client is a wrapper for various clients.
type Client struct {
	RPC  *rpc.Client
	GRPC *grpc.Client
}

// NewClient creates a new Client with the given configuration.
func NewClient(rpcURL string, grpcURL string) (*Client, error) {

	//log.Debug().Msg("connecting clients")

	rpcClient, err := rpc.NewClient(rpcURL, DefaultRPCTimeout)
	if err != nil {
		return &Client{}, err
	}

	grpcClient, err := grpc.NewClient(grpcURL, DefaultGRPCTimeout)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		RPC:  rpcClient,
		GRPC: grpcClient,
	}, nil
}

// GetRPCClient returns RPC client.
func (c *Client) GetRPCClient() *rpc.Client {
	return c.RPC
}

// GetGRPCClient returns GRPC client.
func (c *Client) GetGRPCClient() *grpc.Client {
	return c.GRPC
}

// Stop defers the node stop execution to the RPC and GRPC clients.
func (c Client) Stop() error {
	err := c.RPC.Stop()
	if err != nil {
		return err
	}

	err = c.GRPC.Close()
	if err != nil {
		return err
	}
	return nil
}
