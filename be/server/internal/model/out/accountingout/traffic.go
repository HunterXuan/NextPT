package accountingout

type TrafficGetMeOut struct {
	Uploaded   uint64  `json:"uploaded" description:"总上传量 (bytes)"`
	Downloaded uint64  `json:"downloaded" description:"总下载量 (bytes)"`
	ShareRatio float64 `json:"shareRatio" description:"分享率"`
	SeedTime   uint64  `json:"seedTime" description:"总做种时间 (秒)"`
	LeechTime  uint64  `json:"leechTime" description:"总下载时间 (秒)"`
}

type TrafficHistoryItem struct {
	Date       string `json:"date" description:"日期 / 月份"`
	Uploaded   uint64 `json:"uploaded"`
	Downloaded uint64 `json:"downloaded"`
	SeedTime   uint64 `json:"seedTime"`
	LeechTime  uint64 `json:"leechTime"`
	Bonus      string `json:"bonus" description:"获得魔力值"`
}

type TrafficHistoryListOut struct {
	List []TrafficHistoryItem `json:"list"`
}
