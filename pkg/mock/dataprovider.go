package mock

import (
	"context"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"

	dataProviderProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
	dataProviderGrpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
)

type MockDataProviderClient struct {
	mock.Mock
}

var _ dataProviderGrpc.DataProviderClient = (*MockDataProviderClient)(nil)

func (mdpc *MockDataProviderClient) GetCards(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[dataProviderProto.GetCardsResponse], error) {
	args := mdpc.Called(ctx, req)
	return args.Get(0).(*connect.Response[dataProviderProto.GetCardsResponse]), args.Error(1)
}

func (mdpc *MockDataProviderClient) GetMapCards(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[dataProviderProto.GetMapCardsResponse], error) {
	args := mdpc.Called(ctx, req)
	return args.Get(0).(*connect.Response[dataProviderProto.GetMapCardsResponse]), args.Error(1)
}
