package adminin

type CatalogCategoryCreateInp struct {
	NameI18N  map[string]any `json:"nameI18N" v:"required" description:"多语言名称映射"`
	Slug      string         `json:"slug" v:"required" description:"URL-friendly 标识"`
	SortOrder int            `json:"sortOrder" description:"排序值"`
	Enabled   *bool          `json:"enabled" description:"是否启用，默认启用"`
}

type CatalogCategoryUpdateInp struct {
	Id        uint           `json:"id" in:"path" v:"required" description:"分类ID"`
	NameI18N  map[string]any `json:"nameI18N" description:"多语言名称映射"`
	Slug      *string        `json:"slug" description:"URL-friendly 标识"`
	SortOrder *int           `json:"sortOrder" description:"排序值"`
	Enabled   *bool          `json:"enabled" description:"是否启用"`
}

type CatalogCategoryDeleteInp struct {
	Id uint `json:"id" in:"path" v:"required" description:"分类ID"`
}

type CatalogCategoryListInp struct {
}
