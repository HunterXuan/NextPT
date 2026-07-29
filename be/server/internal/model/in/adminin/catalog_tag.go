package adminin

type CatalogTagGroupCreateInp struct {
	NameI18N    map[string]any `json:"nameI18N" v:"required" description:"多语言名称映射"`
	Slug        string         `json:"slug" v:"required|max-length:50" description:"稳定分组标识"`
	CategoryIds []uint         `json:"categoryIds" description:"适用分类 ID，空表示全部"`
	SortOrder   int            `json:"sortOrder" description:"排序值"`
}

type CatalogTagGroupUpdateInp struct {
	Id          uint           `json:"id" in:"path" v:"required" description:"标签组 ID"`
	NameI18N    map[string]any `json:"nameI18N" v:"required" description:"多语言名称映射"`
	CategoryIds []uint         `json:"categoryIds" description:"适用分类 ID，空表示全部"`
	SortOrder   int            `json:"sortOrder" description:"排序值"`
}

type CatalogTagGroupDeleteInp struct {
	Id uint `json:"id" in:"path" v:"required" description:"标签组 ID"`
}

type CatalogTagCreateInp struct {
	GroupId   uint           `json:"groupId" in:"path" v:"required" description:"标签组 ID"`
	NameI18N  map[string]any `json:"nameI18N" v:"required" description:"多语言名称映射"`
	Value     string         `json:"value" v:"required|max-length:50" description:"稳定标签值"`
	SortOrder int            `json:"sortOrder" description:"排序值"`
}

type CatalogTagUpdateInp struct {
	Id        uint           `json:"id" in:"path" v:"required" description:"标签 ID"`
	NameI18N  map[string]any `json:"nameI18N" v:"required" description:"多语言名称映射"`
	SortOrder int            `json:"sortOrder" description:"排序值"`
}

type CatalogTagDeleteInp struct {
	Id uint `json:"id" in:"path" v:"required" description:"标签 ID"`
}

type CatalogTagGroupListInp struct{}
