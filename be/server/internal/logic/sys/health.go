package sys

import (
	"context"
	"sync"
	"time"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	sysHealthCacheTTL     = 3 * time.Second
	sysHealthCheckTimeout = 2 * time.Second
)

type healthProbe func(ctx context.Context) error

type sSysHealth struct {
	mu            sync.Mutex
	now           func() time.Time
	checkDatabase healthProbe
	checkRedis    healthProbe
	checkedAt     time.Time
	cached        model.SysHealthReadiness
}

func init() {
	service.RegisterSysHealth(NewHealth())
}

func NewHealth() *sSysHealth {
	return newHealth(
		func(ctx context.Context) error {
			_, err := g.DB().GetValue(ctx, "SELECT 1")
			return err
		},
		func(ctx context.Context) error {
			_, err := g.Redis().Do(ctx, "PING")
			return err
		},
	)
}

func newHealth(checkDatabase, checkRedis healthProbe) *sSysHealth {
	return &sSysHealth{
		now:           time.Now,
		checkDatabase: checkDatabase,
		checkRedis:    checkRedis,
	}
}

func (s *sSysHealth) Readiness(ctx context.Context) model.SysHealthReadiness {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	if !s.checkedAt.IsZero() && now.Sub(s.checkedAt) < sysHealthCacheTTL {
		return s.cached
	}

	checks := map[string]model.SysHealthCheck{
		"database": s.check(ctx, s.checkDatabase),
		"redis":    s.check(ctx, s.checkRedis),
	}
	status := consts.SysHealthStatusOK
	for _, check := range checks {
		if check.Status != consts.SysHealthStatusOK {
			status = consts.SysHealthStatusUnavailable
			break
		}
	}

	s.cached = model.SysHealthReadiness{
		Status:    status,
		CheckedAt: now,
		Checks:    checks,
	}
	s.checkedAt = now
	return s.cached
}

func (s *sSysHealth) check(ctx context.Context, probe healthProbe) model.SysHealthCheck {
	checkCtx, cancel := context.WithTimeout(ctx, sysHealthCheckTimeout)
	defer cancel()

	startedAt := time.Now()
	status := consts.SysHealthStatusOK
	if probe == nil || probe(checkCtx) != nil {
		status = consts.SysHealthStatusUnavailable
	}
	return model.SysHealthCheck{
		Status:    status,
		LatencyMs: time.Since(startedAt).Milliseconds(),
	}
}
