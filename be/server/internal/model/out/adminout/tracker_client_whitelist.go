package adminout

import "server/internal/model"

type TrackerClientWhitelistListOut struct {
	List []model.TrackerClientWhitelist `json:"list"`
}

type TrackerClientWhitelistCreateOut struct {
	Id uint `json:"id"`
}

type TrackerClientWhitelistCheckOut struct {
	Matched   bool   `json:"matched"`
	BlockedBy string `json:"blockedBy"`
}
