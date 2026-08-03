package iamin

type PasswordResetRequestCreateInp struct {
	Email string `json:"email" v:"required|email#{#iam.user.email_req}|{#iam.user.email_fmt}"`
}

type PasswordResetCreateInp struct {
	Token       string `json:"token" v:"required|max-length:128#{#iam.user.password_reset_token_req}|{#iam.user.password_reset_token_invalid}"`
	NewPassword string `json:"newPassword" v:"required|length:6,30#{#iam.user.new_password_req}|{#iam.user.password_len}"`
}
