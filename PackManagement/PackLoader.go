package PackManagement

/*
PackLoader.go - Validate and preview .bfk pack files

Functions:
- ValidateBfkPack: Open zip, find character folders with frames
- GetPackPreviewImage: Extract first frame as base64 for preview
- GetPackInfo: Return pack metadata including characters and preview
*/

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

type PackInfo struct {
	FilePath     string   `json:"filePath"`
	PackName     string   `json:"packName"`
	Characters   []string `json:"characters"`
	PreviewImage string   `json:"previewImage"`
	Error        string   `json:"error,omitempty"`
}

func ValidateBfkPack(filePath string) (*PackInfo, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open pack: %w", err)
	}
	defer reader.Close()

	packName := strings.TrimSuffix(filepath.Base(filePath), ".bfk")
	characters := make(map[string]bool)
	var firstImagePath string
	var firstImageData []byte

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		parts := strings.Split(file.Name, "/")
		if len(parts) < 2 {
			continue
		}

		charName := parts[0]
		fileName := parts[len(parts)-1]
		ext := strings.ToLower(filepath.Ext(fileName))

		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			characters[charName] = true

			if firstImageData == nil {
				rc, err := file.Open()
				if err == nil {
					data, err := io.ReadAll(rc)
					rc.Close()
					if err == nil {
						firstImagePath = file.Name
						firstImageData = data
					}
				}
			}
		}
	}

	if len(characters) == 0 {
		return nil, fmt.Errorf("no valid character folders found in pack")
	}

	charList := make([]string, 0, len(characters))
	for c := range characters {
		charList = append(charList, c)
	}
	sort.Strings(charList)

	previewImage := ""
	if firstImageData != nil {
		ext := strings.ToLower(filepath.Ext(firstImagePath))
		mimeType := "image/png"
		if ext == ".jpg" || ext == ".jpeg" {
			mimeType = "image/jpeg"
		}
		previewImage = fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(firstImageData))
	}

	return &PackInfo{
		FilePath:     filePath,
		PackName:     packName,
		Characters:   charList,
		PreviewImage: previewImage,
	}, nil
}

func GetPackInfo(filePath string) PackInfo {
	info, err := ValidateBfkPack(filePath)
	if err != nil {
		return PackInfo{
			FilePath: filePath,
			Error:    err.Error(),
		}
	}
	return *info
}
