package consts

const (
	IamLoginLogResultFail    = 0
	IamLoginLogResultSuccess = 1
)

const (
	IamLoginLogFailReasonUserNotFound       = "user_not_found"
	IamLoginLogFailReasonInvalidPassword    = "invalid_password"
	IamLoginLogFailReasonAccountUnavailable = "account_unavailable"
	IamLoginLogFailReasonRoleMissing        = "role_missing"
	IamLoginLogFailReasonTokenCreateFailed  = "token_create_failed"
)
