package modout

import (
	"server/internal/consts"
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type ReportTargetSummary struct {
	Type       string               `json:"type"`
	Id         uint64               `json:"id"`
	Title      string               `json:"title"`
	ParentType string               `json:"parentType"`
	ParentId   uint64               `json:"parentId"`
	Status     string               `json:"status"`
	Author     model.IamUserSummary `json:"author"`
}

func (m ReportTargetSummary) SiteMessageTarget() (string, uint64) {
	if m.Status != consts.ModReportTargetStatusNormal {
		return "", 0
	}

	targetType, targetId := m.Type, m.Id
	if m.ParentType != "" {
		targetType, targetId = m.ParentType, m.ParentId
	}

	switch targetType {
	case consts.SiteMessageTargetTypeCatalogTorrent, consts.SiteMessageTargetTypeForumTopic:
		return targetType, targetId
	default:
		return "", 0
	}
}

type ReportItem struct {
	Id           uint64               `json:"id"`
	ReporterId   uint64               `json:"reporterId"`
	Reporter     model.IamUserSummary `json:"reporter"`
	TargetType   string               `json:"targetType"`
	TargetId     uint64               `json:"targetId"`
	Target       ReportTargetSummary  `json:"target"`
	Reason       string               `json:"reason"`
	Status       int                  `json:"status"`
	DealtBy      uint64               `json:"dealtBy"`
	DealtUser    model.IamUserSummary `json:"dealtUser"`
	DealtComment string               `json:"dealtComment"`
	DealtAt      *gtime.Time          `json:"dealtAt"`
	CreatedAt    *gtime.Time          `json:"createdAt"`
}

type ListReportsOut struct {
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
	List  []ReportItem `json:"list"`
}
