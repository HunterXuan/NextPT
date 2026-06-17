package iamin

type UserCreateInp struct {
	Username   string `json:"username" v:"required|length:3,20#{#iam.user.username_req}|{#iam.user.username_len}"`
	Email      string `json:"email" v:"required|email#{#iam.user.email_req}|{#iam.user.email_fmt}"`
	Password   string `json:"password" v:"required|length:6,30#{#iam.user.password_req}|{#iam.user.password_len}"`
	InviteHash string `json:"inviteHash"`
}

type UserProfileUpdateInp struct {
	Avatar    string `json:"avatar" v:"max-length:500#{#iam.user.avatar_len}"`
	Info      string `json:"info" v:"max-length:5000#{#iam.user.info_len}"`
	Signature string `json:"signature" v:"max-length:500#{#iam.user.signature_len}"`
}

type UserPasswordChangeInp struct {
	OldPassword string `json:"oldPassword" v:"required#{#iam.user.old_password_req}"`
	NewPassword string `json:"newPassword" v:"required|length:6,30#{#iam.user.new_password_req}|{#iam.user.password_len}"`
}
