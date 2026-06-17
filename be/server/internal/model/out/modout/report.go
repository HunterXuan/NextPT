package modout

import "github.com/gogf/gf/v2/os/gtime"

type ReportItem struct {
	Id           uint64      `json:"id"`
	ReporterId   uint64      `json:"reporter_id"`
	TargetType   string      `json:"target_type"`
	TargetId     uint64      `json:"target_id"`
	Reason       string      `json:"reason"`
	Status       int         `json:"status"`
	DealtBy      uint64      `json:"dealt_by"`
	DealtComment string      `json:"dealt_comment"`
	DealtAt      *gtime.Time `json:"dealt_at"`
	CreatedAt    *gtime.Time `json:"created_at"`
}

type ListReportsOut struct {
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
	List  []ReportItem `json:"list"`
}
