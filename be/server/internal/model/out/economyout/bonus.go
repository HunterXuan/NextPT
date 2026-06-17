package economyout

import "github.com/gogf/gf/v2/os/gtime"

type BonusLogItem struct {
	Id           uint64      `json:"id"`
	UserId       uint64      `json:"userId"`
	Amount       float64     `json:"amount"`
	BalanceAfter float64     `json:"balanceAfter"`
	Action       string      `json:"action"`
	TargetType   string      `json:"targetType"`
	TargetId     uint64      `json:"targetId"`
	Remark       string      `json:"remark"`
	CreatedAt    *gtime.Time `json:"createdAt"`
}

type BonusLogsOut struct {
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Total int            `json:"total"`
	List  []BonusLogItem `json:"list"`
}

type HourlyBonusOut struct {
	HourlyBonus float64 `json:"hourlyBonus"`
}
