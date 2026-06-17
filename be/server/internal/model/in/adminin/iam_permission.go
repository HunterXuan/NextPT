package adminin

type IamPermissionListInp struct {
	RoleId uint `json:"roleId" v:"required"`
}

type IamUserPermissionGrantInp struct {
	UserId  uint64 `json:"userId" in:"path" v:"required"`
	PermKey string `json:"permKey" v:"required"`
	IsDeny  bool   `json:"isDeny"`
}

type IamUserPermissionRevokeInp struct {
	UserId  uint64 `json:"userId" in:"path" v:"required"`
	PermKey string `json:"permKey" v:"required"`
	IsDeny  bool   `json:"isDeny"`
}
