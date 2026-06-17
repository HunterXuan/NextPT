package modout

import "github.com/gogf/gf/v2/os/gtime"

type UserModItem struct {
	Id         uint64      `json:"id"`
	UserId     uint64      `json:"user_id"`
	ModType    int         `json:"mod_type"`
	Reason     string      `json:"reason"`
	ExpireAt   *gtime.Time `json:"expire_at"`
	ModBy      uint64      `json:"mod_by"`
	ModComment string      `json:"mod_comment"`
	IsActive   bool        `json:"is_active"`
	CreatedAt  *gtime.Time `json:"created_at"`
}

type ListUserOut struct {
	List []UserModItem `json:"list"`
}
