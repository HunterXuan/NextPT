package forumout

import (
	"server/internal/model"

	"github.com/gogf/gf/v2/encoding/gjson"
)

type NodeCategoryItem struct {
	Id       uint        `json:"id"`
	NameI18N *gjson.Json `json:"nameI18n"`
	DescI18N *gjson.Json `json:"descI18n"`
	Nodes    []NodeItem  `json:"nodes"`
}

type NodeListOut struct {
	List []NodeCategoryItem `json:"list"`
}

type NodeItem struct {
	Id         uint        `json:"id"`
	Slug       string      `json:"slug"`
	NameI18N   *gjson.Json `json:"nameI18n"`
	DescI18N   *gjson.Json `json:"descI18n"`
	TopicCount uint        `json:"topicCount"`
	ReplyCount uint        `json:"replyCount"`
}

type TopicListItem struct {
	Id            uint64               `json:"id"`
	Subject       string               `json:"subject"`
	Author        model.IamUserSummary `json:"author"`
	IsLocked      bool                 `json:"isLocked"`
	IsSticky      bool                 `json:"isSticky"`
	Views         uint                 `json:"views"`
	ReplyCount    uint                 `json:"replyCount"`
	LastReplyAt   string               `json:"lastReplyAt"`
	LastReplyUser model.IamUserSummary `json:"lastReplyUser"`
	CreatedAt     string               `json:"createdAt"`
}

type TopicListOut struct {
	List  []TopicListItem `json:"list"`
	Total int             `json:"total"`
	Node  NodeItem        `json:"node"`
}

type TopicBookmarkListOut struct {
	List  []TopicListItem `json:"list"`
	Total int             `json:"total"`
}

type TopicDetailOut struct {
	TopicListItem
	Content      string      `json:"content"`
	Appends      *gjson.Json `json:"appends"`
	NodeId       uint        `json:"nodeId"`
	IsLiked      bool        `json:"isLiked"`
	IsBookmarked bool        `json:"isBookmarked"`
}

type ReplyListItem struct {
	Id        uint64               `json:"id"`
	Author    model.IamUserSummary `json:"author"`
	Content   string               `json:"content"`
	CreatedAt string               `json:"createdAt"`
	IsLiked   bool                 `json:"isLiked"`
}

type ReplyListOut struct {
	List  []ReplyListItem `json:"list"`
	Total int             `json:"total"`
}

type TopicCreateOut struct {
	Id uint64 `json:"id"`
}

type ReplyCreateOut struct {
	Id uint64 `json:"id"`
}
