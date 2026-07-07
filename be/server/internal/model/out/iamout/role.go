package iamout

type RoleItem struct {
	Id       uint              `json:"id"`
	Level    int               `json:"level"`
	Name     string            `json:"name"`
	NameI18N map[string]string `json:"nameI18N"`
	Rules    map[string]any    `json:"rules"`
	IsStaff  bool              `json:"isStaff"`
}

type RoleListOut struct {
	Roles []RoleItem `json:"roles"`
}
