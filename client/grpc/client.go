package grpc

import (
	"context"
	"fmt"

	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	tmtype "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Client wraps GRPC client connection.
type Client struct {
	tmservice.ServiceClient
	*grpc.ClientConn
}

// NewClient creates GRPC client.
func NewClient(grpcURL string, timeout int64) (*Client, error) {
	var grpcopts []grpc.DialOption
	var conn *grpc.ClientConn
	var err error

	grpcopts = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}
	conn, err = grpc.DialContext(context.Background(), grpcURL, grpcopts...)
	if err != nil {
		return &Client{}, fmt.Errorf("failed to connect GRPC client: %s", err)
	}
	tmclient := tmservice.NewServiceClient(conn)

	return &Client{tmclient, conn}, nil
}

// IsNotFound returns not found status.
func IsNotFound(err error) bool {
	return status.Convert(err).Code() == codes.NotFound
}

func (c *Client) GetLBlock(ctx context.Context) (*tmtype.Block, error) {
	LatestBlockProto, err := c.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestBlockProto: %s", err)
	}
	LatestBlock, err := tmtype.BlockFromProto(LatestBlockProto.Block)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestBlock: %s", err)
	}
	return LatestBlock, nil
}
func (c *Client) GetLByBlock(ctx context.Context, block_num int64) (*tmtype.Block, error) {
	BlockProto, err := c.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: block_num})
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestBlockProto: %s", err)
	}
	Block, err := tmtype.BlockFromProto(BlockProto.Block)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestBlock: %s", err)
	}
	return Block, nil
}

func (c *Client) GetLValidatorSet(ctx context.Context) ([]*tmservice.Validator, error) {
	LatestValidatorSetProto, err := c.GetLatestValidatorSet(context.Background(), &tmservice.GetLatestValidatorSetRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestValidatorSetProto: %s", err)
	}
	Validators := LatestValidatorSetProto.GetValidators()

	return Validators, nil
}

func (c *Client) GetSyncInfo(ctx context.Context) (bool, error) {
	LatestSyncingProto, err := c.GetSyncing(context.Background(), &tmservice.GetSyncingRequest{})
	if err != nil {
		return false, fmt.Errorf("failed to get LatestValidatorSetProto: %s", err)
	}
	Syncing := LatestSyncingProto.GetSyncing()
	return Syncing, nil
}
