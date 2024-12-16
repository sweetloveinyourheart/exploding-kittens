package actions

import (
	"context"

	"connectrpc.com/connect"

	proto "github.com/sweetloveinyourheart/planning-pocker/proto/code/userserver/go"
)

func (a *actions) GetUser(ctx context.Context, request *connect.Request[proto.GetUserRequest]) (response *connect.Response[proto.GetUserResponse], err error) {
	return nil, nil
}
