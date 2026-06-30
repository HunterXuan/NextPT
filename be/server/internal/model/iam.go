package model

type IamUserSummary struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
