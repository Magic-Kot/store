package middlewarex

type sensitiveDataMasker interface {
	Mask([]byte) []byte
}
