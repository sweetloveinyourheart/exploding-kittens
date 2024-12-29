package schemas

type CreateNewGuestUserRequest struct {
	Username string `json:"username" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Status   int    `json:"status" validate:"required,oneof=0 1 2"`
}
