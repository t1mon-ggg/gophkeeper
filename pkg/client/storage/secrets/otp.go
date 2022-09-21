package secrets

const scopeOTP = "otp"

// OTP - type for otp secret
type OTP struct {
	Method        string
	Issuer        string
	Secret        string
	AccountName   string
	RecoveryCodes []string
}

// NewOTP - create otp secret
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

// Scope - secret scope
func (s *OTP) Scope() string {
	return scopeOTP
}

// Value - secret value
func (s *OTP) Value() interface{} {
	return s
}
