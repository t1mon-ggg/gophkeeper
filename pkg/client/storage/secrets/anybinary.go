package secrets

const scopeAnyBinary = "anybinary"

type AnyBinary struct {
	Bytes []byte
}

func NewBinary(b []byte) *AnyBinary {
	bb := AnyBinary{
		Bytes: b,
	}
	return &bb
}

func (s *AnyBinary) Scope() string {
	return scopeAnyBinary
}

func (s *AnyBinary) Value() interface{} {
	return s
}
