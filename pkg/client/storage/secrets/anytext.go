package secrets

const scopeAnyText = "anytext"

type AnyText struct {
	Text string
}

func NewText(text string) *AnyText {
	txt := AnyText{
		Text: text,
	}
	return &txt
}

func (s *AnyText) Scope() string {
	return scopeAnyText
}

func (s *AnyText) Value() interface{} {
	return s
}
