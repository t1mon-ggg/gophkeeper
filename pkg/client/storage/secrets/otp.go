package secrets

const scopeOTP = "otp"

type OTP struct {
	Method        string
	Issuer        string
	Secret        string
	RecoveryCodes []string
}

func NewOTP(method, issuer, secret string, rc ...string) *OTP {
	otp := OTP{
		Method:        method,
		Issuer:        issuer,
		Secret:        secret,
		RecoveryCodes: rc,
	}
	return &otp
}

func (s *OTP) Scope() string {
	return scopeOTP
}

func (s *OTP) Value() interface{} {
	return s
}
