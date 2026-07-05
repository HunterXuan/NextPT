package adminin

type IamUserListInp struct {
	Search string `json:"search"`
	Order  string `json:"order" d:"id desc"`
	Page   int    `json:"page" d:"1"`
	Size   int    `json:"size" d:"10"`
}

type IamUserUpdateInp struct {
	Id      uint64  `json:"id" in:"path" v:"required"`
	Status  *int    `json:"status" description:"0=pending 1=confirmed 2=disabled"`
	Role    *uint   `json:"role" description:"角色ID"`
	Passkey *string `json:"passkey"`
}

type IamUserStatDetailInp struct {
	Id uint64 `json:"id" in:"path" v:"required"`
}

type IamUserStatUpdateInp struct {
	Id             uint64   `json:"id" in:"path" v:"required"`
	UploadedDiff   *int64   `json:"uploadedDiff" description:"上传量增量(正加负减)"`
	DownloadedDiff *int64   `json:"downloadedDiff" description:"下载量增量(正加负减)"`
	BonusDiff      *float64 `json:"bonusDiff" description:"魔力值增量(正加负减)"`
}
