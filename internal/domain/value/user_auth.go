package value

type UserAuth struct {
	PersonID PersonID `db:"person_id"`
	Password Password `db:"password"`
}
