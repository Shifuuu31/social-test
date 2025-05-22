package utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func IsAllowedFile(filename string, file multipart.File) bool {
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

func UploadHandler(file multipart.File, handler *multipart.FileHeader) (string, int) {
	// if r.Method != http.MethodPost {
	// 	utils.EncodeJson(w, 403, nil)
	// 	return
	// }

	// // Parse up to 10MB of form data
	// err := r.ParseMultipartForm(10 << 20) // 10MB
	// if err != nil {
	// 	http.Error(w, "Failed to parse form", http.StatusBadRequest)
	// 	return
	// }

	// Check if it's an image
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if (ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" ){
		return "", 400
	}

	if IsAllowedFile(handler.Filename, file) {
		return "", 400
	}

	// Make sure the uploads directory exists
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return "", 500
	}

	// Save the file
	dst, err := os.Create(filepath.Join("images", handler.Filename)) // TODO might wann add user spicific folder assignment
	if err != nil {
		return "", 500
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", 500
	}
	// = filepath.Join("images", handler.Filename)

	return filepath.Join("images", handler.Filename), 200
}
