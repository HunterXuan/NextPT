package tracker

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"net"
	"regexp"
	"strings"
	"time"

	v1 "server/api/tracker/v1"
	"server/internal/consts"
	libtracker "server/internal/library/tracker"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/trackerin"
	"server/internal/model/out/trackerout"
	"server/internal/service"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	DefaultNumWant = 50
	MaxNumWant     = 100

	retrievePeerCandidateMultiplier = 6
	retrievePeerCandidateMax        = 600
)

type sTrackerPeerUsecase struct{}

func init() {
	service.RegisterTrackerPeerUsecase(NewTrackerPeerUsecase())
}

func NewTrackerPeerUsecase() *sTrackerPeerUsecase {
	return &sTrackerPeerUsecase{}
}

func (s *sTrackerPeerUsecase) Announce(ctx context.Context, actor *model.Actor, req *v1.AnnounceReq) (*trackerout.AnnounceOut, error) {
	r := ghttp.RequestFromCtx(ctx)
	infoHash := req.InfoHash
	peerId := req.PeerId

	if len(infoHash) != 20 || len(peerId) != 20 {
		return nil, errors.New("invalid info_hash or peer_id")
	}

	// 1. 获取种子信息 (带10分钟缓存)
	torrent, err := s.getTorrentByInfoHashCache(ctx, req.InfoHash)
	if err != nil {
		return nil, err
	}
	if torrent == nil {
		return nil, errors.New("unregistered torrent")
	}

	// 获取当前用户
	user := actor
	if user == nil {
		return nil, errors.New("user context missing")
	}

	// 统一 Tracker 级别的可见性与封禁拦截
	if err := service.CatalogTorrentDomain().CheckTorrentVisiblePolicy(ctx, user, torrent); err != nil {
		// 返回标准的 tracker 失败文本，最终会被外层包装成 bencode 的 failure reason
		return nil, errors.New("torrent not found or banned")
	}

	// 解析客户端 IP 与端口
	ipv4, ipv6, port, err := s.resolveClientAddress(r, req.Port)
	if err != nil {
		return nil, err
	}

	// 设置默认 NumWant
	numWant := s.resolveNumWant(req.NumWant)

	isSeeder := req.Left == 0

	// 只投递异步账务所需的最小事件数据，避免队列契约依赖 catalog 实体结构。
	event := &trackerin.AnnounceEvent{
		Event:      req.Event,
		TorrentId:  torrent.Id,
		Uploaded:   req.Uploaded,
		Downloaded: req.Downloaded,
		Left:       req.Left,
		UserId:     user.Id,
		PeerId:     peerId,
		Ipv4:       ipv4,
		Ipv6:       ipv6,
		Port:       port,
		IsSeeder:   isSeeder,
		Now:        gtime.Now(),
		UserAgent:  r.Header.Get("User-Agent"),
	}

	// 非阻塞发送，避免拖垮 Announce 主线程
	if err := service.TrackerEventUsecase().PushAnnounceEvent(event); err != nil {
		// 如果 Redis 入队失败，仅当是 stopped/completed 这样的终端/一次性事件才报错拦截，迫使客户端重试
		// 常规间隔汇报失败可以放行，下一次汇报的 diff 会自动弥补。
		if event.Event == consts.TrackerAnnounceEventStopped || event.Event == consts.TrackerAnnounceEventCompleted {
			return nil, errors.New("tracker service temporarily unavailable, please retry")
		}
	}

	// 获取活动 Peer
	peers, err := s.retrievePeers(ctx, torrent.Id, peerId, isSeeder, numWant, ipv4, ipv6)
	if err != nil {
		return nil, err
	}

	// 解析 Tracker 汇报间隔
	interval, minInterval := s.resolveInterval(ctx, torrent)

	out := &trackerout.AnnounceOut{
		Interval:    interval,
		MinInterval: minInterval,
		Complete:    int(torrent.Seeders),
		Incomplete:  int(torrent.Leechers),
	}

	out.Peers, out.Peers6 = s.encodePeers(req.Compact, peers)

	return out, nil
}

// resolveNumWant 确定最终返回的 peer 数量
func (s *sTrackerPeerUsecase) resolveNumWant(numWant int) int {
	if numWant <= 0 {
		return DefaultNumWant
	}
	if numWant > MaxNumWant {
		return MaxNumWant
	}
	return numWant
}

// resolveClientAddress 提取并验证客户端的 IPv4, IPv6 以及端口
func (s *sTrackerPeerUsecase) resolveClientAddress(r *ghttp.Request, port int) (string, string, int, error) {
	if port < 1 || port > 65535 {
		return "", "", 0, errors.New("invalid port")
	}

	var ipv4, ipv6 string

	// 辅助闭包：解析并赋值（如果已经满载则短路跳过）
	tryAssignIP := func(ipStr string) {
		if ipv4 != "" && ipv6 != "" {
			return
		}

		ipStr = strings.TrimSpace(ipStr)
		if ipStr == "" {
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			return
		}

		// 过滤内网和回环地址（严格模式下 Tracker 必须拒绝内网 IP 避免污染 Peer 列表）
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return
		}

		if p4 := ip.To4(); p4 != nil {
			if ipv4 == "" {
				ipv4 = p4.String() // 确保以标准 a.b.c.d 格式保存
			}
		} else {
			if ipv6 == "" {
				ipv6 = ip.String()
			}
		}
	}

	// 1. 最高优先级：真实的物理连接 IP (防 IP 欺骗)
	// 按照 BT BEP7 规范：如果客户端通过 IPv4 连接，则强制使用其实际 IPv4；它传的 ipv6 参数将被采纳。反之亦然。
	tryAssignIP(r.GetClientIp())

	// 2. 次优先级：客户端主动申报的 IPv4 / IPv6 (用于双栈补全)
	tryAssignIP(r.GetQuery("ipv4").String())
	tryAssignIP(r.GetQuery("ipv6").String())

	// 3. 最低优先级：传统 ip 字段 (兼容老客户端补全)
	tryAssignIP(r.GetQuery("ip").String())

	if ipv4 == "" && ipv6 == "" {
		return "", "", 0, errors.New("invalid or private ip address")
	}

	return ipv4, ipv6, port, nil
}

// retrievePeers 获取同伴列表并进行缓存、打散、双栈连通性匹配
func (s *sTrackerPeerUsecase) retrievePeers(ctx context.Context, torrentId uint64, peerId string, isSeeder bool, numWant int, reqIpv4, reqIpv6 string) ([]*entity.TrackerPeer, error) {
	now := time.Now().Unix()
	var peerIds []string
	candidateLimit := s.calcRetrievePeerCandidateLimit(numWant)

	if isSeeder {
		// 做种者只需要连下线 (leechers)
		leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, torrentId)
		lValues, err := s.retrievePeerMembers(ctx, leechersKey, now, candidateLimit)
		if err != nil {
			return nil, err
		}
		peerIds = lValues
	} else {
		// 下载者既要连上线 (seeders) 也要连下线 (leechers)
		seedersKey := service.SysCache().KeyTrackerTorrentSeeders(ctx, torrentId)
		leechersKey := service.SysCache().KeyTrackerTorrentLeechers(ctx, torrentId)

		sValues, err := s.retrievePeerMembers(ctx, seedersKey, now, candidateLimit)
		if err != nil {
			return nil, err
		}
		lValues, err := s.retrievePeerMembers(ctx, leechersKey, now, candidateLimit)
		if err != nil {
			return nil, err
		}
		peerIds = append(peerIds, sValues...)
		peerIds = append(peerIds, lValues...)
	}

	if len(peerIds) == 0 {
		return nil, nil
	}

	rand.Shuffle(len(peerIds), func(i, j int) {
		peerIds[i], peerIds[j] = peerIds[j], peerIds[i]
	})

	// 批量构建 Detail Key
	detailKeys := make([]any, 0, min(candidateLimit, len(peerIds)))
	for _, member := range peerIds {
		_, pId, ok := libtracker.PeerParseZsetMember(member)
		if !ok {
			continue
		}
		if pId == peerId {
			continue // 剔除自己
		}
		detailKeys = append(detailKeys, service.SysCache().KeyTrackerPeerDetail(ctx, torrentId, pId))
		if len(detailKeys) >= candidateLimit {
			break
		}
	}

	if len(detailKeys) == 0 {
		return nil, nil
	}

	// MGET 批量拉取详情
	bytesList, err := g.Redis().Do(ctx, "MGET", detailKeys...)
	if err != nil {
		return nil, err
	}

	var candidates []*entity.TrackerPeer
	for _, val := range bytesList.Vars() {
		if val.IsNil() || val.String() == "" {
			continue
		}
		var peer entity.TrackerPeer
		if err := json.Unmarshal(val.Bytes(), &peer); err != nil {
			continue
		}

		// 连通性过滤策略
		if reqIpv4 != "" && reqIpv6 == "" {
			if peer.Ipv4 == "" {
				continue
			}
		} else if reqIpv6 != "" && reqIpv4 == "" {
			if peer.Ipv6 == "" {
				continue
			}
		}

		candidates = append(candidates, &peer)
	}

	// 内存洗牌 (Fisher-Yates Shuffle)
	if len(candidates) > 0 {
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
	}

	// 截断 numWant
	if len(candidates) > numWant {
		candidates = candidates[:numWant]
	}

	return candidates, nil
}

func (s *sTrackerPeerUsecase) retrievePeerMembers(ctx context.Context, key string, now int64, limit int) ([]string, error) {
	if limit <= 0 {
		return nil, nil
	}

	countVal, err := g.Redis().Do(ctx, "ZCOUNT", key, now, "+inf")
	if err != nil {
		return nil, err
	}
	count := countVal.Int()
	if count <= 0 {
		return nil, nil
	}

	offset := 0
	if count > limit {
		offset = rand.Intn(count - limit + 1)
	}

	values, err := g.Redis().Do(ctx, "ZRANGEBYSCORE", key, now, "+inf", "LIMIT", offset, limit)
	if err != nil || values.IsNil() {
		return nil, err
	}
	return values.Strings(), nil
}

func (s *sTrackerPeerUsecase) calcRetrievePeerCandidateLimit(numWant int) int {
	limit := max(numWant*retrievePeerCandidateMultiplier, DefaultNumWant)
	if limit > retrievePeerCandidateMax {
		return retrievePeerCandidateMax
	}
	return limit
}

// resolveInterval 根据系统配置和种子活跃度动态计算下一次汇报的时间间隔
func (s *sTrackerPeerUsecase) resolveInterval(ctx context.Context, _ *entity.CatalogTorrent) (interval, minInterval int) {
	// interval = 客户端最小回报时间 (优先读取动态业务配置)
	interval = s.getTrackerConfigCache(ctx, consts.SiteConfigTrackerAnnounceInterval).Int()
	minInterval = s.getTrackerConfigCache(ctx, consts.SiteConfigTrackerAnnounceMinInterval).Int()

	// TODO: 后续可根据 torrent.Seeders 和 torrent.Leechers 动态调整
	// 比如：如果是一个死种（0 seeders, 0 leechers），可以让客户端 1-2 小时再来汇报，节约服务器资源

	return interval, minInterval
}

func (s *sTrackerPeerUsecase) getTrackerConfigCache(ctx context.Context, key string) *gvar.Var {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, key)
	val, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.SiteConfigDomain().GetByPath(ctx, key).Val(), nil
	}, 5*time.Minute)
	if err != nil || val.IsNil() {
		return service.SiteConfigDomain().GetByPath(ctx, key)
	}
	return gvar.New(val.Val())
}

// encodePeers 根据客户端请求生成对应的 peers 响应格式（支持 Compact 和字典模式，支持 IPv4/IPv6 双栈）
func (s *sTrackerPeerUsecase) encodePeers(compact int, peers []*entity.TrackerPeer) (peersOut any, peers6Out any) {
	if compact == 1 {
		// 生成 compact 格式
		var compactData4 []byte
		var compactData6 []byte
		for _, p := range peers {
			if p.Ipv4 != "" {
				ip := net.ParseIP(p.Ipv4)
				if ip != nil {
					if ip4 := ip.To4(); ip4 != nil {
						compactData4 = append(compactData4, ip4...)
						portBytes := make([]byte, 2)
						binary.BigEndian.PutUint16(portBytes, uint16(p.Port))
						compactData4 = append(compactData4, portBytes...)
					}
				}
			}
			if p.Ipv6 != "" {
				ip := net.ParseIP(p.Ipv6)
				if ip != nil {
					if ip16 := ip.To16(); ip16 != nil {
						compactData6 = append(compactData6, ip16...)
						portBytes := make([]byte, 2)
						binary.BigEndian.PutUint16(portBytes, uint16(p.Port))
						compactData6 = append(compactData6, portBytes...)
					}
				}
			}
		}

		if len(compactData6) > 0 {
			peers6Out = string(compactData6)
		}
		return string(compactData4), peers6Out
	}

	// 普通字典格式
	var peerDicts []trackerout.AnnouncePeer
	for _, p := range peers {
		if p.Ipv4 != "" {
			peerDicts = append(peerDicts, trackerout.AnnouncePeer{
				PeerId: string(p.PeerId),
				Ip:     p.Ipv4,
				Port:   int(p.Port),
			})
		}
		if p.Ipv6 != "" {
			peerDicts = append(peerDicts, trackerout.AnnouncePeer{
				PeerId: string(p.PeerId),
				Ip:     p.Ipv6,
				Port:   int(p.Port),
			})
		}
	}
	return peerDicts, nil
}

// Scrape 批量查询种子的做种、下载和完成数
func (s *sTrackerPeerUsecase) Scrape(ctx context.Context, actor *model.Actor, infoHashes []string) (*trackerout.ScrapeOut, error) {
	if len(infoHashes) == 0 {
		return &trackerout.ScrapeOut{Files: make(map[string]trackerout.ScrapeFile)}, nil
	}

	out := &trackerout.ScrapeOut{
		Files: make(map[string]trackerout.ScrapeFile, len(infoHashes)),
	}

	for _, hashStr := range infoHashes {
		// 复用 Resource 模块的高效缓存 (10分钟有效期)，完美避免 DB 扫表
		t, err := s.getTorrentByInfoHashCache(ctx, hashStr)
		if err != nil || t == nil {
			// 如果种子不存在，按 BT 协议规范可以直接忽略该 info_hash，不予返回
			continue
		}

		if err := service.CatalogTorrentDomain().CheckTorrentVisiblePolicy(ctx, actor, t); err != nil {
			continue
		}

		out.Files[hashStr] = trackerout.ScrapeFile{
			Complete:   t.Seeders,
			Downloaded: t.TimesCompleted,
			Incomplete: t.Leechers,
		}
	}

	return out, nil
}

// CachedAgentRule 用于缓存编译好的正则规则，避免每次请求都重新编译
type CachedAgentRule struct {
	PeerIdPrefix string
	AgentRegex   *regexp.Regexp // 如果为空，表示不强制校验 Agent
}

// CheckClientWhitelist 检查客户端白名单 (内存全量缓存)
func (s *sTrackerPeerUsecase) CheckClientWhitelist(ctx context.Context, peerId, userAgent string) bool {
	// 获取全量白名单缓存 (1小时过期)
	val, err := gcache.GetOrSetFunc(ctx, service.SysCache().KeyTrackerClientWhitelist(ctx), func(ctx context.Context) (any, error) {
		// 调用 Domain 层获取纯粹的白名单数据
		list, err := service.TrackerPeerDomain().GetEnabledClientWhitelists(ctx)
		if err != nil {
			return nil, err
		}

		var compiledRules []CachedAgentRule
		for _, rule := range list {
			var re *regexp.Regexp
			if rule.AgentPattern != "" {
				// 如果正则配置错误，则忽略该条规则或跳过正则解析
				re, _ = regexp.Compile(rule.AgentPattern)
			}
			compiledRules = append(compiledRules, CachedAgentRule{
				PeerIdPrefix: rule.PeerIdPrefix,
				AgentRegex:   re,
			})
		}

		return compiledRules, nil
	}, time.Hour)

	if err != nil || val.IsNil() {
		return false
	}

	var whitelist []CachedAgentRule
	if err := val.Scan(&whitelist); err != nil {
		return false
	}

	// 遍历白名单进行匹配
	for _, rule := range whitelist {
		// 1. PeerIdPrefix 匹配 (如果有)
		if rule.PeerIdPrefix != "" && !strings.HasPrefix(peerId, rule.PeerIdPrefix) {
			continue
		}

		// 2. AgentPattern 匹配 (使用预编译正则)
		if rule.AgentRegex != nil && !rule.AgentRegex.MatchString(userAgent) {
			continue
		}

		// 只要有一条规则完全命中，就放行
		return true
	}

	return false
}

// getTorrentByInfoHashCache 根据 info_hash 获取种子信息并在 Tracker 应用层进行内存/Redis缓存
func (s *sTrackerPeerUsecase) getTorrentByInfoHashCache(ctx context.Context, infoHash string) (*entity.CatalogTorrent, error) {
	cacheKey := service.SysCache().KeyCatalogTorrentInfoHash(ctx, hex.EncodeToString([]byte(infoHash)))
	v, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.CatalogTorrentDomain().GetTorrentByInfoHash(ctx, infoHash)
	}, 10*time.Minute)

	if err != nil || v.IsNil() {
		return nil, err
	}

	var torrent entity.CatalogTorrent
	if err := v.Scan(&torrent); err != nil {
		return nil, err
	}

	return &torrent, nil
}
