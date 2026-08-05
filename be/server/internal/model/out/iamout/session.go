package iamout

type RegisterOut struct {
}

type SessionCreateOut struct {
	Token            string `json:"token"`
	TwoStepRequired  bool   `json:"twoStepRequired"`
	TwoStepChallenge string `json:"twoStepChallenge"`
}
