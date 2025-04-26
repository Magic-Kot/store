package value

type AccessToken string

func (a AccessToken) String() string {
	return string(a)
}
