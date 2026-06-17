package adminin

type ModCheaterListInp struct {
	Page   int `json:"page" d:"1"`
	Size   int `json:"size" d:"20"`
	Status int `json:"status" d:"-1"`
}

type ModCheaterResolveInp struct {
	Id      uint64 `json:"id" v:"required" in:"path"`
	Status  int    `json:"status" v:"required"`
	Comment string `json:"comment"`
}
