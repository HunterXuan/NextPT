package adminin

type SysCronLogListInp struct {
	Name   string `json:"name" v:"required" in:"path" dc:"任务名称"`
	Page   int    `json:"page" d:"1"`
	Size   int    `json:"size" d:"20"`
	Status int    `json:"status" d:"-1"`
}
