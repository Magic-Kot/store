package masker

type NopSensitiveDataMasker struct{}

func NewNopSensitiveDataMasker() NopSensitiveDataMasker {
	return NopSensitiveDataMasker{}
}

func (NopSensitiveDataMasker) Mask(input []byte) []byte {
	return input
}
