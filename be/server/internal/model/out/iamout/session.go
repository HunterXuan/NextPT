package iamout

import "server/internal/model"

type RegisterOut struct {
}

type SessionCreateOut struct {
	Token            string `json:"token"`
	TwoStepRequired  bool   `json:"twoStepRequired"`
	TwoStepChallenge string `json:"twoStepChallenge"`
}

type SessionListOut struct {
	List []model.IamSessionItem `json:"list"`
}
