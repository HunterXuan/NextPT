package iamin

type EmailVerificationRequestCreateInp struct {
	Email string `json:"email" v:"required|email#{#iam.user.email_req}|{#iam.user.email_fmt}"`
}

type EmailVerificationCreateInp struct {
	Token string `json:"token" v:"required|max-length:128#{#iam.user.email_verification_token_req}|{#iam.user.email_verification_token_invalid}"`
}
