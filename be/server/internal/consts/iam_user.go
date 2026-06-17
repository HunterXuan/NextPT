package consts

const (
	IamUserStatusPending   = iota // 0: 待验证
	IamUserStatusConfirmed        // 1: 正常
	IamUserStatusDisabled         // 2: 被禁用
)

const (
	IamUserPrivacyLoose  = iota // 0: 宽松
	IamUserPrivacyNormal        // 1: 普通
	IamUserPrivacyStrict        // 2: 严格
)
