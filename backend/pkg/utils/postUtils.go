package utils

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

func isAllowedFile(filename string, file multipart.File) bool {
	// Check extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return false
	}

	// Read first 512 bytes for content sniffing
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false
	}
	// Reset file pointer ofr future reads
	file.Seek(0, 0)

	// Detect MIME type from sample (buffer)
	mimeType := http.DetectContentType(buffer)
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	return slices.Contains(allowedTypes, mimeType)
}
