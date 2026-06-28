package sitein

type SiteConfigListInp struct {
	Group string
}

type SiteConfigUpdateInp struct {
	Group string
	Key   string
	Value any
}
