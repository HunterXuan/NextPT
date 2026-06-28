package adminin

type SiteConfigListInp struct {
	Group string `json:"group" in:"path" v:"required#{#admin.config.group_req}" description:"配置组"`
}

type SiteConfigUpdateInp struct {
	Group string `json:"group" in:"path" v:"required#{#admin.config.group_req}"`
	Key   string `json:"key" in:"path" v:"required#{#admin.config.key_req}"`
	Value any    `json:"value" description:"配置值"`
}
