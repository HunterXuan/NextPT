// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ModUserLog is the golang structure for table mod_user_log.
type ModUserLog struct {
	Id         uint64      `json:"id"         orm:"id"          description:""`
	UserId     uint64      `json:"userId"     orm:"user_id"     description:""`
	ModType    int         `json:"modType"    orm:"mod_type"    description:"1=warned 2=banned 3=leech_warned 4=upload_banned 5=download_banned 6=forum_banned"`
	Reason     string      `json:"reason"     orm:"reason"      description:""`
	ExpireAt   *gtime.Time `json:"expireAt"   orm:"expire_at"   description:"过期时间，NULL=永久"`
	ModBy      uint64      `json:"modBy"      orm:"mod_by"      description:"操作人ID"`
	ModComment string      `json:"modComment" orm:"mod_comment" description:""`
	IsActive   bool        `json:"isActive"   orm:"is_active"   description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`
}
