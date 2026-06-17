package trackerin

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAnnounceEventUnmarshalCompactContract(t *testing.T) {
	data := []byte(`{"event":"started","torrentId":123,"userId":42,"peerId":"peer"}`)

	var event AnnounceEvent
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if event.TorrentId != 123 {
		t.Fatalf("TorrentId = %d, want 123", event.TorrentId)
	}
	if event.UserId != 42 {
		t.Fatalf("UserId = %d, want 42", event.UserId)
	}
}

func TestAnnounceEventMarshalUsesCompactContract(t *testing.T) {
	event := AnnounceEvent{
		Event:     "started",
		TorrentId: 123,
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	msg := string(data)
	if !strings.Contains(msg, `"torrentId":123`) {
		t.Fatalf("marshaled event %s does not contain compact torrentId", msg)
	}
	if strings.Contains(msg, `"torrent":`) || strings.Contains(msg, `"Torrent":`) {
		t.Fatalf("marshaled event %s should not contain full torrent payload", msg)
	}
}
