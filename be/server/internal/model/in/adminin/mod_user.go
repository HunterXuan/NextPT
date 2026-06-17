package adminin

type ModUserApplyInp struct {
	Id       uint64 `json:"id" v:"required" in:"path"`
	Type     uint   `json:"type" v:"required"`
	Reason   string `json:"reason"`
	Duration int64  `json:"duration" description:"0 means permanent, else duration in seconds"`
}

type ModUserRemoveInp struct {
	Id    uint64 `json:"id" v:"required" in:"path"`
	ModId uint64 `json:"modId" v:"required"`
}
