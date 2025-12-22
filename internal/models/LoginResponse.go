package models

type LoginResponse struct {
	Model
	Token string `json:"token"`
	User  Users  `json:"user"`
}
