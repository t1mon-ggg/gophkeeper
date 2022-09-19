package secrets

const scopeOTP = "otp"

type OTP struct {
	Method        string
	Issuer        string
	Secret        string
	AccountName   string
	RecoveryCodes []string
}

func NewOTP(method, issuer, secret, accountname string, recoverycodes ...string) *OTP {
	if method != "TOTP" && method != "HOTP" {
		return nil
	}
	otp := OTP{
		Method:        method,
		Issuer:        issuer,
		Secret:        secret,
		AccountName:   accountname,
		RecoveryCodes: recoverycodes,
	}
	return &otp
}

func (s *OTP) Scope() string {
	return scopeOTP
}

func (s *OTP) Value() interface{} {
	return s
}
