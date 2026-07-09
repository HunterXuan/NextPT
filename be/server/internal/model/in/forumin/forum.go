package forumin

type NodeListInp struct {
	Scope string `json:"scope" d:"read" v:"in:read,create"`
}

type TopicListInp struct {
	Slug string `json:"slug" in:"path"` // Path param for routing
	Page int    `json:"page" d:"1"`
	Size int    `json:"size" d:"20"`
}

type TopicGetHotInp struct {
	Size int `json:"size" d:"5" v:"min:1|max:10"`
}

type TopicCreateInp struct {
	NodeId  uint   `json:"nodeId" v:"required"`
	Subject string `json:"subject" v:"required|length:2,200"`
	Content string `json:"content" v:"required|min-length:2"`
}

type TopicUpdateInp struct {
	Id      uint64 `json:"id" in:"path" v:"required"`
	NodeId  uint   `json:"nodeId" v:"required"`
	Subject string `json:"subject" v:"required|length:2,200"`
	Content string `json:"content" v:"required|min-length:2"`
}

type TopicAppendInp struct {
	Id      uint64 `json:"id" in:"path" v:"required"`
	Content string `json:"content" v:"required"`
}

type TopicDetailInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type ReplyListInp struct {
	Id   uint64 `json:"id" in:"path" v:"required"`
	Page int    `json:"page" d:"1"`
	Size int    `json:"size" d:"50"`
}

type ReplyCreateInp struct {
	Id      uint64 `json:"id" in:"path" v:"required"`
	ReplyTo uint64 `json:"replyTo" d:"0"`
	Content string `json:"content" v:"required|min-length:2"`
}

type TopicToggleLikeInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type TopicRewardInp struct {
	Id     uint64  `json:"id" in:"path" v:"required"`
	Amount float64 `json:"amount" v:"required|min:1"`
}

type TopicRewardListInp struct {
	Id   uint64 `json:"id" in:"path" v:"required"`
	Page int    `json:"page" d:"1" v:"min:1"`
	Size int    `json:"size" d:"20" v:"min:1|max:100"`
}

type TopicReportInp struct {
	Id     uint64 `json:"id" in:"path" v:"required"`
	Reason string `json:"reason" v:"required|max-length:500"`
}

type TopicBookmarkInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type TopicUnbookmarkInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type TopicBookmarkListInp struct {
	Page int `json:"page" d:"1"`
	Size int `json:"size" d:"20"`
}

type ReplyToggleLikeInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type ReplyRewardInp struct {
	Id     uint64  `json:"id" in:"path" v:"required"`
	Amount float64 `json:"amount" v:"required|min:1"`
}

type ReplyReportInp struct {
	Id     uint64 `json:"id" in:"path" v:"required"`
	Reason string `json:"reason" v:"required|max-length:500"`
}
