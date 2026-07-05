package adminout

import "github.com/gogf/gf/v2/os/gtime"

type IamPermissionListOut struct {
	Permissions []string `json:"permissions"`
}

type IamUserPermissionDetailOut struct {
	UserAcls []IamUserAclSummary `json:"userAcls"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	Size     int                 `json:"size"`
}

type IamUserAclSummary struct {
	Id         uint64      `json:"id"`
	UserId     uint64      `json:"userId"`
	PermKey    string      `json:"permKey"`
	RawPermKey string      `json:"rawPermKey"`
	IsDeny     bool        `json:"isDeny"`
	SourceType int         `json:"sourceType"`
	SourceId   uint64      `json:"sourceId"`
	ExpireAt   *gtime.Time `json:"expireAt"`
	IsActive   bool        `json:"isActive"`
	CreatedAt  *gtime.Time `json:"createdAt"`
}
