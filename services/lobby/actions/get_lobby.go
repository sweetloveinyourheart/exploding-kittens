package actions

import (
	"context"

	"connectrpc.com/connect"

	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/lobbyserver/go"
)

func (a *actions) GetLobbyData(ctx context.Context, request *connect.Request[proto.GetLobbyDataRequest]) (response *connect.Response[proto.GetLobbyDataResponse], err error) {
	return nil, nil
}
