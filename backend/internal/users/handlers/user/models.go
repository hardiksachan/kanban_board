package user

type UpdateRequest struct {
	ProfileURL  string `json:"profile_url,omitempty"`
	DisplayName string `json:"display_name,omitempty" validate:"required,min=2,max=20"`
}

type GetResponse struct {
	UserID      string `json:"user_id"`
	ProfileURL  string `json:"profile_url,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}
