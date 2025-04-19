package entity

type PostAuthSignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
