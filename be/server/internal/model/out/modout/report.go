package modout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type ReportTargetSummary struct {
	Type       string               `json:"type"`
	Id         uint64               `json:"id"`
	Title      string               `json:"title"`
	ParentType string               `json:"parent_type"`
	ParentId   uint64               `json:"parent_id"`
	Status     string               `json:"status"`
	Author     model.IamUserSummary `json:"author"`
}

type ReportItem struct {
	Id           uint64               `json:"id"`
	ReporterId   uint64               `json:"reporter_id"`
	Reporter     model.IamUserSummary `json:"reporter"`
	TargetType   string               `json:"target_type"`
	TargetId     uint64               `json:"target_id"`
	Target       ReportTargetSummary  `json:"target"`
	Reason       string               `json:"reason"`
	Status       int                  `json:"status"`
	DealtBy      uint64               `json:"dealt_by"`
	DealtUser    model.IamUserSummary `json:"dealt_user"`
	DealtComment string               `json:"dealt_comment"`
	DealtAt      *gtime.Time          `json:"dealt_at"`
	CreatedAt    *gtime.Time          `json:"created_at"`
}

type ListReportsOut struct {
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
	List  []ReportItem `json:"list"`
}
