package v1

import (
	"server/internal/model/in/sitein"
	"server/internal/model/out/siteout"

	"github.com/gogf/gf/v2/frame/g"
)

type MessageListReq struct {
	g.Meta `path:"/messages" method:"get" tags:"Site" summary:"获取站内消息" perm:"read:site/message:*"`
	sitein.MessageListInp
}

type MessageListRes struct {
	siteout.MessageListOut
}

type MessageReadReq struct {
	g.Meta `path:"/messages/{id}:read" method:"post" tags:"Site" summary:"标记消息已读" perm:"read:site/message:*"`
	sitein.MessageReadInp
}

type MessageReadRes struct{}

type MessageReadAllReq struct {
	g.Meta `path:"/messages:readAll" method:"post" tags:"Site" summary:"标记全部消息已读" perm:"read:site/message:*"`
	sitein.MessageReadAllInp
}

type MessageReadAllRes struct{}

type ChatMessageListReq struct {
	g.Meta `path:"/chat-messages" method:"get" tags:"Site" summary:"获取聊天室消息" perm:"read:site/chat-message:*"`
	sitein.ChatMessageListInp
}

type ChatMessageListRes struct {
	siteout.ChatMessageListOut
}

type ChatMessageCreateReq struct {
	g.Meta `path:"/chat-messages" method:"post" tags:"Site" summary:"发送聊天室消息" perm:"create:site/chat-message:*"`
	sitein.ChatMessageCreateInp
}

type ChatMessageCreateRes struct {
	siteout.ChatMessageItem
}
