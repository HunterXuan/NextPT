package httpx

import (
	"mime"
	"strings"
	"testing"
)

func TestHeaderContentDispositionAttachmentSanitizesFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "path traversal uses basename",
			filename: "../../etc/passwd",
			want:     "passwd",
		},
		{
			name:     "windows path uses basename",
			filename: `C:\Users\me\movie.torrent`,
			want:     "movie.torrent",
		},
		{
			name:     "empty basename falls back",
			filename: "..",
			want:     "download",
		},
		{
			name:     "blank falls back",
			filename: " \t\n ",
			want:     "download",
		},
		{
			name:     "control bytes removed",
			filename: "bad\r\nname.torrent",
			want:     "badname.torrent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := HeaderContentDispositionAttachment(tt.filename)
			mediaType, params, err := mime.ParseMediaType(header)
			if err != nil {
				t.Fatalf("ParseMediaType(%q) error = %v", header, err)
			}
			if mediaType != "attachment" {
				t.Fatalf("media type = %q, want attachment", mediaType)
			}
			if params["filename"] != tt.want {
				t.Fatalf("filename = %q, want %q, header = %q", params["filename"], tt.want, header)
			}
			if strings.ContainsAny(header, "\r\n") {
				t.Fatalf("header contains newline: %q", header)
			}
		})
	}
}

func TestHeaderDeviceIdHash(t *testing.T) {
	hash := HeaderDeviceIdHash("  Browser-Device-001  ")
	if hash == "" {
		t.Fatal("HeaderDeviceIdHash() returned an empty hash")
	}
	if hash != HeaderDeviceIdHash("browser-device-001") {
		t.Fatal("HeaderDeviceIdHash() must normalize casing and whitespace")
	}
	if HeaderDeviceIdHash(" \t\n") != "" {
		t.Fatal("HeaderDeviceIdHash() must return an empty hash for an empty value")
	}
}
