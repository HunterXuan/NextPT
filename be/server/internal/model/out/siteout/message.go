package siteout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/os/gtime"
)

type MessageItem struct {
	Id         uint64               `json:"id"`
	SenderId   uint64               `json:"senderId"`
	Sender     model.IamUserSummary `json:"sender"`
	ReceiverId uint64               `json:"receiverId"`
	Receiver   model.IamUserSummary `json:"receiver"`
	Title      string               `json:"title"`
	Content    string               `json:"content"`
	TargetType string               `json:"targetType"`
	TargetId   uint64               `json:"targetId"`
	IsRead     bool                 `json:"isRead"`
	ReadAt     *gtime.Time          `json:"readAt"`
	CreatedAt  *gtime.Time          `json:"createdAt"`
}

type MessageListOut struct {
	List  []*MessageItem `json:"list"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

type MessageCreateOut struct {
	Count int `json:"count"`
}
