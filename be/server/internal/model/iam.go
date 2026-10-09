package model

import (
	"fmt"
	"html"
	"math"
	"sort"
	"strings"
	"time"

	"server/internal/model/entity"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type IamSession struct {
	Id         string `json:"id"`
	UserId     uint64 `json:"userId"`
	DeviceHash string `json:"deviceHash,omitempty"`
	Ip         string `json:"ip"`
	UserAgent  string `json:"userAgent"`
	CreatedAt  int64  `json:"createdAt"`
	LastSeenAt int64  `json:"lastSeenAt"`
}

type IamSessionItem struct {
	Id         string      `json:"id"`
	Ip         string      `json:"ip"`
	UserAgent  string      `json:"userAgent"`
	CreatedAt  *gtime.Time `json:"createdAt"`
	LastSeenAt *gtime.Time `json:"lastSeenAt"`
	Current    bool        `json:"current"`
}

func (s IamSession) Item(currentSessionId string) IamSessionItem {
	return IamSessionItem{
		Id:         s.Id,
		Ip:         s.Ip,
		UserAgent:  s.UserAgent,
		CreatedAt:  gtime.NewFromTime(time.UnixMilli(s.CreatedAt)),
		LastSeenAt: gtime.NewFromTime(time.UnixMilli(s.LastSeenAt)),
		Current:    s.Id == currentSessionId,
	}
}

const (
	iamRankBytesPerGiB = 1024 * 1024 * 1024
)

const (
	iamRankRulePromotion = "promotion"
	iamRankRuleDemotion  = "demotion"
)

const (
	iamRankFieldAccountAgeDaysGte = "accountAgeDaysGte"
	iamRankFieldUploadedGiBGte    = "uploadedGiBGte"
	iamRankFieldUploadedGiBGt     = "uploadedGiBGt"
	iamRankFieldUploadedGiBLte    = "uploadedGiBLte"
	iamRankFieldDownloadedGiBGte  = "downloadedGiBGte"
	iamRankFieldDownloadedGiBGt   = "downloadedGiBGt"
	iamRankFieldDownloadedGiBLte  = "downloadedGiBLte"
	iamRankFieldRatioGt           = "ratioGt"
	iamRankFieldRatioGte          = "ratioGte"
	iamRankFieldRatioLt           = "ratioLt"
	iamRankFieldRatioLte          = "ratioLte"
)

type IamUserSummary struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type IamAccountActionMail struct {
	Kind      string
	Recipient string
	Subject   string
	Greeting  string
	Intro     string
	ActionURL string
	Action    string
	Expiry    string
	Note      string
}

func (m IamAccountActionMail) Bodies() (textBody, htmlBody string) {
	textBody = strings.Join([]string{m.Greeting, m.Intro, m.Action, m.ActionURL, m.Expiry, m.Note}, "\n\n")
	htmlBody = fmt.Sprintf(
		`<!doctype html><html><body style="margin:0;background:#f8fafc;font-family:Arial,sans-serif;color:#0f172a"><div style="max-width:560px;margin:0 auto;padding:32px 20px"><div style="background:#fff;border:1px solid #e2e8f0;border-radius:8px;padding:28px"><h1 style="margin:0 0 20px;font-size:20px">%s</h1><p style="margin:0 0 12px;line-height:1.7">%s</p><p style="margin:0 0 20px;line-height:1.7">%s</p><p style="margin:0 0 20px"><a href="%s" style="display:inline-block;border-radius:6px;background:#0284c7;padding:11px 18px;color:#fff;text-decoration:none;font-weight:600">%s</a></p><p style="margin:0 0 8px;color:#475569;line-height:1.7">%s</p><p style="margin:0;color:#64748b;line-height:1.7">%s</p></div></div></body></html>`,
		html.EscapeString(m.Subject), html.EscapeString(m.Greeting), html.EscapeString(m.Intro),
		html.EscapeString(m.ActionURL), html.EscapeString(m.Action), html.EscapeString(m.Expiry), html.EscapeString(m.Note),
	)
	return
}

type IamUserPermissionListOptions struct {
	SourceType   *int
	WildcardOnly bool
	Page         int
	Size         int
}

type IamLoginLogListOptions struct {
	UserId uint64
	Result *int
	Page   int
	Size   int
}

type IamLoginLogItem struct {
	Id         uint64      `json:"id"`
	UserId     uint64      `json:"userId"`
	Ip         string      `json:"ip"`
	UserAgent  string      `json:"userAgent"`
	Result     int         `json:"result"`
	FailReason string      `json:"failReason"`
	CreatedAt  *gtime.Time `json:"createdAt"`
}

func NewIamLoginLogItem(log entity.IamLoginLog) IamLoginLogItem {
	return IamLoginLogItem{
		Id:         log.Id,
		UserId:     log.UserId,
		Ip:         log.Ip,
		UserAgent:  log.UserAgent,
		Result:     log.Result,
		FailReason: log.FailReason,
		CreatedAt:  log.CreatedAt,
	}
}

type IamRankCandidate struct {
	UserId     uint64
	RoleId     uint
	Uploaded   uint64
	Downloaded uint64
	CreatedAt  *gtime.Time
}

type IamRankSyncResult struct {
	Scanned  int
	Promoted int
	Demoted  int
	Banned   int
	Skipped  int
}

type IamRankResolveResult struct {
	Target    *IamRankNode
	ShouldBan bool
}

type IamRankStats struct {
	accountAgeDays float64
	uploadedGiB    float64
	downloadedGiB  float64
	ratio          float64
}

type IamRankChain struct {
	nodes    []IamRankNode
	nodeById map[uint]*IamRankNode
}

type IamRankNode struct {
	Role entity.IamRole
	Prev *IamRankNode
	Next *IamRankNode
}

func NewIamRankStats(candidate IamRankCandidate, now *gtime.Time) IamRankStats {
	stats := IamRankStats{
		uploadedGiB:   float64(candidate.Uploaded) / iamRankBytesPerGiB,
		downloadedGiB: float64(candidate.Downloaded) / iamRankBytesPerGiB,
	}
	if candidate.Downloaded == 0 {
		if candidate.Uploaded > 0 {
			stats.ratio = math.Inf(1)
		}
	} else {
		stats.ratio = float64(candidate.Uploaded) / float64(candidate.Downloaded)
	}
	if candidate.CreatedAt != nil && now != nil {
		stats.accountAgeDays = now.Time.Sub(candidate.CreatedAt.Time).Hours() / 24
	}
	return stats
}

func NewIamRankChain(roles []entity.IamRole) *IamRankChain {
	rankRoles := make([]entity.IamRole, 0, len(roles))
	for _, role := range roles {
		if role.Id == 0 || role.IsStaff {
			continue
		}
		rankRoles = append(rankRoles, role)
	}
	sort.SliceStable(rankRoles, func(i, j int) bool {
		if rankRoles[i].Level == rankRoles[j].Level {
			return rankRoles[i].Id < rankRoles[j].Id
		}
		return rankRoles[i].Level < rankRoles[j].Level
	})

	chain := &IamRankChain{
		nodes:    make([]IamRankNode, len(rankRoles)),
		nodeById: make(map[uint]*IamRankNode, len(rankRoles)),
	}
	for i, role := range rankRoles {
		chain.nodes[i].Role = role
	}
	for i := range chain.nodes {
		node := &chain.nodes[i]
		if i > 0 {
			node.Prev = &chain.nodes[i-1]
		}
		if i < len(chain.nodes)-1 {
			node.Next = &chain.nodes[i+1]
		}
		chain.nodeById[node.Role.Id] = node
	}
	return chain
}

func (c *IamRankChain) Len() int {
	if c == nil {
		return 0
	}
	return len(c.nodes)
}

func (c *IamRankChain) NodeByRoleId(roleId uint) *IamRankNode {
	if c == nil {
		return nil
	}
	return c.nodeById[roleId]
}

func (n *IamRankNode) ResolveTarget(stats IamRankStats) IamRankResolveResult {
	if n == nil {
		return IamRankResolveResult{}
	}
	if n.MatchDemotion(stats) {
		if n.Prev == nil {
			return IamRankResolveResult{ShouldBan: true}
		}
		return IamRankResolveResult{Target: n.DemotionTarget(stats)}
	}
	return IamRankResolveResult{Target: n.PromotionTarget(stats)}
}

func (n *IamRankNode) PromotionTarget(stats IamRankStats) *IamRankNode {
	var target *IamRankNode
	for node := n.Next; node != nil; node = node.Next {
		if node.MatchPromotion(stats) {
			target = node
		}
	}
	return target
}

func (n *IamRankNode) DemotionTarget(stats IamRankStats) *IamRankNode {
	var lowest *IamRankNode
	for node := n.Prev; node != nil; node = node.Prev {
		if node.MatchPromotion(stats) {
			return node
		}
		lowest = node
	}
	return lowest
}

func (n *IamRankNode) MatchPromotion(stats IamRankStats) bool {
	if n == nil {
		return false
	}
	return matchIamRankRule(n.Role.Rules, iamRankRulePromotion, stats)
}

func (n *IamRankNode) MatchDemotion(stats IamRankStats) bool {
	if n == nil {
		return false
	}
	return matchIamRankRule(n.Role.Rules, iamRankRuleDemotion, stats)
}

func matchIamRankRule(rulesJSON *gjson.Json, ruleType string, stats IamRankStats) bool {
	if rulesJSON == nil || ruleType == "" {
		return false
	}

	groups := rulesJSON.Get(ruleType).Maps()
	if len(groups) == 0 {
		return false
	}

	for _, group := range groups {
		if matchIamRankRuleGroup(group, stats) {
			return true
		}
	}
	return false
}

func matchIamRankRuleGroup(group map[string]any, stats IamRankStats) bool {
	if len(group) == 0 {
		return false
	}

	for key, raw := range group {
		value := gconv.Float64(raw)
		switch key {
		case iamRankFieldAccountAgeDaysGte:
			if stats.accountAgeDays < value {
				return false
			}
		case iamRankFieldUploadedGiBGte:
			if stats.uploadedGiB < value {
				return false
			}
		case iamRankFieldUploadedGiBGt:
			if stats.uploadedGiB <= value {
				return false
			}
		case iamRankFieldUploadedGiBLte:
			if stats.uploadedGiB > value {
				return false
			}
		case iamRankFieldDownloadedGiBGte:
			if stats.downloadedGiB < value {
				return false
			}
		case iamRankFieldDownloadedGiBGt:
			if stats.downloadedGiB <= value {
				return false
			}
		case iamRankFieldDownloadedGiBLte:
			if stats.downloadedGiB > value {
				return false
			}
		case iamRankFieldRatioGt:
			if stats.ratio <= value {
				return false
			}
		case iamRankFieldRatioGte:
			if stats.ratio < value {
				return false
			}
		case iamRankFieldRatioLt:
			if stats.ratio >= value {
				return false
			}
		case iamRankFieldRatioLte:
			if stats.ratio > value {
				return false
			}
		default:
			return false
		}
	}
	return true
}
