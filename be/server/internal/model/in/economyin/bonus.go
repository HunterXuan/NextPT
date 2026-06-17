package economyin

import "github.com/gogf/gf/v2/os/gtime"

type BonusLogsInp struct {
	Page   int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size   int    `json:"size" d:"20" v:"min:1|max:100" description:"每页数量"`
	Action string `json:"action" description:"变动类型，如 torrent_reward_sent"`
}

type HourlyBonusInp struct{}

type BonusFormulaConfig struct {
	T0, N0, B0, L, BasePoints float64
}

type BonusPeerSnapshot struct {
	UserId    uint64
	TorrentId uint64
}

type BonusTorrentSnapshot struct {
	Id        uint64
	Size      uint64
	Seeders   uint
	CreatedAt *gtime.Time
}
