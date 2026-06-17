package catalogin

type CommentCreateInp struct {
	Id      uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Content string `json:"content" v:"required|length:3,10000#{#catalog.comment.content_req}|{#catalog.comment.content_len}" description:"评论内容"`
}

type CommentListInp struct {
	Id   uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Page int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int    `json:"size" d:"20" v:"max:100" description:"每页数量"`
}

type CommentReportInp struct {
	Id     uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Cid    uint64 `json:"cid" in:"path" v:"required#{#catalog.comment.cid_req}" description:"评论ID"`
	Reason string `json:"reason" v:"required|length:5,500#{#catalog.comment.report_reason_req}|{#catalog.comment.report_reason_len}" description:"举报原因"`
}

type CommentToggleLikeInp struct {
	Id  uint64 `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Cid uint64 `json:"cid" in:"path" v:"required#{#catalog.comment.cid_req}" description:"评论ID"`
}

type CommentRewardInp struct {
	Id     uint64  `json:"id" in:"path" v:"required#{#catalog.torrent.id_req}" description:"种子ID"`
	Cid    uint64  `json:"cid" in:"path" v:"required#{#catalog.comment.cid_req}" description:"评论ID"`
	Amount float64 `json:"amount" v:"required|min:1#{#catalog.reward.amount_req}|{#catalog.reward.amount_min}" description:"赞赏金额"`
}
