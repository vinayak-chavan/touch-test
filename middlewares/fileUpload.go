package middlewares

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// Middleware to upload file in local folder
func SaveFile(file *multipart.FileHeader) (string, error) {

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	timestamp := time.Now().Unix()
	dstPath := filepath.Join("uploads", fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename)))
	dst, err := os.Create(dstPath)

	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}
	return dstPath, nil
}
