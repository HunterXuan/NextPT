package tracker

import (
	"context"

	v1 "server/api/tracker/v1"
	"server/internal/library/contexts"
	"server/internal/service"

	"github.com/anacrolix/torrent/bencode"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) Announce(ctx context.Context, req *v1.AnnounceReq) (res *v1.AnnounceRes, err error) {
	r := ghttp.RequestFromCtx(ctx)

	// URL query parameter is raw binary string, but struct binding might fail or mutate it in some edge cases.
	// To be perfectly safe, we directly extract info_hash and peer_id from the query string.
	req.InfoHash = r.GetQuery("info_hash").String()
	req.PeerId = r.GetQuery("peer_id").String()

	out, err := service.TrackerPeerUsecase().Announce(ctx, contexts.GetActor(ctx), req)
	if err != nil {
		r.Response.Write(bencodeError(err.Error()))
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
