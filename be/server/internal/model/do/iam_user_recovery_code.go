// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserRecoveryCode is the golang structure of table iam_user_recovery_code for DAO operations like Where/Data.
type IamUserRecoveryCode struct {
	g.Meta    `orm:"table:iam_user_recovery_code, do:true"`
	Id        any         //
	UserId    any         //
	CodeHash  any         //
	UsedAt    *gtime.Time //
	CreatedAt *gtime.Time //
}
