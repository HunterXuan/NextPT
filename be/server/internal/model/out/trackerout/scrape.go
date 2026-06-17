package trackerout

type ScrapeOut struct {
	Files map[string]ScrapeFile `bencode:"files"`
}

type ScrapeFile struct {
	Complete   uint `bencode:"complete"`
	Downloaded uint `bencode:"downloaded"`
	Incomplete uint `bencode:"incomplete"`
}
