package value

type Password string

func (p Password) String() string { return string(p) }
