package adminin

type ForumCategoryCreateInp struct {
	NameI18N    string `json:"nameI18N" v:"required" description:"多语言名称映射(JSON string)"`
	DescI18N    string `json:"descI18N" v:"required" description:"多语言描述映射(JSON string)"`
	SortOrder   int    `json:"sortOrder"`
	MinRoleView int    `json:"minRoleView"`
}

type ForumCategoryUpdateInp struct {
	Id          uint    `json:"id" in:"path" v:"required"`
	NameI18N    *string `json:"nameI18N"`
	DescI18N    *string `json:"descI18N"`
	SortOrder   *int    `json:"sortOrder"`
	MinRoleView *int    `json:"minRoleView"`
}

type ForumCategoryDeleteInp struct {
	Id uint `json:"id" in:"path" v:"required"`
}

type ForumCategoryListInp struct {
}
