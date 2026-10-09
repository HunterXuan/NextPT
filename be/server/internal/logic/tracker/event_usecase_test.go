package tracker

import (
	"context"
	"errors"
	"testing"

	_ "server/internal/logic/mod"
	_ "server/internal/logic/sys"
	"server/internal/model/entity"
	"server/internal/model/in/modin"
	"server/internal/model/in/trackerin"
	"server/internal/service"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func TestRecordTrafficAnomaly(t *testing.T) {
	original := service.ModCheaterUsecase()
	defer service.RegisterModCheaterUsecase(original)
	for _, tt := range []struct {
		name, reason string
		torrent      *entity.CatalogTorrent
		err          error
		wantCalls    int
	}{
		{name: "normal"},
		{name: "anomaly", reason: "speed exceeded", torrent: &entity.CatalogTorrent{Seeders: 3, Leechers: 4}, wantCalls: 1},
		{name: "missing torrent", reason: "speed exceeded", wantCalls: 1},
		{name: "record failure is nonfatal", reason: "speed exceeded", err: errors.New("database unavailable"), wantCalls: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeCheaterRecorder{err: tt.err}
			service.RegisterModCheaterUsecase(fake)
			event := &trackerin.AnnounceEvent{UserId: 42, TorrentId: 7, Uploaded: 123, Downloaded: 456}
			NewTrackerEventUsecase().recordTrafficAnomaly(context.Background(), event, tt.torrent, 30, tt.reason)
			if fake.calls != tt.wantCalls {
				t.Fatalf("Record called %d times, want %d", fake.calls, tt.wantCalls)
			}
			if tt.wantCalls == 0 {
				return
			}
			if fake.in.UserId != 42 || fake.in.TorrentId != 7 || fake.in.Uploaded != 123 || fake.in.Downloaded != 456 || fake.in.AnnounceTime != 30 || fake.in.HitCount != 1 || fake.in.Comment != tt.reason {
				t.Fatalf("unexpected log: %#v", fake.in)
			}
			if tt.torrent != nil && (fake.in.Seeders != 3 || fake.in.Leechers != 4) {
				t.Fatal("missing swarm counts")
			}
		})
	}
}

type fakeCheaterRecorder struct {
	service.IModCheaterUsecase
	err   error
	calls int
	in    modin.RecordCheaterLogInp
}

func (f *fakeCheaterRecorder) Record(ctx context.Context, in modin.RecordCheaterLogInp) error {
	f.calls++
	f.in = in
	return f.err
}
