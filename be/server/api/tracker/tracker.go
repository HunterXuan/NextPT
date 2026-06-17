// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package tracker

import (
	"context"

	"server/api/tracker/v1"
)

type ITrackerV1 interface {
	Announce(ctx context.Context, req *v1.AnnounceReq) (res *v1.AnnounceRes, err error)
	Download(ctx context.Context, req *v1.DownloadReq) (res *v1.DownloadRes, err error)
	Scrape(ctx context.Context, req *v1.ScrapeReq) (res *v1.ScrapeRes, err error)
}
