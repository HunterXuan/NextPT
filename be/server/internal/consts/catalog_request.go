package consts

import "time"

const (
	CatalogRequestTypeTorrent = 1
	CatalogRequestTypeReseed  = 2
)

const (
	CatalogRequestStatusOpen      = 0
	CatalogRequestStatusClaimed   = 1
	CatalogRequestStatusSubmitted = 2
	CatalogRequestStatusCompleted = 3
	CatalogRequestStatusCancelled = 4
)

const (
	CatalogRequestViewAll     = "all"
	CatalogRequestViewCreated = "created"
	CatalogRequestViewClaimed = "claimed"
)

const (
	CatalogRequestTorrentClaimDuration = 72 * time.Hour
	CatalogRequestReseedClaimDuration  = 24 * time.Hour
)
