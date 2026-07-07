package accountingout

type TrafficGetMeOut struct {
	Uploaded      uint64  `json:"uploaded" description:"总入账上传量 (bytes)"`
	Downloaded    uint64  `json:"downloaded" description:"总入账下载量 (bytes)"`
	RawUploaded   uint64  `json:"rawUploaded" description:"真实总上传量 (bytes)"`
	RawDownloaded uint64  `json:"rawDownloaded" description:"真实总下载量 (bytes)"`
	ShareRatio    float64 `json:"shareRatio" description:"分享率"`
	SeedTime      uint64  `json:"seedTime" description:"总做种时间 (秒)"`
	LeechTime     uint64  `json:"leechTime" description:"总下载时间 (秒)"`
}

type TrafficHistoryItem struct {
	Date          string `json:"date" description:"日期 / 月份"`
	Uploaded      uint64 `json:"uploaded" description:"入账上传量 (bytes)"`
	Downloaded    uint64 `json:"downloaded" description:"入账下载量 (bytes)"`
	RawUploaded   uint64 `json:"rawUploaded" description:"真实上传量 (bytes)"`
	RawDownloaded uint64 `json:"rawDownloaded" description:"真实下载量 (bytes)"`
	SeedTime      uint64 `json:"seedTime"`
	LeechTime     uint64 `json:"leechTime"`
	Bonus         string `json:"bonus" description:"获得魔力值"`
}

type TrafficHistoryListOut struct {
	List []TrafficHistoryItem `json:"list"`
}
