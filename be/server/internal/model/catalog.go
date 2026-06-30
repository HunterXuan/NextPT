package model

type CatalogUploadConfig struct {
	Title  CatalogUploadTitleConfig   `json:"title"`
	Fields []CatalogUploadFieldConfig `json:"fields"`
}

type CatalogUploadTitleConfig struct {
	Mode                string                   `json:"mode"`
	AllowManualOverride bool                     `json:"allowManualOverride"`
	Parts               []CatalogUploadTitlePart `json:"parts"`
}

type CatalogUploadTitlePart struct {
	Field     string `json:"field"`
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Separator string `json:"separator"`
}

type CatalogUploadFieldConfig struct {
	Key         string                     `json:"key"`
	Type        string                     `json:"type"`
	Label       map[string]any             `json:"label"`
	Description map[string]any             `json:"description,omitempty"`
	Placeholder map[string]any             `json:"placeholder,omitempty"`
	Required    bool                       `json:"required"`
	Options     *CatalogUploadFieldOptions `json:"options,omitempty"`
}

type CatalogUploadFieldOptions struct {
	Source string                    `json:"source"`
	Slug   string                    `json:"slug,omitempty"`
	Items  []CatalogUploadOptionItem `json:"items,omitempty"`
}

type CatalogUploadOptionItem struct {
	Value string         `json:"value"`
	Label map[string]any `json:"label"`
}
