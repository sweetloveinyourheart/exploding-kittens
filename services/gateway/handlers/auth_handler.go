package handlers

import (
	"net/http"

	"connectrpc.com/connect"

	grpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/gateway/common/requests"
	"github.com/sweetloveinyourheart/exploding-kittens/services/gateway/common/responses"
	"github.com/sweetloveinyourheart/exploding-kittens/services/gateway/schemas"
)

func (h *handler) GuestLogin(w http.ResponseWriter, r *http.Request) {
	credentials, err := requests.ParseBodyWithValidation[schemas.GuestLoginRequest](r)
	if err != nil {
		responses.BadRequestException(w, err)
		return
	}

	resp, err := h.userServerClient.SignIn(h.ctx, connect.NewRequest(&grpc.SignInRequest{
		UserId: credentials.UserID,
	}))
	if err != nil {
		responses.BadRequestException(w, err)
		return
	}

	responses.Ok(w, resp)
}
