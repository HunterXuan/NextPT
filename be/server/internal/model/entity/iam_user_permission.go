// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserPermission is the golang structure for table iam_user_permission.
type IamUserPermission struct {
	Id         uint64      `json:"id"         orm:"id"          description:""`
	UserId     uint64      `json:"userId"     orm:"user_id"     description:""`
	PermKey    string      `json:"permKey"    orm:"perm_key"    description:"权限标识符，如 update:torrent:123 或 update:torrent:*"`
	SourceType int         `json:"sourceType" orm:"source_type" description:"来源类型: 1=manual 2=user_mod"`
	SourceId   uint64      `json:"sourceId"   orm:"source_id"   description:"来源记录ID，manual=0，user_mod=mod_user_log.id"`
	ExpireAt   *gtime.Time `json:"expireAt"   orm:"expire_at"   description:"权限过期时间，NULL=永久"`
	IsActive   bool        `json:"isActive"   orm:"is_active"   description:"是否生效"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`
}
