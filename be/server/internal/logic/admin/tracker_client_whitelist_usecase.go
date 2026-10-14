package admin

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"server/internal/consts"
	libtracker "server/internal/library/tracker"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/model/out/adminout"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
)

type sAdminTrackerClientWhitelistUsecase struct{}

func NewAdminTrackerClientWhitelistUsecase() *sAdminTrackerClientWhitelistUsecase {
	return &sAdminTrackerClientWhitelistUsecase{}
}
func init() {
	service.RegisterAdminTrackerClientWhitelistUsecase(NewAdminTrackerClientWhitelistUsecase())
}

func (s *sAdminTrackerClientWhitelistUsecase) List(ctx context.Context, actor *model.Actor) (*adminout.TrackerClientWhitelistListOut, error) {
	list, err := service.TrackerPeerDomain().AdminListClientWhitelists(ctx)
	if err != nil {
		return nil, err
	}
	return &adminout.TrackerClientWhitelistListOut{List: list}, nil
}

func (s *sAdminTrackerClientWhitelistUsecase) Create(ctx context.Context, actor *model.Actor, in adminin.TrackerClientWhitelistCreateInp) (*adminout.TrackerClientWhitelistCreateOut, error) {
	rule := model.TrackerClientWhitelist{Family: in.Family, PeerIdPattern: in.PeerIdPattern, AgentPattern: in.AgentPattern, Exceptions: in.Exceptions, Enabled: in.Enabled, Comment: in.Comment}
	if err := s.validate(ctx, &rule); err != nil {
		return nil, err
	}
	id, err := service.TrackerPeerDomain().CreateClientWhitelist(ctx, rule)
	if err != nil {
		return nil, err
	}
	rule.Id = id
	s.changed(ctx, actor, consts.SiteAuditActionCreate, rule)
	return &adminout.TrackerClientWhitelistCreateOut{Id: id}, nil
}

func (s *sAdminTrackerClientWhitelistUsecase) Update(ctx context.Context, actor *model.Actor, in adminin.TrackerClientWhitelistUpdateInp) error {
	rule, err := s.get(ctx, in.Id)
	if err != nil {
		return err
	}
	if in.Family != nil {
		rule.Family = *in.Family
	}
	if in.PeerIdPattern != nil {
		rule.PeerIdPattern = *in.PeerIdPattern
	}
	if in.AgentPattern != nil {
		rule.AgentPattern = *in.AgentPattern
	}
	if in.Exceptions != nil {
		rule.Exceptions = *in.Exceptions
	}
	if in.Enabled != nil {
		rule.Enabled = *in.Enabled
	}
	if in.Comment != nil {
		rule.Comment = *in.Comment
	}
	if err := s.validate(ctx, rule); err != nil {
		return err
	}
	if err := service.TrackerPeerDomain().UpdateClientWhitelist(ctx, *rule); err != nil {
		return err
	}
	s.changed(ctx, actor, consts.SiteAuditActionUpdate, *rule)
	return nil
}

func (s *sAdminTrackerClientWhitelistUsecase) Delete(ctx context.Context, actor *model.Actor, in adminin.TrackerClientWhitelistDeleteInp) error {
	rule, err := s.get(ctx, in.Id)
	if err != nil {
		return err
	}
	if err := service.TrackerPeerDomain().DeleteClientWhitelist(ctx, in.Id); err != nil {
		return err
	}
	s.changed(ctx, actor, consts.SiteAuditActionDelete, *rule)
	return nil
}

func (s *sAdminTrackerClientWhitelistUsecase) Check(ctx context.Context, actor *model.Actor, in adminin.TrackerClientWhitelistCheckInp) (*adminout.TrackerClientWhitelistCheckOut, error) {
	config := model.TrackerClientWhitelist{Family: "test", PeerIdPattern: in.PeerIdPattern, AgentPattern: in.AgentPattern, Exceptions: in.Exceptions}
	if err := s.validate(ctx, &config); err != nil {
		return nil, err
	}
	rule, err := s.compile(ctx, config)
	if err != nil {
		return nil, err
	}
	matched, blockedBy := rule.Check(in.PeerId, in.UserAgent)
	return &adminout.TrackerClientWhitelistCheckOut{Matched: matched, BlockedBy: blockedBy}, nil
}

func (s *sAdminTrackerClientWhitelistUsecase) get(ctx context.Context, id uint) (*model.TrackerClientWhitelist, error) {
	rule, err := service.TrackerPeerDomain().GetClientWhitelist(ctx, id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, gerror.New(gi18n.T(ctx, "admin.tracker.client_not_found"))
	}
	return rule, nil
}

func (s *sAdminTrackerClientWhitelistUsecase) validate(ctx context.Context, rule *model.TrackerClientWhitelist) error {
	rule.Family, rule.PeerIdPattern = strings.TrimSpace(rule.Family), strings.TrimSpace(rule.PeerIdPattern)
	rule.AgentPattern, rule.Comment = strings.TrimSpace(rule.AgentPattern), strings.TrimSpace(rule.Comment)
	if rule.Family == "" {
		return gerror.New(gi18n.T(ctx, "admin.tracker.client_family_required"))
	}
	if rule.Exceptions == nil {
		rule.Exceptions = []model.TrackerClientException{}
	}
	if len(rule.Exceptions) > 50 {
		return gerror.New(gi18n.T(ctx, "admin.tracker.client_exceptions_invalid"))
	}
	for i := range rule.Exceptions {
		e := &rule.Exceptions[i]
		e.Name, e.Comment = strings.TrimSpace(e.Name), strings.TrimSpace(e.Comment)
		e.PeerIdPattern, e.AgentPattern = strings.TrimSpace(e.PeerIdPattern), strings.TrimSpace(e.AgentPattern)
		if e.Name == "" || utf8.RuneCountInString(e.Name) > 60 || utf8.RuneCountInString(e.Comment) > 200 || utf8.RuneCountInString(e.PeerIdPattern) > 200 || utf8.RuneCountInString(e.AgentPattern) > 200 {
			return gerror.New(gi18n.T(ctx, "admin.tracker.client_exceptions_invalid"))
		}
	}
	_, err := s.compile(ctx, *rule)
	return err
}

func (s *sAdminTrackerClientWhitelistUsecase) compile(ctx context.Context, config model.TrackerClientWhitelist) (libtracker.ClientWhitelist, error) {
	rule, err := libtracker.CompileClientWhitelist(config.PeerIdPattern, config.AgentPattern, config.Exceptions)
	if err != nil {
		key := "admin.tracker.client_pattern_invalid"
		if errors.Is(err, libtracker.ErrEmptyClientRule) {
			key = "admin.tracker.client_rule_required"
		}
		return rule, gerror.New(gi18n.T(ctx, key))
	}
	return rule, nil
}

func (s *sAdminTrackerClientWhitelistUsecase) changed(ctx context.Context, actor *model.Actor, action string, rule model.TrackerClientWhitelist) {
	key := service.SysCache().KeyTrackerClientWhitelist(ctx)
	if _, err := gcache.Remove(ctx, key); err != nil {
		glog.Warningf(ctx, "clear client whitelist cache failed: %v", err)
	}
	if err := service.SysCache().PublishInvalidate(ctx, key); err != nil {
		glog.Warningf(ctx, "publish client whitelist cache invalidation failed: %v", err)
	}
	service.SiteAuditUsecase().Record(ctx, actor, sitein.AuditRecordInp{
		Action: action, TargetType: consts.SiteAuditTargetTypeTrackerClientWhitelist, TargetId: uint64(rule.Id), Level: consts.SiteAuditLevelCritical,
		Detail: map[string]any{"family": rule.Family, "peerIdPattern": rule.PeerIdPattern, "agentPattern": rule.AgentPattern, "exceptions": rule.Exceptions, "enabled": rule.Enabled, "comment": rule.Comment},
	})
}
