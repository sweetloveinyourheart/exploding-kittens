package schemas

type GuestLoginRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}
