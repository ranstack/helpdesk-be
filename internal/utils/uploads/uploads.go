// Package uploads provides file upload and management utilities.
//
// Usage patterns:
// - Avatar updates: Get old user -> Save new avatar -> Update DB -> Delete old avatar
// - Ticket attachments: Save images -> Store URLs in DB -> Delete on attachment/ticket removal
// - User deletion: Get user -> Delete from DB -> Delete avatar file
package uploads

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	appErrors "helpdesk/internal/utils/errors"
)

const (
	MaxImageSize   = 5 * 1024 * 1024
	MaxFileSize    = 10 * 1024 * 1024
	ImageAvatarDir = "uploads/image/avatar"
	ImageTicketDir = "uploads/image/ticket"
	FileDir        = "uploads/file"
)

var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

var AllowedFileExtensions = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".txt":  true,
}

func EnsureUploadDirs() error {
	dirs := []string{ImageAvatarDir, ImageTicketDir, FileDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create upload directory %s: %w", dir, err)
		}
	}

	return nil
}

func ValidateImageFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > MaxImageSize {
		return appErrors.BadRequest("Image size exceeds maximum limit of 5MB")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !AllowedImageExtensions[ext] {
		return appErrors.BadRequest("Invalid image type. Only jpg, jpeg, png, and webp are allowed")
	}

	return nil
}

func ValidateDocumentFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > MaxFileSize {
		return appErrors.BadRequest("File size exceeds maximum limit of 10MB")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !AllowedFileExtensions[ext] {
		return appErrors.BadRequest("Invalid file type. Only pdf, doc, docx, xls, xlsx, and txt are allowed")
	}

	return nil
}

func SaveImageFile(fileHeader *multipart.FileHeader) (string, error) {
	if err := ValidateImageFile(fileHeader); err != nil {
		return "", err
	}

	return saveFile(fileHeader, ImageAvatarDir)
}

func SaveAvatarImage(fileHeader *multipart.FileHeader) (string, error) {
	if err := ValidateImageFile(fileHeader); err != nil {
		return "", err
	}

	return saveFile(fileHeader, ImageAvatarDir)
}

func SaveTicketImage(fileHeader *multipart.FileHeader) (string, error) {
	if err := ValidateImageFile(fileHeader); err != nil {
		return "", err
	}

	return saveFile(fileHeader, ImageTicketDir)
}

func SaveDocumentFile(fileHeader *multipart.FileHeader) (string, error) {
	if err := ValidateDocumentFile(fileHeader); err != nil {
		return "", err
	}

	return saveFile(fileHeader, FileDir)
}

func saveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return "/" + filepath.ToSlash(filePath), nil
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	cleanPath := strings.TrimPrefix(filePath, "/")

	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(cleanPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func DeleteFiles(filePaths []string) []error {
	var errors []error

	for _, path := range filePaths {
		if err := DeleteFile(path); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
