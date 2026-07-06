package site

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/sitein"
	"server/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type sSiteAuditUsecase struct{}

func init() {
	service.RegisterSiteAuditUsecase(NewSiteAuditUsecase())
}

func NewSiteAuditUsecase() *sSiteAuditUsecase {
	return &sSiteAuditUsecase{}
}

func (s *sSiteAuditUsecase) Record(ctx context.Context, actor *model.Actor, in sitein.AuditRecordInp) {
	if in.Action == "" || in.TargetType == "" {
		return
	}
	userId := in.UserId
	if userId == 0 && actor != nil {
		userId = actor.Id
	}
	if err := service.SiteAuditDomain().Create(ctx, sitein.AuditCreateInp{
		UserId:     userId,
		Action:     in.Action,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		Detail:     s.detailString(in.Detail),
		Ip:         s.requestIp(ctx),
		Level:      in.Level,
		CreatedAt:  gtime.Now(),
	}); err != nil {
		g.Log().Warningf(ctx, "record site audit failed: %v", err)
	}
}

func (s *sSiteAuditUsecase) requestIp(ctx context.Context) string {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return r.GetClientIp()
}

func (s *sSiteAuditUsecase) detailString(detail any) string {
	switch v := detail.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	default:
		data, err := gjson.Encode(v)
		if err != nil {
			return ""
		}
		return string(data)
	}
}
