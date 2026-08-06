package httpx

import (
	"crypto/sha256"
	"encoding/hex"
	"mime"
	"path"
	"strings"
)

// HeaderDeviceIdHash returns a stable, non-reversible hash for a device ID header.
func HeaderDeviceIdHash(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:])
}

// HeaderContentDispositionAttachment formats a safe Content-Disposition value for downloads.
func HeaderContentDispositionAttachment(filename string) string {
	name := sanitizeAttachmentFilename(filename)
	if name == "" {
		name = "download"
	}

	value := mime.FormatMediaType("attachment", map[string]string{
		"filename": name,
	})
	if value == "" {
		return "attachment"
	}
	return value
}

func sanitizeAttachmentFilename(filename string) string {
	name := strings.ReplaceAll(filename, "\\", "/")
	name = path.Base(name)
	name = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, name)
	name = strings.TrimSpace(name)
	if name == "." || name == ".." || name == "/" {
		return ""
	}
	return name
}
