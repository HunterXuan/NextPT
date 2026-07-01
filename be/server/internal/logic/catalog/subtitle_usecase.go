package catalog

import (
	"context"

	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/database/gdb"

	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/catalogin"
	"server/internal/model/in/modin"
	"server/internal/model/out/catalogout"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type sCatalogSubtitleUsecase struct{}

func NewCatalogSubtitleUsecase() *sCatalogSubtitleUsecase {
	return &sCatalogSubtitleUsecase{}
}

func init() {
	service.RegisterCatalogSubtitleUsecase(NewCatalogSubtitleUsecase())
}

func (s *sCatalogSubtitleUsecase) formatSubtitles(ctx context.Context, actor *model.Actor, entities []entity.CatalogSubtitle) []catalogout.SubtitleListItem {
	var userIds []uint64
	for _, e := range entities {
		if s.canViewSubtitleOwner(actor, e) {
			userIds = append(userIds, e.UserId)
		}
	}

	uploaderMap := make(map[uint64]model.IamUserSummary)
	if len(userIds) > 0 {
		var users []entity.IamUser
		users, _ = service.IamUserDomain().GetUsersByIds(ctx, userIds)
		for _, u := range users {
			uploaderMap[u.Id] = model.IamUserSummary{
				Id:       u.Id,
				Username: u.Username,
			}
		}
	}

	var list []catalogout.SubtitleListItem
	for _, e := range entities {
		uploader := uploaderMap[e.UserId]
		if uploader.Id == 0 && s.canViewSubtitleOwner(actor, e) && e.UserId > 0 {
			uploader.Id = e.UserId
		}
		list = append(list, catalogout.SubtitleListItem{
			Id:        e.Id,
			TorrentId: e.TorrentId,
			Uploader:  uploader,
			Anonymous: e.Anonymous,
			FileName:  e.FileName,
			Language:  e.Language,
			Size:      e.FileSize,
			CreatedAt: e.CreatedAt.String(),
		})
	}
	return list
}

func (s *sCatalogSubtitleUsecase) canViewSubtitleOwner(actor *model.Actor, subtitle entity.CatalogSubtitle) bool {
	return !subtitle.Anonymous || (actor != nil && (actor.IsStaff || actor.Id == subtitle.UserId))
}

func (s *sCatalogSubtitleUsecase) List(ctx context.Context, actor *model.Actor, in catalogin.SubtitleListInp) (*catalogout.SubtitleListOut, error) {
	subs, total, err := service.CatalogSubtitleDomain().QuerySubtitles(ctx, 0, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	return &catalogout.SubtitleListOut{
		List:  s.formatSubtitles(ctx, actor, subs),
		Total: total,
	}, nil
}

func (s *sCatalogSubtitleUsecase) ListByTorrent(ctx context.Context, actor *model.Actor, in catalogin.TorrentSubtitleListInp) (*catalogout.SubtitleListOut, error) {
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return nil, err
	}

	subs, total, err := service.CatalogSubtitleDomain().QuerySubtitles(ctx, in.Id, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	return &catalogout.SubtitleListOut{
		List:  s.formatSubtitles(ctx, actor, subs),
		Total: total,
	}, nil
}

func (s *sCatalogSubtitleUsecase) Upload(ctx context.Context, actor *model.Actor, in catalogin.SubtitleUploadInp) (uint64, error) {
	if actor == nil {
		return 0, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	// Verify target torrent is visible to the actor
	_, err := service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, in.Id)
	if err != nil {
		return 0, err
	}

	// 1. Read file
	f, err := in.File.Open()
	if err != nil {
		return 0, gerror.Wrap(err, gi18n.T(ctx, "catalog.subtitle.open_upload_failed"))
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return 0, gerror.Wrap(err, gi18n.T(ctx, "catalog.subtitle.read_upload_failed"))
	}
	ext := filepath.Ext(in.File.Filename)

	var subtitleId uint64
	// 2. Transaction 1: Insert DB Record
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := service.CatalogSubtitleDomain().InsertSubtitle(ctx, in.Id, actor.Id, in.File.Filename, ext, len(data), in.Language, in.Anonymous)
		if err != nil {
			return err
		}
		subtitleId = id
		return nil
	})
	if err != nil {
		return 0, err
	}

	// 3. S3 Upload (Outside DB Transaction)
	s3Key := fmt.Sprintf("subs/%d/%d%s", in.Id, subtitleId, ext)
	err = service.SysStorage().Upload(ctx, s3Key, data, "application/octet-stream")
	if err != nil {
		// Rollback compensation
		_ = service.CatalogSubtitleDomain().DeleteSubtitle(ctx, subtitleId)
		return 0, gerror.Wrap(err, gi18n.T(ctx, "catalog.subtitle.upload_s3_failed"))
	}

	// 4. Transaction 2: Update status & Grant permissions
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = service.CatalogSubtitleDomain().UpdateSubtitleStoragePath(ctx, subtitleId, s3Key)
		if err != nil {
			return err
		}

		_ = service.IamPermissionDomain().GrantUserPermission(ctx, actor.Id, fmt.Sprintf("update:catalog/subtitle:%d", subtitleId), false)
		return nil
	})

	if err != nil {
		return 0, err
	}
	return subtitleId, nil
}

func (s *sCatalogSubtitleUsecase) Download(ctx context.Context, actor *model.Actor, in catalogin.SubtitleDownloadInp) (*catalogout.SubtitleDownloadOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	sub, err := service.CatalogSubtitleDomain().GetSubtitleById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// Verify target torrent is visible to the actor
	_, err = service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, sub.TorrentId)
	if err != nil {
		return nil, err
	}

	targetPath := service.SysStorage().GetLocalPath(ctx, sub.StoragePath)

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.subtitle.mkdir_failed"))
		}
		err = service.SysStorage().Download(ctx, sub.StoragePath, targetPath)
		if err != nil {
			return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.subtitle.download_s3_failed"))
		}
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "catalog.subtitle.read_cache_failed"))
	}

	_ = service.CatalogSubtitleDomain().IncrementDownloadCount(ctx, in.Id)

	return &catalogout.SubtitleDownloadOut{
		Bytes:    data,
		FileName: sub.FileName,
		MimeType: "application/octet-stream",
	}, nil
}

func (s *sCatalogSubtitleUsecase) Update(ctx context.Context, actor *model.Actor, in catalogin.SubtitleUpdateInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	sub, err := service.CatalogSubtitleDomain().GetSubtitleById(ctx, in.Id)
	if err != nil {
		return err
	}

	// Verify target torrent is visible to the actor
	_, err = service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, sub.TorrentId)
	if err != nil {
		return err
	}

	return service.CatalogSubtitleDomain().UpdateSubtitle(ctx, in.Id, in.Language, actor.Id, false)
}

func (s *sCatalogSubtitleUsecase) Report(ctx context.Context, actor *model.Actor, in catalogin.SubtitleReportInp) error {
	if actor == nil {
		return gerror.New(gi18n.T(ctx, "catalog.general.unauthorized"))
	}

	sub, err := service.CatalogSubtitleDomain().GetSubtitleById(ctx, in.Id)
	if err != nil {
		return err
	}

	// Verify target torrent is visible to the actor
	_, err = service.CatalogTorrentDomain().LoadVisibleTorrent(ctx, actor, sub.TorrentId)
	if err != nil {
		return err
	}

	return service.ModReportUsecase().Create(ctx, actor, modin.CreateReportInp{
		TargetType: "subtitle",
		TargetId:   in.Id,
		Reason:     in.Reason,
	})
}
