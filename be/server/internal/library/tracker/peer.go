package tracker

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
)

// PeerFormatZsetMember 格式化 ZSET 成员为 userId:peerId 字符串
func PeerFormatZsetMember(userId uint64, peerId string) string {
	return fmt.Sprintf("%d:%s", userId, peerId)
}

// PeerParseZsetMember 解析 ZSET 成员中的 userId 和 peerId
func PeerParseZsetMember(member string) (userId uint64, peerId string, ok bool) {
	firstColon := strings.Index(member, ":")
	if firstColon == -1 {
		return 0, "", false
	}
	return gconv.Uint64(member[:firstColon]), member[firstColon+1:], true
}

// PeerFormatUserSetValue 格式化用户做种/下载 Set 的值为 torrentId:peerId
func PeerFormatUserSetValue(torrentId uint64, peerId string) string {
	return fmt.Sprintf("%d:%s", torrentId, peerId)
}

// PeerParseUserSetValue 解析用户做种/下载 Set 的值，提取 torrentId 和 peerId
func PeerParseUserSetValue(val string) (torrentId uint64, peerId string, ok bool) {
	firstColon := strings.Index(val, ":")
	if firstColon == -1 {
		return 0, "", false
	}
	return gconv.Uint64(val[:firstColon]), val[firstColon+1:], true
}
