package models

type User struct {
	ID       int    `json:"id"           validate:"required,min=1"`
	Username string `json:"login"        validate:"min=4,max=20"`
	Email    string `json:"email"        validate:"email"`
	//Age       int    `json:"age"         validate:"gte=0,lte=120"`
}

type UserLogin struct {
	ID       int    `json:"id"`
	Username string `json:"login"        validate:"required,min=4,max=20"`
	Password string `json:"password"     validate:"required,min=6,max=20"`
	//CreditCard string `json:"credit_card"`
}

type UserAuthorization struct {
	ID       int    `json:"id"`
	Username string `json:"login"        validate:"required,min=1,max=20"`
	Password string `json:"password"     validate:"required,min=1,max=20"`
}
