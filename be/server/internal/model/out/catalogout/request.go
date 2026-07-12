package catalogout

import "server/internal/model"

type RequestListItem struct {
	Id              uint64                `json:"id"`
	RequestType     uint                  `json:"requestType"`
	CategoryId      uint                  `json:"categoryId"`
	TargetTorrentId uint64                `json:"targetTorrentId"`
	ResultTorrentId uint64                `json:"resultTorrentId"`
	Title           string                `json:"title"`
	RewardAmount    float64               `json:"rewardAmount"`
	Status          uint                  `json:"status"`
	Requester       model.IamUserSummary  `json:"requester"`
	Claimer         *model.IamUserSummary `json:"claimer"`
	ClaimExpiresAt  string                `json:"claimExpiresAt"`
	SubmittedAt     string                `json:"submittedAt"`
	CompletedAt     string                `json:"completedAt"`
	CreatedAt       string                `json:"createdAt"`
	UpdatedAt       string                `json:"updatedAt"`
}

type RequestListOut struct {
	List  []RequestListItem `json:"list"`
	Total int               `json:"total"`
}

type RequestDetailOut struct {
	RequestListItem
	Description   string                       `json:"description"`
	TargetTorrent *model.CatalogTorrentSummary `json:"targetTorrent"`
	ResultTorrent *model.CatalogTorrentSummary `json:"resultTorrent"`
	CancelReason  string                       `json:"cancelReason"`
	Actions       model.CatalogRequestActions  `json:"actions"`
}

type RequestCreateOut struct {
	Id uint64 `json:"id"`
}
