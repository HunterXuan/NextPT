package modout

import "github.com/gogf/gf/v2/os/gtime"

type UserModItem struct {
	Id         uint64      `json:"id"`
	UserId     uint64      `json:"userId"`
	ModType    int         `json:"modType"`
	Reason     string      `json:"reason"`
	ExpireAt   *gtime.Time `json:"expireAt"`
	ModBy      uint64      `json:"modBy"`
	ModComment string      `json:"modComment"`
	IsActive   bool        `json:"isActive"`
	CreatedAt  *gtime.Time `json:"createdAt"`
}

type ListUserOut struct {
	List []UserModItem `json:"list"`
}
