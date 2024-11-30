package grpcClient

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/amirvalhalla/tron-wallet/enums"
	"github.com/amirvalhalla/tron-wallet/grpcClient/proto/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var c *GrpcClient
var singleton sync.Once

// GrpcClient controller structure
type GrpcClient struct {
	Address     string
	Conn        *grpc.ClientConn
	Client      api.WalletClient
	grpcTimeout time.Duration
	opts        []grpc.DialOption
	apiKey      string
}

func GetGrpcClient(node enums.Node, opts ...grpc.DialOption) (*GrpcClient, error) {

	var err error

	singleton.Do(func() {
		temp := &GrpcClient{
			Address:     string(node),
			grpcTimeout: 5 * time.Second,
			apiKey:      os.Getenv("TRON_PRO_API_KEY"),
			opts:        opts,
		}

		if e := temp.Start(temp.opts...); e != nil {
			err = e
		}

		c = temp
	})

	return c, err
}

func (g *GrpcClient) Start(opts ...grpc.DialOption) error {
	var err error
	if len(g.Address) == 0 {
		g.Address = "grpc.trongrid.io:50051"
	}
	g.opts = opts
	g.Conn, err = grpc.Dial(g.Address, opts...)

	if err != nil {
		return fmt.Errorf("connecting grpc client: %v", err)
	}
	g.Client = api.NewWalletClient(g.Conn)
	return nil
}

func (g *GrpcClient) getContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), g.grpcTimeout)
	if len(g.apiKey) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", g.apiKey)
	}
	return ctx, cancel
}
