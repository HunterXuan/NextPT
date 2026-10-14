package tracker

import (
	"strings"
	"testing"
)

func TestLegacyPrefixLiteralMigration(t *testing.T) {
	for _, prefix := range []string{"-qB", "-TR", "a.b[0]", `foo\Ebar`, `^$()*+?{}|\`} {
		pattern := `^\Q` + strings.ReplaceAll(prefix, `\E`, `\E\\E\Q`) + `\E`
		rule, err := CompileClientRule(pattern, "")
		if err != nil {
			t.Fatal(err)
		}
		if !rule.Match(prefix+"rest", "") || rule.Match("x"+prefix+"rest", "") {
			t.Fatalf("prefix migration changed matching: %q", prefix)
		}
	}
}

func TestClientWhitelistLocalExceptions(t *testing.T) {
	for _, tt := range []struct {
		name        string
		peer, agent string
		exceptions  []ClientException
		allowed     bool
		blockedBy   string
	}{
		{name: "allow", peer: "-UT2210-", agent: "uTorrent/2210", allowed: true},
		{name: "wrong agent", peer: "-UT2210-", agent: "qBittorrent/4.5.0"},
		{name: "blocked by peer only", peer: "-UT2210-", agent: "uTorrent/2210", exceptions: []ClientException{{Name: "bad peer", PeerIdPattern: "^-UT2210-", Enabled: true}}, blockedBy: "bad peer"},
		{name: "exception requires both", peer: "-UT2210-", agent: "uTorrent/2210", exceptions: []ClientException{{Name: "other build", PeerIdPattern: "^-UT2210-", AgentPattern: "^uTorrent/2200", Enabled: true}}, allowed: true},
		{name: "disabled exception", peer: "-UT2210-", agent: "uTorrent/2210", exceptions: []ClientException{{Name: "disabled", PeerIdPattern: "^-UT2210-"}}, allowed: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := CompileClientWhitelist("^-UT2[0-9]{3}-", "^uTorrent/2[0-9]{3}", tt.exceptions)
			if err != nil {
				t.Fatal(err)
			}
			allowed, blockedBy := rule.Check(tt.peer, tt.agent)
			if allowed != tt.allowed || blockedBy != tt.blockedBy {
				t.Fatalf("got %v %q", allowed, blockedBy)
			}
		})
	}
	if _, err := CompileClientWhitelist("^-UT", "", []ClientException{{Name: "invalid", AgentPattern: "["}}); err == nil {
		t.Fatal("invalid exception accepted")
	}
}

func TestClientRule(t *testing.T) {
	for _, tt := range []struct {
		name, prefix, pattern, peer, agent string
		invalid, matched                   bool
	}{
		{name: "empty", invalid: true},
		{name: "invalid regex", prefix: "-qB", pattern: "[", invalid: true},
		{name: "prefix only", prefix: "-qB", peer: "-qB4500-abcdefghijk", matched: true},
		{name: "pattern only", pattern: "^qBittorrent/", agent: "qBittorrent/4.5.0", matched: true},
		{name: "both match", prefix: "-qB", pattern: "^qBittorrent/", peer: "-qB4500-abcdefghijk", agent: "qBittorrent/4.5.0", matched: true},
		{name: "wrong prefix", prefix: "-qB", pattern: "^qBittorrent/", peer: "-TR4000-abcdefghijk", agent: "qBittorrent/4.5.0"},
		{name: "wrong agent", prefix: "-qB", pattern: "^qBittorrent/", peer: "-qB4500-abcdefghijk", agent: "Transmission/4.0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := CompileClientRule(tt.prefix, tt.pattern)
			if (err != nil) != tt.invalid {
				t.Fatalf("error = %v", err)
			}
			if err == nil && rule.Match(tt.peer, tt.agent) != tt.matched {
				t.Fatal("unexpected match result")
			}
		})
	}
}
