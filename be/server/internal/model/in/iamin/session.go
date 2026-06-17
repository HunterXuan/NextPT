package iamin

type SessionCreateInp struct {
	Username string `json:"username" v:"required#{#iam.session.username_req}"`
	Password string `json:"password" v:"required#{#iam.session.password_req}"`
}
