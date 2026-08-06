package sys

import (
	"context"
	"errors"
	"testing"
	"time"

	"server/internal/consts"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func TestReadinessCachesHealthyChecks(t *testing.T) {
	var databaseCalls, redisCalls int
	health := newHealth(
		func(context.Context) error {
			databaseCalls++
			return nil
		},
		func(context.Context) error {
			redisCalls++
			return nil
		},
	)
	now := time.Date(2026, 8, 6, 0, 0, 0, 0, time.UTC)
	health.now = func() time.Time { return now }

	if readiness := health.Readiness(context.Background()); !readiness.IsReady() {
		t.Fatalf("Readiness() status = %q, want ready", readiness.Status)
	}
	if readiness := health.Readiness(context.Background()); !readiness.IsReady() {
		t.Fatalf("cached Readiness() status = %q, want ready", readiness.Status)
	}
	if databaseCalls != 1 || redisCalls != 1 {
		t.Fatalf("probe calls = database:%d redis:%d, want 1 each", databaseCalls, redisCalls)
	}

	now = now.Add(sysHealthCacheTTL)
	_ = health.Readiness(context.Background())
	if databaseCalls != 2 || redisCalls != 2 {
		t.Fatalf("expired cache calls = database:%d redis:%d, want 2 each", databaseCalls, redisCalls)
	}
}

func TestReadinessReportsUnavailableDependency(t *testing.T) {
	health := newHealth(
		func(context.Context) error { return errors.New("database unavailable") },
		func(context.Context) error { return nil },
	)
	readiness := health.Readiness(context.Background())
	if readiness.Status != consts.SysHealthStatusUnavailable {
		t.Fatalf("Readiness() status = %q, want %q", readiness.Status, consts.SysHealthStatusUnavailable)
	}
	if readiness.Checks["database"].Status != consts.SysHealthStatusUnavailable {
		t.Fatalf("database check = %q, want unavailable", readiness.Checks["database"].Status)
	}
	if readiness.Checks["redis"].Status != consts.SysHealthStatusOK {
		t.Fatalf("redis check = %q, want ok", readiness.Checks["redis"].Status)
	}
}
