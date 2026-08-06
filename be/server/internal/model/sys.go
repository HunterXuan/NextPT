package model

import (
	"time"

	"server/internal/consts"
)

type SysHealthReadiness struct {
	Status    string                    `json:"status"`
	CheckedAt time.Time                 `json:"checkedAt"`
	Checks    map[string]SysHealthCheck `json:"checks"`
}

type SysHealthCheck struct {
	Status    string `json:"status"`
	LatencyMs int64  `json:"latencyMs"`
}

func (h SysHealthReadiness) IsReady() bool {
	return h.Status == consts.SysHealthStatusOK
}
