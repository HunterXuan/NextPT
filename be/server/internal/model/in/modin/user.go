package modin

type ApplyModInp struct {
	UserId          uint64 `json:"user_id" v:"required" in:"path"`
	ModType         int    `json:"mod_type" v:"required|in:1,2,3,4,5,6"`
	Reason          string `json:"reason" v:"required"`
	DurationDays    int    `json:"duration_days" description:"0 = Permanent"`
	DurationSeconds int64  `json:"duration_seconds" description:"0 = Permanent"`
}

type RemoveModInp struct {
	UserId uint64 `json:"user_id" v:"required" in:"path"`
	Id     uint64 `json:"id" v:"required"`
}

type ListUserInp struct {
	UserId uint64 `json:"user_id" v:"required"`
}
