package tracker

import (
	"context"

	v1 "server/api/tracker/v1"
	"server/internal/library/contexts"
	"server/internal/service"

	"github.com/anacrolix/torrent/bencode"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) Scrape(ctx context.Context, req *v1.ScrapeReq) (res *v1.ScrapeRes, err error) {
	r := ghttp.RequestFromCtx(ctx)

	// info_hash is raw bytes, we should parse it from URL query directly to support multiple
	infoHashes := r.GetQuery("info_hash").Strings()

	out, err := service.TrackerPeerUsecase().Scrape(ctx, contexts.GetActor(ctx), infoHashes)
	if err != nil {
		r.Response.Write(bencodeError("Database error"))
		r.ExitAll()
		return
	}

	b, err := bencode.Marshal(out)
	if err != nil {
		r.Response.Write(bencodeError("Encode error"))
		r.ExitAll()
		return
	}

	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.Response.Write(b)
	r.ExitAll()

	return nil, nil
}

func bencodeError(reason string) []byte {
	m := map[string]string{"failure reason": reason}
	b, _ := bencode.Marshal(m)
	return b
}
