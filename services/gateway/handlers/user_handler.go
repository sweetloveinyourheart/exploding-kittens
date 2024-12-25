package handlers

import (
	"net/http"

	"connectrpc.com/connect"

	grpc "github.com/sweetloveinyourheart/planning-pocker/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/planning-pocker/services/gateway/common/requests"
	"github.com/sweetloveinyourheart/planning-pocker/services/gateway/common/responses"
)

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := requests.GetVar(r, "id")
	if err != nil {
		responses.BadRequestException(w, err)
		return
	}

	user, err := h.userServerClient.GetUser(h.ctx, connect.NewRequest(&grpc.GetUserRequest{UserId: userID}))
	if err != nil {
		responses.NotFoundException(w, err)
		return
	}

	responses.Ok(w, user)
}
