package admin

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	_ "server/internal/logic/site"
	_ "server/internal/logic/sys"
	_ "server/internal/logic/tracker"
	"server/internal/model"
	"server/internal/model/in/adminin"
	"server/internal/model/in/sitein"
	"server/internal/service"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
)

func TestWhitelistErrorTranslations(t *testing.T) {
	path, err := filepath.Abs("../../../manifest/i18n")
	if err != nil {
		t.Fatal(err)
	}
	translator := gi18n.New(gi18n.Options{Path: path})
	for _, language := range []string{"zh-CN", "zh-TW", "en-US"} {
		if _, err := gjson.LoadPath(filepath.Join(path, language, "admin.yaml"), gjson.Options{}); err != nil {
			t.Fatalf("invalid %s admin translations: %v", language, err)
		}
		ctx := gi18n.WithLanguage(context.Background(), language)
		for _, key := range []string{"admin.tracker.client_not_found", "admin.tracker.client_rule_required", "admin.tracker.client_pattern_invalid"} {
			if translator.GetContent(ctx, key) == "" {
				t.Fatalf("missing %s translation for %s", language, key)
			}
		}
	}
}

func TestWhitelistValidationAndCheck(t *testing.T) {
	usecase := NewAdminTrackerClientWhitelistUsecase()
	ctx := context.Background()
	for _, rule := range []model.TrackerClientWhitelist{
		{Family: " "},
		{Family: "qBittorrent"},
		{Family: "qBittorrent", PeerIdPattern: "-qB", AgentPattern: "["},
	} {
		if err := usecase.validate(ctx, &rule); err == nil {
			t.Fatal("invalid rule accepted")
		}
	}
	rule := model.TrackerClientWhitelist{Family: " qBittorrent ", PeerIdPattern: " -qB ", AgentPattern: " ^qBittorrent/ "}
	if err := usecase.validate(ctx, &rule); err != nil {
		t.Fatal(err)
	}
	if rule.Family != "qBittorrent" || rule.PeerIdPattern != "-qB" {
		t.Fatal("fields not normalized")
	}
	result, err := usecase.Check(ctx, &model.Actor{Id: 1}, adminin.TrackerClientWhitelistCheckInp{PeerIdPattern: rule.PeerIdPattern, AgentPattern: rule.AgentPattern, PeerId: "-qB4500-abcdefghijk", UserAgent: "qBittorrent/4.5.0"})
	if err != nil || !result.Matched {
		t.Fatalf("Check failed: %v", err)
	}
	result, err = usecase.Check(ctx, &model.Actor{Id: 1}, adminin.TrackerClientWhitelistCheckInp{PeerIdPattern: "^-qB", PeerId: "-qB4500-abcdefghijk", Exceptions: []model.TrackerClientException{{Name: "blocked", PeerIdPattern: "^-qB4500-", Enabled: true}}})
	if err != nil || result.Matched || result.BlockedBy != "blocked" {
		t.Fatalf("exception result not reported: %#v %v", result, err)
	}
	rule.Exceptions = []model.TrackerClientException{{Name: "invalid"}}
	if err := usecase.validate(ctx, &rule); err == nil {
		t.Fatal("empty exception patterns accepted")
	}
}

func TestWhitelistMutationInvalidatesCacheAndAudits(t *testing.T) {
	ctx := context.Background()
	originalDomain, originalCache, originalAudit := service.TrackerPeerDomain(), service.SysCache(), service.SiteAuditUsecase()
	defer service.RegisterTrackerPeerDomain(originalDomain)
	defer service.RegisterSysCache(originalCache)
	defer service.RegisterSiteAuditUsecase(originalAudit)
	key := "nextpt:test:client-whitelist"
	defer gcache.Remove(ctx, key)
	for _, action := range []string{"create", "update", "delete", "failed update"} {
		t.Run(action, func(t *testing.T) {
			domain := &fakeWhitelistDomain{rule: &model.TrackerClientWhitelist{Id: 7, Family: "qBittorrent", PeerIdPattern: "-qB", AgentPattern: "^qBittorrent/", Enabled: true, Comment: "old note"}}
			if action == "failed update" {
				domain.err = errors.New("write failed")
			}
			cache, audit := &fakeWhitelistCache{key: key}, &fakeWhitelistAudit{}
			service.RegisterTrackerPeerDomain(domain)
			service.RegisterSysCache(cache)
			service.RegisterSiteAuditUsecase(audit)
			gcache.Set(ctx, key, "old rules", 0)
			usecase, actor := NewAdminTrackerClientWhitelistUsecase(), &model.Actor{Id: 42}
			var err error
			switch action {
			case "create":
				_, err = usecase.Create(ctx, actor, adminin.TrackerClientWhitelistCreateInp{TrackerClientWhitelistFields: adminin.TrackerClientWhitelistFields{Family: "qBittorrent", PeerIdPattern: "-qB", Enabled: true}})
			case "delete":
				err = usecase.Delete(ctx, actor, adminin.TrackerClientWhitelistDeleteInp{Id: 7})
			default:
				enabled := false
				err = usecase.Update(ctx, actor, adminin.TrackerClientWhitelistUpdateInp{Id: 7, Enabled: &enabled})
				if domain.written.Family != "qBittorrent" || domain.written.Comment != "old note" || domain.written.Enabled {
					t.Fatal("partial update lost fields or failed to disable rule")
				}
			}
			if action == "failed update" {
				if err == nil || cache.published || audit.calls != 0 {
					t.Fatal("failed write must not invalidate or audit")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			cached, _ := gcache.Get(ctx, key)
			if !cached.IsNil() || !cache.published || audit.calls != 1 || audit.in.Action != action || audit.in.TargetId != 7 {
				t.Fatal("cache invalidation or audit missing")
			}
		})
	}
}

type fakeWhitelistDomain struct {
	service.ITrackerPeerDomain
	rule    *model.TrackerClientWhitelist
	written model.TrackerClientWhitelist
	err     error
}

func (f *fakeWhitelistDomain) GetClientWhitelist(ctx context.Context, id uint) (*model.TrackerClientWhitelist, error) {
	return f.rule, nil
}
func (f *fakeWhitelistDomain) CreateClientWhitelist(ctx context.Context, rule model.TrackerClientWhitelist) (uint, error) {
	f.written = rule
	return 7, f.err
}
func (f *fakeWhitelistDomain) UpdateClientWhitelist(ctx context.Context, rule model.TrackerClientWhitelist) error {
	f.written = rule
	return f.err
}
func (f *fakeWhitelistDomain) DeleteClientWhitelist(ctx context.Context, id uint) error { return f.err }

type fakeWhitelistCache struct {
	service.ISysCache
	key       string
	published bool
}

func (f *fakeWhitelistCache) KeyTrackerClientWhitelist(ctx context.Context) string { return f.key }
func (f *fakeWhitelistCache) PublishInvalidate(ctx context.Context, key string) error {
	f.published = key == f.key
	return nil
}

type fakeWhitelistAudit struct {
	service.ISiteAuditUsecase
	calls int
	in    sitein.AuditRecordInp
}

func (f *fakeWhitelistAudit) Record(ctx context.Context, actor *model.Actor, in sitein.AuditRecordInp) {
	f.calls++
	f.in = in
}
