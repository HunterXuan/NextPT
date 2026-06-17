// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCommentReward is the golang structure of table catalog_comment_reward for DAO operations like Where/Data.
type CatalogCommentReward struct {
	g.Meta    `orm:"table:catalog_comment_reward, do:true"`
	Id        any         //
	UserId    any         //
	CommentId any         //
	Amount    any         // 打赏金额
	CreatedAt *gtime.Time //
}
