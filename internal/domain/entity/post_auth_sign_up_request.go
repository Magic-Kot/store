package entity

type PostAuthSignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
