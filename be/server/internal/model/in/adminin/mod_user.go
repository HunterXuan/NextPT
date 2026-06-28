package adminin

type ModUserListInp struct {
	Id uint64 `json:"id" v:"required" in:"path"`
}

type ModUserApplyInp struct {
	Id       uint64 `json:"id" v:"required" in:"path"`
	Type     uint   `json:"type" v:"required|in:1,2,3,4,5,6"`
	Reason   string `json:"reason" v:"required"`
	Duration int64  `json:"duration" v:"min:0" description:"0 means permanent, else duration in seconds"`
}

type ModUserRemoveInp struct {
	Id    uint64 `json:"id" v:"required" in:"path"`
	ModId uint64 `json:"modId" v:"required" in:"path"`
}
