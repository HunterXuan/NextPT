package modin

type CreateReportInp struct {
	TargetType string `json:"target_type" v:"required"`
	TargetId   uint64 `json:"target_id" v:"required"`
	Reason     string `json:"reason" v:"required"`
}

type ListReportsInp struct {
	Page       int    `json:"page" d:"1" v:"min:1"`
	Size       int    `json:"size" d:"20" v:"min:1|max:100"`
	Status     int    `json:"status" d:"-1" description:"-1=All, 0=Pending, 1=Resolved, 2=Rejected"`
	TargetType string `json:"target_type"`
}

type ResolveReportInp struct {
	Id      uint64 `json:"id" v:"required" in:"path"`
	Status  int    `json:"status" v:"required|in:1,2"`
	Comment string `json:"comment"`
}
