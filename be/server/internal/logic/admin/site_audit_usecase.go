package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/out/adminout"
	"server/internal/service"
)

type sAdminSiteAuditUsecase struct{}

func init() {
	service.RegisterAdminSiteAuditUsecase(NewAdminSiteAuditUsecase())
}

func NewAdminSiteAuditUsecase() *sAdminSiteAuditUsecase {
	return &sAdminSiteAuditUsecase{}
}

func (s *sAdminSiteAuditUsecase) List(ctx context.Context, actor *model.Actor, in adminin.SiteAuditListInp) (*adminout.SiteAuditListOut, error) {
	logs, total, err := service.SiteAuditDomain().AdminListAudits(ctx, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	var list []*adminout.SiteAuditItem
	for _, l := range logs {
		list = append(list, &adminout.SiteAuditItem{
			Id:         l.Id,
			UserId:     l.UserId,
			Action:     l.Action,
			TargetType: l.TargetType,
			TargetId:   l.TargetId,
			Level:      l.Level,
			Ip:         l.Ip,
			Detail:     l.Detail,
			CreatedAt:  l.CreatedAt,
		})
	}

	return &adminout.SiteAuditListOut{
		List:  list,
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}
