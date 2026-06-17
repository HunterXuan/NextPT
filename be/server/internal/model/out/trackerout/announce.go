package trackerout

type AnnouncePeer struct {
	PeerId string `bencode:"peer id"` // 20 bytes
	Ip     string `bencode:"ip"`
	Port   int    `bencode:"port"`
}

type AnnounceOut struct {
	Interval    int    `bencode:"interval"`
	MinInterval int    `bencode:"min interval,omitempty"`
	TrackerId   string `bencode:"tracker id,omitempty"`
	Complete    int    `bencode:"complete"`
	Incomplete  int    `bencode:"incomplete"`
	// Peers can be a string (binary compact format) or a list of dictionaries (non-compact)
	Peers  any `bencode:"peers"`
	Peers6 any `bencode:"peers6,omitempty"`
}
