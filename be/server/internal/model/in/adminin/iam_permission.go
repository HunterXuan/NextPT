package adminin

type IamPermissionListInp struct {
	RoleId uint `json:"roleId" v:"required"`
}

type IamUserPermissionDetailInp struct {
	UserId       uint64 `json:"userId" in:"path" v:"required"`
	SourceType   *int   `json:"sourceType" in:"query"`
	WildcardOnly bool   `json:"wildcardOnly" in:"query"`
	Page         int    `json:"page" in:"query" d:"1"`
	Size         int    `json:"size" in:"query" d:"20"`
}

type IamUserPermissionGrantInp struct {
	UserId   uint64   `json:"userId" in:"path" v:"required"`
	PermKeys []string `json:"permKeys" v:"required"`
	IsDeny   bool     `json:"isDeny"`
}

type IamUserPermissionRevokeInp struct {
	UserId uint64   `json:"userId" in:"path" v:"required"`
	Ids    []uint64 `json:"ids" v:"required"`
}
