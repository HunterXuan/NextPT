// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EconomyRewardRecord is the golang structure of table economy_reward_record for DAO operations like Where/Data.
type EconomyRewardRecord struct {
	g.Meta     `orm:"table:economy_reward_record, do:true"`
	Id         any         //
	TargetType any         // catalog_torrent/catalog_comment/forum_topic/forum_reply
	TargetId   any         //
	FromUserId any         // 赞赏者ID
	ToUserId   any         // 接收者ID
	Amount     any         // 赞赏金额
	CreatedAt  *gtime.Time //
}
