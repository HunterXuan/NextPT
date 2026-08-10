package sitein

type TaskClaimInp struct {
	Key string `json:"key" in:"path" v:"required|length:3,64"`
}

type UserTaskRewardClaimInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}
