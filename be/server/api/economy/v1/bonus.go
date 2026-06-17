package v1

import (
	"server/internal/model/in/economyin"
	"server/internal/model/out/economyout"

	"github.com/gogf/gf/v2/frame/g"
)

type BonusLogsReq struct {
	g.Meta `path:"/users/me/bonus-logs" method:"get" tags:"Economy" summary:"获取魔力值流水"`
	economyin.BonusLogsInp
}

type BonusLogsRes struct {
	economyout.BonusLogsOut
}

type HourlyBonusReq struct {
	g.Meta `path:"/users/me/hourly-bonus" method:"get" tags:"Economy" summary:"获取当前用户每小时预期魔力值"`
	economyin.HourlyBonusInp
}

type HourlyBonusRes struct {
	economyout.HourlyBonusOut
}
