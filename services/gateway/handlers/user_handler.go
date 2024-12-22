package handlers

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/gorilla/mux"

	grpc "github.com/sweetloveinyourheart/planning-pocker/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/planning-pocker/services/gateway/common/responses"
)

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	routerVars := mux.Vars(r)

	userID := routerVars["id"]
	if userID == "" {
		responses.BadRequestException(w, "User ID is required")
		return
	}

	user, err := h.userServerClient.GetUser(h.ctx,
		connect.NewRequest(&grpc.GetUserRequest{
			UserId: userID,
		}),
	)
	if err != nil {
		responses.NotFoundException(w, err.Error())
		return
	}

	responses.Ok(w, user)
}
