package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sweetloveinyourheart/planning-pocker/services/gateway/common/responses"
)

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	routerVars := mux.Vars(r)
	userID := routerVars["id"]
	if userID == "" {
		responses.BadRequestException(w, "User ID is required")
		return
	}

	responses.JSON(w, 200, "Get user successfully", true)
}
