package iamin

type SessionCreateInp struct {
	Username string `json:"username" v:"required#{#iam.session.username_req}"`
	Password string `json:"password" v:"required#{#iam.session.password_req}"`
}

type SessionTwoStepVerifyInp struct {
	Challenge string `json:"challenge" v:"required|max-length:128#{#iam.two_step.challenge_req}|{#iam.two_step.challenge_invalid}"`
	Code      string `json:"code" v:"required|max-length:64#{#iam.two_step.code_req}|{#iam.two_step.code_invalid}"`
}
