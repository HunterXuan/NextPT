// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EconomyBonusLog is the golang structure of table economy_bonus_log for DAO operations like Where/Data.
type EconomyBonusLog struct {
	g.Meta       `orm:"table:economy_bonus_log, do:true"`
	Id           any         //
	UserId       any         //
	Amount       any         // 正=获得 负=消耗
	BalanceAfter any         // 变动后余额
	Action       any         // torrent_reward/post_reward/daily_bonus/...
	TargetType   any         //
	TargetId     any         //
	Period       any         // 结算周期/幂等键
	Remark       any         //
	CreatedAt    *gtime.Time //
}
