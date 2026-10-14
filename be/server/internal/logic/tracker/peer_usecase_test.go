package tracker

import (
	"context"
	"testing"

	"server/internal/model"
	"server/internal/service"

	"github.com/gogf/gf/v2/os/gcache"
)

func TestClientWhitelistFailsClosedAndCachesRules(t *testing.T) {
	ctx := context.Background()
	key := service.SysCache().KeyTrackerClientWhitelist(ctx)
	original := service.TrackerPeerDomain()
	defer service.RegisterTrackerPeerDomain(original)
	defer gcache.Remove(ctx, key)
	fake := &fakeWhitelistPeerDomain{rules: []*model.TrackerClientWhitelist{
		{PeerIdPattern: "-BAD", AgentPattern: "["},
		{},
		{PeerIdPattern: "-qB", AgentPattern: "^qBittorrent/"},
	}}
	service.RegisterTrackerPeerDomain(fake)
	gcache.Remove(ctx, key)
	usecase := NewTrackerPeerUsecase()
	if usecase.CheckClientWhitelist(ctx, "-BAD123", "unknown") {
		t.Fatal("invalid regex must not allow a client")
	}
	if usecase.CheckClientWhitelist(ctx, "-qB4500-abcdefghijk", "Transmission/4.0") {
		t.Fatal("both fields must match")
	}
	if !usecase.CheckClientWhitelist(ctx, "-qB4500-abcdefghijk", "qBittorrent/4.5.0") {
		t.Fatal("valid client should match cached compiled regex")
	}
	if fake.calls != 1 {
		t.Fatalf("loaded rules %d times, want 1", fake.calls)
	}
}

func TestClientWhitelistExceptionDoesNotBlockOtherRules(t *testing.T) {
	ctx := context.Background()
	key := service.SysCache().KeyTrackerClientWhitelist(ctx)
	original := service.TrackerPeerDomain()
	defer service.RegisterTrackerPeerDomain(original)
	defer gcache.Remove(ctx, key)
	fake := &fakeWhitelistPeerDomain{rules: []*model.TrackerClientWhitelist{
		{PeerIdPattern: "^-UT2[0-9]{3}-", AgentPattern: "^uTorrent/2[0-9]{3}", Exceptions: []model.TrackerClientException{{Name: "specific version", PeerIdPattern: "^-UT2210-", Enabled: true}}},
	}}
	service.RegisterTrackerPeerDomain(fake)
	gcache.Remove(ctx, key)
	usecase := NewTrackerPeerUsecase()
	if usecase.CheckClientWhitelist(ctx, "-UT2210-abcdefghijk", "uTorrent/2210") {
		t.Fatal("own exception must exclude the client")
	}
	fake.rules = append(fake.rules, &model.TrackerClientWhitelist{PeerIdPattern: "^-UT2210-", AgentPattern: "^uTorrent/2210"})
	gcache.Remove(ctx, key)
	if !usecase.CheckClientWhitelist(ctx, "-UT2210-abcdefghijk", "uTorrent/2210") {
		t.Fatal("a separate allow rule must still be able to allow the client")
	}
}

type fakeWhitelistPeerDomain struct {
	service.ITrackerPeerDomain
	rules []*model.TrackerClientWhitelist
	calls int
}

func (f *fakeWhitelistPeerDomain) GetEnabledClientWhitelists(ctx context.Context) ([]*model.TrackerClientWhitelist, error) {
	f.calls++
	return f.rules, nil
}
