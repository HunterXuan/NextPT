package adminin

type ModReportListInp struct {
	Page       int    `json:"page" d:"1"`
	Size       int    `json:"size" d:"20"`
	Status     int    `json:"status" d:"-1"`
	TargetType string `json:"targetType"`
}

type ModReportResolveInp struct {
	Id      uint64 `json:"id" v:"required" in:"path"`
	Status  int    `json:"status" v:"required"`
	Comment string `json:"comment"`
}
