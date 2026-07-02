package economyout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type RewardSummary struct {
	UserId       uint64      `json:"userId" orm:"user_id"`
	Amount       float64     `json:"amount" orm:"amount"`
	RewardCount  uint        `json:"rewardCount" orm:"reward_count"`
	LastRewardAt *gtime.Time `json:"lastRewardAt" orm:"last_reward_at"`
}

type RewardItem struct {
	User         model.IamUserSummary `json:"user" description:"赞赏者"`
	Amount       float64              `json:"amount" description:"累计赞赏金额"`
	RewardCount  uint                 `json:"rewardCount" description:"赞赏次数"`
	LastRewardAt string               `json:"lastRewardAt" description:"最近赞赏时间"`
}

type RewardListOut struct {
	List  []RewardItem `json:"list" description:"赞赏列表"`
	Total int          `json:"total" description:"总记录数"`
}
