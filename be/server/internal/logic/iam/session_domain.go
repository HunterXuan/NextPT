package iam

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"server/internal/model"
	"server/internal/service"

	"github.com/goflyfox/gtoken/v2/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const iamSessionLastSeenInterval = 5 * time.Minute
const iamKnownDeviceTTL = 365 * 24 * time.Hour

const iamRegisterDeviceScript = `
local added = redis.call("SADD", KEYS[1], ARGV[1])
redis.call("EXPIRE", KEYS[1], tonumber(ARGV[2]))
return added
`

type sIamSessionDomain struct {
	gfToken    gtoken.Token
	middleware gtoken.Middleware
	timeout    time.Duration
}

func init() {
	service.RegisterIamSessionDomain(NewIamSessionDomain())
}

func NewIamSessionDomain() *sIamSessionDomain {
	ctx := context.TODO()
	cacheMode := g.Cfg().MustGet(ctx, "iam.cacheMode", gtoken.CacheModeCache).Int8()
	cachePreKey := g.Cfg().MustGet(ctx, "iam.cachePreKey", "nextpt:iam:").String()
	timeout := g.Cfg().MustGet(ctx, "iam.timeout", "10d").Duration()
	if timeout <= 0 {
		timeout = 10 * 24 * time.Hour
	}
	maxRefresh := g.Cfg().MustGet(ctx, "iam.maxRefresh", "5d").Duration()
	tokenDelimiter := g.Cfg().MustGet(ctx, "iam.tokenDelimiter", "_").String()

	tm := gtoken.NewDefaultToken(gtoken.Options{
		CacheMode:      cacheMode,
		CachePreKey:    cachePreKey,
		Timeout:        timeout.Milliseconds(),
		MaxRefresh:     maxRefresh.Milliseconds(),
		TokenDelimiter: tokenDelimiter,
	})

	middleware := gtoken.NewDefaultMiddleware(tm)
	middleware.ResFun = func(r *ghttp.Request, err error) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "Unauthorized: " + err.Error(),
			Data:    nil,
		})
		r.ExitAll()
	}

	return &sIamSessionDomain{
		gfToken:    tm,
		middleware: middleware,
		timeout:    timeout,
	}
}

func (s *sIamSessionDomain) GetGFMiddleware() gtoken.Middleware {
	return s.middleware
}

func (s *sIamSessionDomain) Create(ctx context.Context, userId uint64, deviceHash string, ip string, userAgent string) (string, *model.IamSession, error) {
	if userId == 0 || deviceHash == "" {
		return "", nil, errors.New("device identifier required")
	}
	sessionId, err := s.newSessionId()
	if err != nil {
		return "", nil, err
	}
	now := time.Now().UnixMilli()
	session := &model.IamSession{
		Id:         sessionId,
		UserId:     userId,
		DeviceHash: deviceHash,
		Ip:         s.limitString(ip, 64),
		UserAgent:  s.limitString(userAgent, 500),
		CreatedAt:  now,
		LastSeenAt: now,
	}

	token, err := s.gfToken.Generate(ctx, sessionId, g.Map{"userId": userId})
	if err != nil {
		return "", nil, err
	}
	if err = s.save(ctx, session); err != nil {
		_ = s.gfToken.Destroy(ctx, sessionId)
		return "", nil, err
	}
	indexKey := service.SysCache().KeyIamUserSessions(ctx, userId)
	if _, err = g.Redis().Do(ctx, "SADD", indexKey, sessionId); err != nil {
		_ = s.gfToken.Destroy(ctx, sessionId)
		_, _ = g.Redis().Do(ctx, "DEL", service.SysCache().KeyIamSession(ctx, sessionId))
		return "", nil, err
	}
	_, _ = g.Redis().Do(ctx, "PEXPIRE", indexKey, s.timeout.Milliseconds())
	return token, session, nil
}

func (s *sIamSessionDomain) RegisterDevice(ctx context.Context, userId uint64, deviceHash string) (bool, error) {
	if userId == 0 || deviceHash == "" {
		return false, nil
	}
	value, err := g.Redis().Do(
		ctx,
		"EVAL",
		iamRegisterDeviceScript,
		1,
		service.SysCache().KeyIamUserDevices(ctx, userId),
		deviceHash,
		int(iamKnownDeviceTTL.Seconds()),
	)
	if err != nil {
		return false, err
	}
	return value.Int() == 1, nil
}

func (s *sIamSessionDomain) Validate(ctx context.Context, token string, deviceHash string) (*model.IamSession, error) {
	sessionId, err := s.gfToken.Validate(ctx, token)
	if err != nil {
		return nil, err
	}
	session, err := s.Get(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if session == nil || session.UserId == 0 {
		_ = s.gfToken.Destroy(ctx, sessionId)
		return nil, errors.New("session invalid")
	}
	if deviceHash == "" || session.DeviceHash == "" || subtle.ConstantTimeCompare([]byte(session.DeviceHash), []byte(deviceHash)) != 1 {
		return nil, errors.New("session device mismatch")
	}

	now := time.Now().UnixMilli()
	if now-session.LastSeenAt >= iamSessionLastSeenInterval.Milliseconds() {
		session.LastSeenAt = now
		if err = s.save(ctx, session); err != nil {
			return nil, err
		}
		_, _ = g.Redis().Do(ctx, "PEXPIRE", service.SysCache().KeyIamUserSessions(ctx, session.UserId), s.timeout.Milliseconds())
	}
	return session, nil
}

func (s *sIamSessionDomain) Get(ctx context.Context, sessionId string) (*model.IamSession, error) {
	value, err := g.Redis().Do(ctx, "GET", service.SysCache().KeyIamSession(ctx, sessionId))
	if err != nil {
		return nil, err
	}
	if value == nil || value.IsNil() {
		return nil, nil
	}
	var session model.IamSession
	if err = json.Unmarshal([]byte(value.String()), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *sIamSessionDomain) ListByUser(ctx context.Context, userId uint64) ([]model.IamSession, error) {
	indexKey := service.SysCache().KeyIamUserSessions(ctx, userId)
	value, err := g.Redis().Do(ctx, "SMEMBERS", indexKey)
	if err != nil {
		return nil, err
	}
	if value == nil || value.IsNil() {
		return []model.IamSession{}, nil
	}
	ids := value.Strings()
	sessions := make([]model.IamSession, 0, len(ids))
	for _, sessionId := range ids {
		session, loadErr := s.Get(ctx, sessionId)
		if loadErr != nil {
			return nil, loadErr
		}
		if session == nil || session.UserId != userId {
			_, _ = g.Redis().Do(ctx, "SREM", indexKey, sessionId)
			_ = s.gfToken.Destroy(ctx, sessionId)
			continue
		}
		if _, _, loadErr = s.gfToken.Get(ctx, sessionId); loadErr != nil {
			_, _ = g.Redis().Do(ctx, "SREM", indexKey, sessionId)
			_, _ = g.Redis().Do(ctx, "DEL", service.SysCache().KeyIamSession(ctx, sessionId))
			continue
		}
		sessions = append(sessions, *session)
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastSeenAt > sessions[j].LastSeenAt
	})
	return sessions, nil
}

func (s *sIamSessionDomain) Remove(ctx context.Context, sessionId string) error {
	session, err := s.Get(ctx, sessionId)
	if err != nil {
		return err
	}
	if err = s.gfToken.Destroy(ctx, sessionId); err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "DEL", service.SysCache().KeyIamSession(ctx, sessionId))
	if err != nil {
		return err
	}
	if session != nil && session.UserId > 0 {
		_, err = g.Redis().Do(ctx, "SREM", service.SysCache().KeyIamUserSessions(ctx, session.UserId), sessionId)
	}
	return err
}

func (s *sIamSessionDomain) RemoveByUser(ctx context.Context, userId uint64) error {
	indexKey := service.SysCache().KeyIamUserSessions(ctx, userId)
	value, err := g.Redis().Do(ctx, "SMEMBERS", indexKey)
	if err != nil {
		return err
	}
	var sessionIds []string
	if value != nil && !value.IsNil() {
		sessionIds = value.Strings()
	}
	for _, sessionId := range sessionIds {
		if err = s.gfToken.Destroy(ctx, sessionId); err != nil {
			return err
		}
		if _, err = g.Redis().Do(ctx, "DEL", service.SysCache().KeyIamSession(ctx, sessionId)); err != nil {
			return err
		}
	}
	_ = s.gfToken.Destroy(ctx, strconv.FormatUint(userId, 10))
	_, err = g.Redis().Do(ctx, "DEL", indexKey)
	return err
}

func (s *sIamSessionDomain) save(ctx context.Context, session *model.IamSession) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "SET", service.SysCache().KeyIamSession(ctx, session.Id), payload, "PX", s.timeout.Milliseconds())
	return err
}

func (s *sIamSessionDomain) newSessionId() (string, error) {
	data := make([]byte, 16)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func (s *sIamSessionDomain) limitString(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max]
}
