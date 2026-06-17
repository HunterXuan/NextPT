package adminin

type ForumNodeCreateInp struct {
	CategoryId    uint     `json:"categoryId" v:"required"`
	Slug          string   `json:"slug" v:"required"`
	NameI18N      string   `json:"nameI18N" v:"required"`
	DescI18N      string   `json:"descI18N" v:"required"`
	SortOrder     int      `json:"sortOrder"`
	MinRoleRead   int      `json:"minRoleRead"`
	MinRoleWrite  int      `json:"minRoleWrite"`
	MinRoleCreate int      `json:"minRoleCreate"`
	Moderators    []uint64 `json:"moderators"`
}

type ForumNodeUpdateInp struct {
	Id            uint      `json:"id" in:"path" v:"required"`
	CategoryId    *uint     `json:"categoryId"`
	Slug          *string   `json:"slug"`
	NameI18N      *string   `json:"nameI18N"`
	DescI18N      *string   `json:"descI18N"`
	SortOrder     *int      `json:"sortOrder"`
	MinRoleRead   *int      `json:"minRoleRead"`
	MinRoleWrite  *int      `json:"minRoleWrite"`
	MinRoleCreate *int      `json:"minRoleCreate"`
	Moderators    *[]uint64 `json:"moderators"`
}

type ForumNodeDeleteInp struct {
	Id uint `json:"id" in:"path" v:"required"`
}

type ForumNodeListInp struct {
}
