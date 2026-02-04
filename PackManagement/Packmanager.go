package PackManagement

/*
Packmanager.go - Install .bfk packs to Frames directory

Functions:
- InstallPack: Extract character folders from zip to Frames directory
*/

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func InstallPack(bfkPath, framesPath string) error {
	reader, err := zip.OpenReader(bfkPath)
	if err != nil {
		return fmt.Errorf("failed to open pack: %w", err)
	}
	defer reader.Close()

	cleanFramesPath := filepath.Clean(framesPath) + string(os.PathSeparator)

	for _, file := range reader.File {
		destPath := filepath.Join(framesPath, file.Name)

		if !strings.HasPrefix(filepath.Clean(destPath)+string(os.PathSeparator), cleanFramesPath) &&
			filepath.Clean(destPath) != filepath.Clean(framesPath) {
			continue
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		dstFile, err := os.Create(destPath)
		if err != nil {
			srcFile.Close()
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}

		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", destPath, err)
		}
	}

	return nil
}
