package model

type IamUserSummary struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type IamUserPermissionListOptions struct {
	SourceType   *int
	WildcardOnly bool
	Page         int
	Size         int
}
