package tracker

import (
	"encoding/json"
	"math/rand"
	"os"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
	"testing"
)

type legacyBootstrapRule struct {
	Id                           int
	Peer, PeerStart, PeerType    string
	Agent, AgentStart, AgentType string
	PeerCount, AgentCount        int
}

func bootstrapClientRules(t *testing.T) []ClientWhitelist {
	t.Helper()
	data, err := os.ReadFile("../../../storage/migration/bootstrap.sql")
	if err != nil {
		t.Fatal(err)
	}
	row := regexp.MustCompile(`(?m)^    \(([0-9]+), '([^']*)', '([^']*)', '([^']*)', '([^']*)', b'1', '([^']*)', NOW\(\), NOW\(\)\)[,;]$`)
	matches := row.FindAllStringSubmatch(string(data), -1)
	if len(matches) != 15 {
		t.Fatalf("expected 15 seed rows, got %d", len(matches))
	}
	ids := make(map[string]bool)
	uTorrentGroups := make(map[string]bool)
	delugeGroups := 0
	var rules []ClientWhitelist
	for index, match := range matches {
		if match[1] != strconv.Itoa(index+1) {
			t.Fatalf("expected consecutive seed ID %d, got %s", index+1, match[1])
		}
		if ids[match[1]] {
			t.Fatalf("duplicate seed ID %s", match[1])
		}
		ids[match[1]] = true
		if strings.HasPrefix(match[2], "Deluge") {
			delugeGroups++
			if match[2] != "Deluge 2.x" {
				t.Fatalf("unexpected Deluge group %s", match[2])
			}
		}
		if strings.HasPrefix(match[2], "uTorrent") {
			if uTorrentGroups[match[2]] {
				t.Fatalf("duplicate uTorrent group %s", match[2])
			}
			uTorrentGroups[match[2]] = true
		}
		if len(match[2]) > 60 || len(match[3]) > 200 || len(match[4]) > 200 || len(match[6]) > 200 {
			t.Fatalf("seed %s exceeds column limits", match[1])
		}
		var exceptions []ClientException
		if err := json.Unmarshal([]byte(match[5]), &exceptions); err != nil {
			t.Fatal(err)
		}
		rule, err := CompileClientWhitelist(match[3], match[4], exceptions)
		if err != nil {
			t.Fatalf("seed %s: %v", match[1], err)
		}
		rules = append(rules, rule)
	}
	for _, family := range []string{"uTorrent 1.x", "uTorrent 2.x", "uTorrent 3.x", "uTorrent Mac 1.x"} {
		if !uTorrentGroups[family] {
			t.Fatalf("missing consolidated group %s", family)
		}
	}
	if len(uTorrentGroups) != 4 {
		t.Fatalf("unexpected uTorrent sub-version groups: %v", uTorrentGroups)
	}
	if delugeGroups != 1 {
		t.Fatalf("expected one Deluge 2.x group, got %d", delugeGroups)
	}
	return rules
}

func bootstrapAllows(rules []ClientWhitelist, peer, agent string) bool {
	for _, rule := range rules {
		if allowed, _ := rule.Check(peer, agent); allowed {
			return true
		}
	}
	return false
}

func TestClientBootstrap(t *testing.T) {
	rules := bootstrapClientRules(t)
	for _, tt := range []struct {
		peer, agent string
		allowed     bool
	}{
		{"-qB4010-", "qBittorrent/4.0.1", true},
		{"-qB4000-", "qBittorrent/4.0.0", false},
		{"-qB4500-", "qBittorrent/4.5.0", false},
		{"-qB4510-", "qBittorrent/4.5.1", false},
		{"-qB4520-", "qBittorrent/4.5.2", true},
		{"-qB41A0-", "qBittorrent/4.1.10", true},
		{"-qB5000-", "qBittorrent/5.0.0.0", true},
		{"-qB3300-", "qBittorrent v3.3.0", true},
		{"-qB3200-", "qBittorrent v3.2.0", false},
		{"-UT2210-", "uTorrent/2210", true},
		{"-UT1610-", "uTorrent/1610", false},
		{"-UT1750-", "uTorrent/1750", false},
		{"-UT1740-", "uTorrent/1740", false},
		{"-UT1800-", "uTorrent/1800", true},
		{"-UT3420-", "uTorrent/3420", true},
		{"-UT3490-", "uTorrent/3490", true},
		{"-UT3500-", "uTorrent/3500", true},
		{"-UT355W-", "uTorrent/3550", true},
		{"-UT340-", "uTorrent/340", false},
		{"-UT3420-", "uTorrent/342", false},
		{"-UT355W-", "uTorrent/355", false},
		{"-UT355W-", "uTorrent/35500", false},
		{"-UT355W-", "uTorrent/3550 (build 12345)", true},
		{"-UM1870-", "uTorrentMac/1870", true},
		{"-UT180B-", "uTorrent/180B(6838)", false},
		{"-UT180B-", "uTorrent/180B(9136)", false},
		{"-UT180B-", "uTorrent/180B(9137)", false},
		{"-UT180B-", "uTorrent/180B(10000)", false},
		{"-TR2940-", "Transmission/2.94", true},
		{"-TR3000-", "Transmission/3.00", true},
		{"-TR4000-", "Transmission/4.0.0", true},
		{"-DE13D0-", "Deluge 1.3.13", false},
		{"-DE13F0-", "Deluge 1.3.15", false},
		{"-DE203s-", "Deluge/2.0.3 libtorrent/1.1.13.0", true},
		{"-DE203s-", "Deluge/2.0.3", true},
		{"-DE211s-", "Deluge/2.1.1 libtorrent/2.0.10.0", true},
		{"-DE20As-", "Deluge/2.0.10", true},
		{"-DE21Gs-", "Deluge/2.1.16", true},
		{"-DE2A0s-", "Deluge/2.10.0", true},
		{"-DE202s-", "Deluge/2.0.2", false},
		{"-DE203s-", "Deluge/2.0.2", false},
		{"-DE202s-", "Deluge/2.0.3", false},
		{"-DE211b-", "Deluge/2.1.1b1", false},
		{"-DE211a-", "Deluge/2.1.1a1", false},
		{"-DE211r-", "Deluge/2.1.1rc1", false},
		{"-DE211D-", "Deluge/2.1.1.dev1", false},
		{"-DE211s-", "Deluge/2.1.1rc1", false},
		{"-DE211s-", "Deluge/2.1.1.dev1", false},
		{"-DE211s-", "Deluge/2.1.1+dev", false},
		{"-DE2110-", "Deluge/2.1.1", false},
		{"-DE211s-", "Deluge/2.1.1junk", false},
		{"-DE311s-", "Deluge/3.1.1", false},
		{"-lt0C00-", "rtorrent/0.8.0/0.12.0", true},
		{"-lt0B00-", "rtorrent/0.8.0/0.12.0", false},
		{"-lt0C00-", "rtorrent/0.8.0/0.11.0", false},
		{"-lt0D00-", "rtorrent/0.9.0/0.11.0", true},
		{"-BR0332-", "BitRocket/0.3.3(32) libtorrent/0.13.0.0", false},
		{"-BR0332-", "BitRocket/0.3.3(32) libtorrent/0.12.0.0", false},
		{"-BR0333-", "BitRocket/0.3.3(33) libtorrent/0.12.0.0", false},
		{"-AZ2504-", "Azureus 2.5.0.4", false},
		{"-AZ3050-", "Azureus 3.0.5.0", false},
		{"-AZ4002-", "Azureus 4.0.0.2", false},
		{"-DE0589-", "Deluge 0.5.8.9", false},
		{"-TR1060-", "Transmission/1.06", false},
		{"-ML2.9.2-", "MLDonkey/2.9.2", false},
		{"unknown", "Rufus/0.6.9", false},
		{"unknown", "SymTorrent", false},
		{"-DE1160-", "Deluge 1.1.6", false},
		{"M6-0-1--", "BitTorrent/6010", true},
		{"-BT7910-", "BitTorrent/7910", true},
		{"M7-0-1--", "BitTorrent/7010", true},
		{"-BT7010-", "BitTorrent/7010", true},
		{"M7-0-0--", "BitTorrent/7000", false},
		{"-BT7000-", "BitTorrent/7000", false},
		{"-BT7910-", "BitTorrent/7000", false},
		{"-BT7910-", "BitTorrent/791", false},
		{"-BT7910-", "BitTorrent/6010", false},
		{"-AZ5600-", "Azureus 5.6.0.0", true},
		{"-qB5000-", "Transmission/4.0.0", false},
		{"unknown", "unknown", false},
	} {
		t.Run(tt.peer+tt.agent, func(t *testing.T) {
			if got := bootstrapAllows(rules, tt.peer+"abcdefghijkl", tt.agent); got != tt.allowed {
				t.Fatalf("allowed = %v, want %v", got, tt.allowed)
			}
		})
	}
}

// Compare the consolidated rules with the old numeric lower bounds under the
// agreed any-passing-rule policy. Literal separators and Peer ID anchors are
// normalized in the fixture; legacy HTTPS and first-match behavior are excluded.
// Only retained families are compared. Retired clients and the simplified
// uTorrent 1.x/3.x, Deluge 2.x and the unified -BT7 minimum version are covered
// by explicit acceptance/rejection cases.
func TestClientBootstrapLegacyBoundaries(t *testing.T) {
	rules := bootstrapClientRules(t)
	data, err := os.ReadFile("testdata/client_whitelist_legacy.json")
	if err != nil {
		t.Fatal(err)
	}
	var legacy []legacyBootstrapRule
	if err := json.Unmarshal(data, &legacy); err != nil {
		t.Fatal(err)
	}
	legacyMatch := func(pattern, start, value, kind string, count int) bool {
		if pattern == "" {
			return true
		}
		re := regexp.MustCompile(pattern)
		base, target := re.FindStringSubmatch(start), re.FindStringSubmatch(value)
		if target == nil || base == nil {
			return false
		}
		number := func(s string) int64 {
			base := 10
			if kind == "hex" {
				base = 16
			}
			n, _ := strconv.ParseInt(s, base, 64)
			return n
		}
		for i := 1; i <= count && i < len(base) && i < len(target); i++ {
			if a, b := number(target[i]), number(base[i]); a != b {
				return a > b
			}
		}
		return true
	}
	expected := func(peer, agent string) bool {
		for _, r := range legacy {
			if !legacyMatch(r.Peer, r.PeerStart, peer, r.PeerType, r.PeerCount) || !legacyMatch(r.Agent, r.AgentStart, agent, r.AgentType, r.AgentCount) {
				continue
			}
			if r.Id == 36 && ((strings.HasPrefix(peer, "-qB4500-") && agent == "qBittorrent/4.5.0") || (strings.HasPrefix(peer, "-qB4510-") && agent == "qBittorrent/4.5.1")) {
				continue
			}
			return true
		}
		return false
	}
	rng := rand.New(rand.NewSource(42))
	var sample func(*syntax.Regexp) string
	sample = func(re *syntax.Regexp) string {
		switch re.Op {
		case syntax.OpLiteral:
			return string(re.Rune)
		case syntax.OpCharClass:
			i := rng.Intn(len(re.Rune)/2) * 2
			return string(re.Rune[i] + rune(rng.Intn(int(re.Rune[i+1]-re.Rune[i]+1))))
		case syntax.OpCapture:
			return sample(re.Sub[0])
		case syntax.OpConcat:
			var result string
			for _, sub := range re.Sub {
				result += sample(sub)
			}
			return result
		case syntax.OpAlternate:
			return sample(re.Sub[rng.Intn(len(re.Sub))])
		case syntax.OpStar, syntax.OpPlus:
			count := rng.Intn(3)
			if re.Op == syntax.OpPlus {
				count++
			}
			var result string
			for i := 0; i < count; i++ {
				result += sample(re.Sub[0])
			}
			return result
		default:
			return ""
		}
	}
	for _, r := range legacy {
		switch r.Id {
		case 3, 6, 17, 18, 23, 24, 27, 35, 36, 41, 43, 44:
		default:
			continue
		}
		peerRE, err := syntax.Parse(r.Peer, syntax.Perl)
		if err != nil {
			t.Fatal(err)
		}
		agentRE, err := syntax.Parse(r.Agent, syntax.Perl)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 500; i++ {
			peer, agent := sample(peerRE), sample(agentRE)
			if i == 0 {
				peer, agent = r.PeerStart, r.AgentStart
			}
			if got, want := bootstrapAllows(rules, peer, agent), expected(peer, agent); got != want {
				t.Fatalf("legacy %d: %q %q: got %v, want %v", r.Id, peer, agent, got, want)
			}
		}
	}
}
