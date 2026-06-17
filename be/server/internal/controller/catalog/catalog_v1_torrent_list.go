package catalog

import (
	"context"

	v1 "server/api/catalog/v1"
	"server/internal/library/contexts"
	"server/internal/model/in/catalogin"
	"server/internal/model/out/catalogout"
	"server/internal/service"
)

func (c *ControllerV1) TorrentList(ctx context.Context, req *v1.TorrentListReq) (res *v1.TorrentListRes, err error) {
	in := catalogin.TorrentListInp{
		Page:       req.Page,
		Size:       req.Size,
		CategoryId: req.CategoryId,
		Type:       req.Type,
	}

	out, err := service.CatalogTorrentUsecase().List(ctx, contexts.GetActor(ctx), in)
	if err != nil {
		return nil, err
	}

	var list []catalogout.TorrentListItem
	for _, item := range out.List {
		list = append(list, catalogout.TorrentListItem{
			Id:         item.Id,
			Name:       item.Name,
			SubTitle:   item.SubTitle,
			CategoryId: item.CategoryId,
			Size:       item.Size,
			FileCount:  item.FileCount,
			Type:       item.Type,
			Seeders:    item.Seeders,
			Leechers:   item.Leechers,
			Snatched:   item.Snatched,
			OwnerId:    item.OwnerId,
			Anonymous:  item.Anonymous,
			CreatedAt:  item.CreatedAt,
		})
	}

	return &v1.TorrentListRes{
		TorrentListOut: catalogout.TorrentListOut{
			List:  list,
			Total: out.Total,
		},
	}, nil
}
