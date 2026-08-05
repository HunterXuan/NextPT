// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserRecoveryCode is the golang structure for table iam_user_recovery_code.
type IamUserRecoveryCode struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`
	CodeHash  string      `json:"codeHash"  orm:"code_hash"  description:""`
	UsedAt    *gtime.Time `json:"usedAt"    orm:"used_at"    description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
