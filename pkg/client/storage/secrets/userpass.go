package secrets

const scopeUserPass = "user-password"

type UserPass struct {
	Username string
	Password string
}

func NewUserPass(username, password string) *UserPass {
	up := UserPass{
		Username: username,
		Password: password,
	}
	return &up
}

func (s *UserPass) Scope() string {
	return scopeUserPass
}

func (s *UserPass) Value() interface{} {
	return s
}
