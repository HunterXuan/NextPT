package consts

import "time"

const (
	IamTwoStepTypeDisabled = 0
	IamTwoStepTypeTOTP     = 1
)

const (
	IamTwoStepSetupTTL          = 10 * time.Minute
	IamTwoStepLoginChallengeTTL = 5 * time.Minute
	IamTwoStepLoginRateWindow   = 5 * time.Minute
	IamTwoStepLoginRateLimit    = 10
	IamTwoStepRecoveryCodeCount = 10
)
