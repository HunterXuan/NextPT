package iam

import (
	"context"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/iamout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

type sIamRoleUsecase struct{}

const (
	iamRankSyncBatchSize = 1000
)

func init() {
	service.RegisterIamRoleUsecase(NewIamRoleUsecase())
}

func NewIamRoleUsecase() *sIamRoleUsecase {
	return &sIamRoleUsecase{}
}

func (s *sIamRoleUsecase) List(ctx context.Context, actor *model.Actor) (*iamout.RoleListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}

	roles, err := service.IamRoleDomain().ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]iamout.RoleItem, 0, len(roles))
	for _, role := range roles {
		items = append(items, iamout.RoleItem{
			Id:       role.Id,
			Level:    role.Level,
			Name:     s.localizeName(ctx, role),
			NameI18N: s.scanNameI18N(role),
			Rules:    s.scanRules(role),
			IsStaff:  role.IsStaff,
		})
	}

	return &iamout.RoleListOut{Roles: items}, nil
}

func (s *sIamRoleUsecase) SyncRanks(ctx context.Context) (*model.IamRankSyncResult, error) {
	roles, err := service.IamRoleDomain().ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	rankChain := model.NewIamRankChain(roles)
	if rankChain.Len() == 0 {
		return &model.IamRankSyncResult{}, nil
	}

	result := &model.IamRankSyncResult{}
	now := gtime.Now()
	systemActor := &model.Actor{}
	defaultLangCtx := gi18n.WithLanguage(ctx, g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String())
	for page := 1; ; page++ {
		candidates, err := service.IamUserDomain().QueryRankCandidates(ctx, page, iamRankSyncBatchSize)
		if err != nil {
			return nil, err
		}
		if len(candidates) == 0 {
			break
		}

		result.Scanned += len(candidates)
		for _, candidate := range candidates {
			currentNode := rankChain.NodeByRoleId(candidate.RoleId)
			if currentNode == nil {
				result.Skipped++
				continue
			}

			stats := model.NewIamRankStats(candidate, now)
			rankTarget := currentNode.ResolveTarget(stats)
			if rankTarget.ShouldBan {
				applied, err := s.applyRankAutoBan(ctx, systemActor, candidate.UserId)
				if err != nil {
					return nil, err
				}
				if !applied {
					result.Skipped++
					continue
				}
				result.Banned++
				continue
			}
			if rankTarget.Target == nil || rankTarget.Target.Role.Id == currentNode.Role.Id {
				result.Skipped++
				continue
			}

			if err := service.IamUserDomain().UpdateUserRole(ctx, candidate.UserId, rankTarget.Target.Role.Id); err != nil {
				return nil, err
			}
			service.IamUserUsecase().InvalidateUserCache(ctx, candidate.UserId)
			titleKey := "site.message.rank.demoted.title"
			if rankTarget.Target.Role.Level > currentNode.Role.Level {
				result.Promoted++
				titleKey = "site.message.rank.promoted.title"
			} else {
				result.Demoted++
			}
			service.SiteMessageUsecase().Notify(ctx, sitein.MessageNotifyInp{
				ReceiverId: candidate.UserId,
				TitleKey:   titleKey,
				ContentKey: "site.message.rank.changed.content",
				ContentArgs: []any{
					s.localizeName(defaultLangCtx, currentNode.Role),
					s.localizeName(defaultLangCtx, rankTarget.Target.Role),
				},
			})
		}

		if len(candidates) < iamRankSyncBatchSize {
			break
		}
	}

	if result.Promoted > 0 || result.Demoted > 0 || result.Banned > 0 {
		glog.Infof(ctx, "[Cron] IAM rank sync completed. Scanned=%d Promoted=%d Demoted=%d Banned=%d Skipped=%d",
			result.Scanned, result.Promoted, result.Demoted, result.Banned, result.Skipped)
	}
	return result, nil
}

func (s *sIamRoleUsecase) scanNameI18N(role entity.IamRole) map[string]string {
	names := make(map[string]string)
	if role.NameI18N == nil {
		return names
	}
	_ = role.NameI18N.Scan(&names)
	return names
}

func (s *sIamRoleUsecase) scanRules(role entity.IamRole) map[string]any {
	rules := make(map[string]any)
	if role.Rules == nil {
		return rules
	}
	_ = role.Rules.Scan(&rules)
	return rules
}

func (s *sIamRoleUsecase) localizeName(ctx context.Context, role entity.IamRole) string {
	names := s.scanNameI18N(role)
	if len(names) == 0 {
		return ""
	}

	lang := gi18n.LanguageFromCtx(ctx)
	if lang != "" && names[lang] != "" {
		return names[lang]
	}

	defaultLang := g.Cfg().MustGet(ctx, "i18n.default", "zh-CN").String()
	if defaultLang != "" && names[defaultLang] != "" {
		return names[defaultLang]
	}

	for _, name := range names {
		if name != "" {
			return name
		}
	}
	return ""
}

func (s *sIamRoleUsecase) applyRankAutoBan(ctx context.Context, actor *model.Actor, userId uint64) (bool, error) {
	banned, err := service.ModUserDomain().HasActiveMod(ctx, userId, []int{consts.ModUserTypeBanned})
	if err != nil {
		return false, err
	}
	if banned {
		return false, nil
	}

	err = service.ModUserUsecase().Apply(ctx, actor, modin.ApplyModInp{
		UserId:  userId,
		ModType: consts.ModUserTypeBanned,
		Reason:  consts.IamUserRankAutoBanReason,
	})
	return err == nil, err
}
