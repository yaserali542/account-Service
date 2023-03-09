package models

type AccountWithJWT struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture []byte `json:"profile_picture"`
	UserName       string `json:"user_name"`
	EmailAddress   string `json:"email_address"`
	JwtToken
}

type JwtToken struct {
	Token string `json:"jwt-token"`
}
