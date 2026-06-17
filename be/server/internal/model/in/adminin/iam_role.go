package adminin

type IamRoleListInp struct {
}

type IamRoleCreateInp struct {
	Level       int            `json:"level" v:"required"`
	NameI18N    map[string]any `json:"nameI18N" v:"required"`
	Rules       map[string]any `json:"rules"`
	Permissions []string       `json:"permissions"`
	IsStaff     bool           `json:"isStaff"`
}

type IamRoleUpdateInp struct {
	Id          uint           `json:"id" in:"path" v:"required"`
	Level       *int           `json:"level"`
	NameI18N    map[string]any `json:"nameI18N"`
	Rules       map[string]any `json:"rules"`
	Permissions []string       `json:"permissions"`
	IsStaff     *bool          `json:"isStaff"`
}

type IamRoleDeleteInp struct {
	Id uint `json:"id" in:"path" v:"required"`
}
