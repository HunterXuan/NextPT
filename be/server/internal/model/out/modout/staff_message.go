package modout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type StaffMessageItem struct {
	Id             uint64               `json:"id"`
	SenderId       uint64               `json:"senderId"`
	Sender         model.IamUserSummary `json:"sender"`
	Subject        string               `json:"subject"`
	Content        string               `json:"content"`
	Status         int                  `json:"status"`
	AnsweredBy     uint64               `json:"answeredBy"`
	AnsweredByUser model.IamUserSummary `json:"answeredByUser"`
	Answer         string               `json:"answer"`
	AnsweredAt     *gtime.Time          `json:"answeredAt"`
	CreatedAt      *gtime.Time          `json:"createdAt"`
	UpdatedAt      *gtime.Time          `json:"updatedAt"`
}

type StaffMessageListOut struct {
	List  []StaffMessageItem `json:"list"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

type StaffMessageCreateOut struct {
	Id uint64 `json:"id"`
}
