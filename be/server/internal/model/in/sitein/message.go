package sitein

type MessageListInp struct {
	Page   int   `json:"page" d:"1" v:"min:1"`
	Size   int   `json:"size" d:"20" v:"min:1|max:100"`
	IsRead *bool `json:"isRead" in:"query"`
}

type MessageReadInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type MessageReadAllInp struct{}

type AdminMessageListInp struct {
	Page       int    `json:"page" d:"1" v:"min:1"`
	Size       int    `json:"size" d:"20" v:"min:1|max:100"`
	ReceiverId uint64 `json:"receiverId" in:"query"`
	IsRead     *bool  `json:"isRead" in:"query"`
}

type AdminMessageCreateInp struct {
	ReceiverIds []uint64 `json:"receiverIds" v:"required"`
	Title       string   `json:"title" v:"required|length:1,200"`
	Content     string   `json:"content" v:"required"`
	TargetType  string   `json:"targetType" v:"max-length:30"`
	TargetId    uint64   `json:"targetId"`
}

type MessageCreateInp struct {
	SenderId   uint64
	ReceiverId uint64
	Title      string
	Content    string
	TargetType string
	TargetId   uint64
}

type MessageNotifyInp struct {
	ActorId     uint64
	ReceiverId  uint64
	TitleKey    string
	TitleArgs   []any
	Content     string
	ContentKey  string
	ContentArgs []any
	TargetType  string
	TargetId    uint64
}
