package entity

type UserData struct {
	ID       int    `db:"id"`
	Age      int    `db:"age"`
	Username string `db:"login"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	Email    string `db:"email"`
	Avatar   string `db:"avatar"`
}
