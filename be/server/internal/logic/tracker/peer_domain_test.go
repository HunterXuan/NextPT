package tracker

import (
	"context"
	"strings"
	"testing"
	"time"

	"server/internal/model/entity"
	"server/internal/model/in/trackerin"

	"github.com/gogf/gf/v2/os/gtime"
)

func TestCalculateTrafficDiff(t *testing.T) {
	const mib = int64(1024 * 1024)
	now := gtime.Now()
	for _, tt := range []struct {
		name                 string
		peer                 *entity.TrackerPeer
		uploaded, downloaded int64
		wantUp, wantDn       int64
		wantAnomaly          bool
	}{
		{name: "first report", uploaded: 100 * mib},
		{name: "normal", peer: &entity.TrackerPeer{Uploaded: uint64(mib), Downloaded: uint64(mib), LastAction: now.Add(-10 * time.Second)}, uploaded: 3 * mib, downloaded: 4 * mib, wantUp: 2 * mib, wantDn: 3 * mib},
		{name: "at limit", peer: &entity.TrackerPeer{LastAction: now.Add(-10 * time.Second)}, uploaded: 60 * mib, downloaded: 40 * mib, wantUp: 60 * mib, wantDn: 40 * mib},
		{name: "combined over limit", peer: &entity.TrackerPeer{LastAction: now.Add(-10 * time.Second)}, uploaded: 60 * mib, downloaded: 50 * mib, wantAnomaly: true},
		{name: "same second small", peer: &entity.TrackerPeer{LastAction: now}, uploaded: mib, wantUp: mib},
		{name: "same second large", peer: &entity.TrackerPeer{LastAction: now}, uploaded: 11 * mib, wantAnomaly: true},
		{name: "counter reset", peer: &entity.TrackerPeer{Uploaded: uint64(100 * mib), Downloaded: uint64(100 * mib), LastAction: now.Add(-10 * time.Second)}, uploaded: mib, downloaded: 2 * mib, wantUp: mib, wantDn: 2 * mib},
		{name: "missing prior time", peer: &entity.TrackerPeer{}, uploaded: mib, wantUp: mib},
	} {
		t.Run(tt.name, func(t *testing.T) {
			up, dn, reason := NewTrackerPeerDomain().CalculateTrafficDiff(context.Background(), &trackerin.AnnounceEvent{Uploaded: tt.uploaded, Downloaded: tt.downloaded, Now: now}, tt.peer)
			if up != tt.wantUp || dn != tt.wantDn || (reason != "") != tt.wantAnomaly {
				t.Fatalf("got (%d, %d, %q), want (%d, %d, anomaly=%v)", up, dn, reason, tt.wantUp, tt.wantDn, tt.wantAnomaly)
			}
			if tt.wantAnomaly && !strings.Contains(reason, "uploaded_delta=") {
				t.Fatal("anomaly reason missing discarded delta")
			}
		})
	}
}
