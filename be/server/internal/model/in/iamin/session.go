package iamin

type SessionCreateInp struct {
	Username string `json:"username" v:"required#{#iam.session.username_req}"`
	Password string `json:"password" v:"required#{#iam.session.password_req}"`
}

type SessionTwoStepVerifyInp struct {
	Challenge string `json:"challenge" v:"required|max-length:128#{#iam.two_step.challenge_req}|{#iam.two_step.challenge_invalid}"`
	Code      string `json:"code" v:"required|max-length:64#{#iam.two_step.code_req}|{#iam.two_step.code_invalid}"`
}

type SessionDeleteInp struct {
	Id string `json:"id" in:"path" v:"required|length:32,32#{#iam.session.id_req}|{#iam.session.id_invalid}"`
}
