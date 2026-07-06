package adminout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type SiteAuditItem struct {
	Id         uint64               `json:"id"`
	UserId     uint64               `json:"userId"`
	Actor      model.IamUserSummary `json:"actor"`
	Action     string               `json:"action"`
	TargetType string               `json:"targetType"`
	TargetId   uint64               `json:"targetId"`
	Level      int                  `json:"level"`
	Ip         string               `json:"ip"`
	Detail     string               `json:"detail"`
	CreatedAt  *gtime.Time          `json:"createdAt"`
}

type SiteAuditListOut struct {
	List  []*SiteAuditItem `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}
