package catalog

import (
	"context"
	"math"
	"testing"

	"server/internal/consts"
)

func TestCatalogRequestNormalizeReward(t *testing.T) {
	usecase := NewCatalogRequestUsecase()
	tests := []struct {
		name    string
		amount  float64
		want    float64
		wantErr bool
	}{
		{name: "integer", amount: 50, want: 50},
		{name: "round to database precision", amount: 10.06, want: 10.1},
		{name: "zero", amount: 0, wantErr: true},
		{name: "negative", amount: -1, wantErr: true},
		{name: "nan", amount: math.NaN(), wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := usecase.normalizeRequestReward(context.Background(), test.amount)
			if test.wantErr {
				if err == nil {
					t.Fatalf("expected error, got amount %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeRequestReward() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("normalizeRequestReward() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestCatalogRequestClaimDuration(t *testing.T) {
	usecase := NewCatalogRequestUsecase()
	if got := usecase.requestClaimDuration(consts.CatalogRequestTypeTorrent); got != consts.CatalogRequestTorrentClaimDuration {
		t.Fatalf("torrent claim duration = %v, want %v", got, consts.CatalogRequestTorrentClaimDuration)
	}
	if got := usecase.requestClaimDuration(consts.CatalogRequestTypeReseed); got != consts.CatalogRequestReseedClaimDuration {
		t.Fatalf("reseed claim duration = %v, want %v", got, consts.CatalogRequestReseedClaimDuration)
	}
}
