package value

type RefreshToken string

func (r RefreshToken) String() string {
	return string(r)
}
