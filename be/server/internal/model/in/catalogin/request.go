package catalogin

type RequestListInp struct {
	Page        int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size        int    `json:"size" d:"20" v:"min:1|max:100" description:"每页数量"`
	Keyword     string `json:"keyword" description:"标题关键词"`
	RequestType uint   `json:"requestType" v:"in:0,1,2" description:"请求类型: 0=全部 1=求种 2=续种"`
	Status      *uint  `json:"status" v:"in:0,1,2,3,4" description:"请求状态"`
	CategoryId  uint   `json:"categoryId" description:"分类ID"`
	View        string `json:"view" d:"all" v:"in:all,created,claimed" description:"列表范围"`
}

type RequestGetInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
}

type RequestCreateInp struct {
	RequestType     uint    `json:"requestType" v:"required|in:1,2#{#catalog.request.type_req}|{#catalog.request.type_invalid}" description:"请求类型"`
	CategoryId      uint    `json:"categoryId" description:"求种分类ID"`
	TargetTorrentId uint64  `json:"targetTorrentId" description:"续种目标种子ID"`
	Title           string  `json:"title" v:"length:0,500#{#catalog.request.title_invalid}" description:"求种标题"`
	Description     string  `json:"description" v:"required|length:3,10000#{#catalog.request.description_req}|{#catalog.request.description_invalid}" description:"请求说明"`
	RewardAmount    float64 `json:"rewardAmount" v:"required|min:0.1#{#catalog.request.reward_req}|{#catalog.request.reward_invalid}" description:"魔力奖励"`
}

type RequestClaimInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
}

type RequestAbandonInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
}

type RequestSubmitInp struct {
	Id              uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
	ResultTorrentId uint64 `json:"resultTorrentId" description:"求种结果种子ID"`
}

type RequestCompleteInp struct {
	Id uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
}

type RequestCancelInp struct {
	Id     uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
	Reason string `json:"reason" v:"length:0,500#{#catalog.request.cancel_reason_invalid}" description:"取消原因"`
}

type RequestCommentCreateInp struct {
	Id      uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
	Content string `json:"content" v:"required|length:3,10000#{#catalog.comment.content_req}|{#catalog.comment.content_len}" description:"评论内容"`
}

type RequestCommentListInp struct {
	Id   uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
	Page int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int    `json:"size" d:"20" v:"max:100" description:"每页数量"`
}

type RequestCommentActionInp struct {
	Id  uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
	Cid uint64 `json:"cid" in:"path" v:"required#{#catalog.comment.cid_req}" description:"评论ID"`
}

type RequestCommentRewardInp struct {
	Id     uint64  `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
	Cid    uint64  `json:"cid" in:"path" v:"required#{#catalog.comment.cid_req}" description:"评论ID"`
	Amount float64 `json:"amount" v:"required|min:1#{#catalog.reward.amount_req}|{#catalog.reward.amount_min}" description:"赞赏金额"`
}

type RequestCommentReportInp struct {
	Id     uint64 `json:"id" in:"path" v:"required#{#catalog.request.id_req}" description:"请求ID"`
	Cid    uint64 `json:"cid" in:"path" v:"required#{#catalog.comment.cid_req}" description:"评论ID"`
	Reason string `json:"reason" v:"required|length:5,500#{#catalog.comment.report_reason_req}|{#catalog.comment.report_reason_len}" description:"举报原因"`
}
