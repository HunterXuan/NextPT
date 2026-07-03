package accountingin

type PeerListInp struct {
	Page   int    `json:"page" d:"1" v:"min:1" description:"页码"`
	Size   int    `json:"size" d:"20" v:"min:1|max:100" description:"每页数量"`
	Status string `json:"status" in:"query" d:"all" v:"in:all,seeding,leeching" description:"状态: all/seeding/leeching"`
}
