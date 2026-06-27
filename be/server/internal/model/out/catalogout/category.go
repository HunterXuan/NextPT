package catalogout

type CategoryItem struct {
	Id   uint           `json:"id"`
	Name map[string]any `json:"name"`
	Slug string         `json:"slug"`
}

type CategoryListOut struct {
	List []CategoryItem `json:"list"`
}

type TagItem struct {
	Id   uint           `json:"id"`
	Name map[string]any `json:"name"`
}

type TagGroupItem struct {
	Id         uint           `json:"id"`
	Name       map[string]any `json:"name"`
	Slug       string         `json:"slug"`
	Categories []uint         `json:"categories"`
	Tags       []TagItem      `json:"tags"`
}

type TagGroupListOut struct {
	List []TagGroupItem `json:"list"`
}
