package secrets

import (
	"time"
)

const scopeCC = "creditcard"

type CreditCard struct {
	Number string
	Holder string
	CVV    uint16
	Expire time.Time
}

func NewCC(number, holder, expire string, cvv uint16) (*CreditCard, error) {
	exp, err := time.Parse("01/06", expire)
	if err != nil {
		return &CreditCard{}, err
	}
	cc := CreditCard{
		Number: number,
		Holder: holder,
		CVV:    cvv,
		Expire: exp.AddDate(0, 1, 0),
	}
	return &cc, nil
}

func (s *CreditCard) Scope() string {
	return scopeCC
}

func (s *CreditCard) Value() interface{} {
	return s
}
