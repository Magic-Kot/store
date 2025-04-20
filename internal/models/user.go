package models

type UserLogin struct {
	ID       int    `json:"id"`
	Username string `json:"login"        validate:"required,min=4,max=20"`
	Password string `json:"password"     validate:"required,min=6,max=20"`
}
