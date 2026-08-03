package consts

import "time"

const (
	IamUserStatusPending   = iota // 0: 待验证
	IamUserStatusConfirmed        // 1: 正常
	IamUserStatusDisabled         // 2: 被禁用
)

const (
	IamPasswordResetTokenTTL       = 30 * time.Minute
	IamPasswordResetRateWindow     = 15 * time.Minute
	IamPasswordResetRateLimitByIp  = 10
	IamPasswordResetRateLimitEmail = 3
)

const (
	IamUserPrivacyLoose  = iota // 0: 宽松
	IamUserPrivacyNormal        // 1: 普通
	IamUserPrivacyStrict        // 2: 严格
)

const (
	IamUserRankAutoBanReason = "rank_demotion_below_minimum"
)
