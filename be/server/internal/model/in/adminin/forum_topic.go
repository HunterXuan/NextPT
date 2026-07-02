package adminin

type ForumTopicLockInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type ForumTopicUnlockInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type ForumTopicPinInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type ForumTopicUnpinInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type ForumTopicMoveInp struct {
	Id     uint64 `json:"id" in:"path" v:"required"`
	NodeId uint   `json:"nodeId" v:"required"`
}

type ForumTopicDeleteInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}
