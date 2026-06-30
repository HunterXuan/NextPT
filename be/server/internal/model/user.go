package model

type UserSummary struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
