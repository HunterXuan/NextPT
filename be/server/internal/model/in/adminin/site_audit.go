package adminin

import "github.com/gogf/gf/v2/os/gtime"

type SiteAuditListInp struct {
	Page       int         `json:"page" d:"1"`
	Size       int         `json:"size" d:"20"`
	Level      *int        `json:"level" in:"query"`
	Action     string      `json:"action" in:"query"`
	TargetType string      `json:"targetType" in:"query"`
	UserId     uint64      `json:"userId" in:"query"`
	StartAt    *gtime.Time `json:"startAt" in:"query"`
	EndAt      *gtime.Time `json:"endAt" in:"query"`
}
