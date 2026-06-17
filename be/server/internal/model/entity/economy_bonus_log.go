// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EconomyBonusLog is the golang structure for table economy_bonus_log.
type EconomyBonusLog struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`
	UserId       uint64      `json:"userId"       orm:"user_id"       description:""`
	Amount       float64     `json:"amount"       orm:"amount"        description:"正=获得 负=消耗"`
	BalanceAfter float64     `json:"balanceAfter" orm:"balance_after" description:"变动后余额"`
	Action       string      `json:"action"       orm:"action"        description:"torrent_reward/post_reward/daily_bonus/..."`
	TargetType   string      `json:"targetType"   orm:"target_type"   description:""`
	TargetId     uint64      `json:"targetId"     orm:"target_id"     description:""`
	Period       string      `json:"period"       orm:"period"        description:"结算周期/幂等键"`
	Remark       string      `json:"remark"       orm:"remark"        description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
}
