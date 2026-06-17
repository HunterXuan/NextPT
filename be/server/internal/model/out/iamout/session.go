package iamout

type RegisterOut struct {
}

type SessionCreateOut struct {
	Token string `json:"token" dc:"认证凭证 JWT"`
}
