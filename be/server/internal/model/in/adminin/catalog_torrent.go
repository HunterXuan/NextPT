package adminin

type CatalogTorrentDeleteInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}
