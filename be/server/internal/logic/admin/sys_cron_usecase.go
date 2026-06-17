package admin

import (
	"context"

	"server/internal/model"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtime"
)

type sAdminSysCronUsecase struct{}

func init() {
	service.RegisterAdminSysCronUsecase(NewAdminSysCronUsecase())
}

func NewAdminSysCronUsecase() *sAdminSysCronUsecase {
	return &sAdminSysCronUsecase{}
}

func (s *sAdminSysCronUsecase) List(ctx context.Context, actor *model.Actor) (*adminout.SysCronListOut, error) {
	entries := gcron.Entries()
	var list []*adminout.SysCronItem
	for _, e := range entries {
		list = append(list, &adminout.SysCronItem{
			Name:         e.Name,
			Status:       e.Status(),
			RegisterTime: gtime.New(e.RegisterTime),
		})
	}
	return &adminout.SysCronListOut{List: list}, nil
}

func (s *sAdminSysCronUsecase) LogList(ctx context.Context, actor *model.Actor, jobName string, page, size int) (*adminout.SysCronLogListOut, error) {
	logs, total, err := service.SysCron().AdminListCronLogs(ctx, jobName, page, size)
	if err != nil {
		return nil, err
	}

	var list []*adminout.SysCronLogItem
	for _, l := range logs {
		list = append(list, &adminout.SysCronLogItem{
			Id:           l.Id,
			JobName:      l.JobName,
			NodeIp:       l.NodeIp,
			Status:       l.Status,
			DurationMs:   l.DurationMs,
			ErrorMessage: l.ErrorMessage,
			CreatedAt:    l.CreatedAt,
			UpdatedAt:    l.UpdatedAt,
		})
	}

	return &adminout.SysCronLogListOut{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}
