package model

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
)

func TestTrackerClientWhitelistJSONMapping(t *testing.T) {
	var rule TrackerClientWhitelist
	row := gdb.Record{
		"family":          gvar.New("uTorrent"),
		"peer_id_pattern": gvar.New("^-UT2[0-9]{3}-"),
		"exceptions":      gvar.New(`[{"name":"bad version","peerIdPattern":"^-UT2200-","agentPattern":"","enabled":true,"comment":"reason"}]`),
	}
	if err := row.Struct(&rule); err != nil {
		t.Fatal(err)
	}
	if rule.PeerIdPattern != "^-UT2[0-9]{3}-" || len(rule.Exceptions) != 1 || !rule.Exceptions[0].Enabled || rule.Exceptions[0].Name != "bad version" {
		t.Fatalf("invalid database mapping: %#v", rule)
	}
}
