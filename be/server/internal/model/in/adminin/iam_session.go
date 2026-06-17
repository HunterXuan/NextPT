package adminin

type IamSessionDeleteInp struct {
	UserId uint64 `json:"userId" in:"path" v:"required"`
}
