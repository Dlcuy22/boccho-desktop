package AnimationEngine

/*
CharacterLoader.go - Utilities for discovering and loading characters

Functions:
- ScanCharacters: Scan Frames directory and return list of available characters
- GetCharacterFramesPath: Get full path to character's frames directory
- GetPreviewImage: Get path to first frame as preview thumbnail
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type CharacterInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	PreviewPath string `json:"previewPath"`
	FrameCount  int    `json:"frameCount"`
}

func ScanCharacters(basePath string) ([]CharacterInfo, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read characters directory: %w", err)
	}

	var characters []CharacterInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		charPath := filepath.Join(basePath, entry.Name())
		frames, err := filepath.Glob(filepath.Join(charPath, "*.png"))
		if err != nil {
			continue
		}

		if len(frames) == 0 {
			continue
		}

		sort.Strings(frames)

		characters = append(characters, CharacterInfo{
			Name:        entry.Name(),
			Path:        charPath,
			PreviewPath: frames[0],
			FrameCount:  len(frames),
		})
	}

	return characters, nil
}

func GetCharacterFramesPath(basePath, characterName string) string {
	return filepath.Join(basePath, characterName)
}

func GetPreviewImage(basePath, characterName string) (string, error) {
	charPath := filepath.Join(basePath, characterName)
	frames, err := filepath.Glob(filepath.Join(charPath, "*.png"))
	if err != nil {
		return "", err
	}

	if len(frames) == 0 {
		return "", fmt.Errorf("no frames found for character %s", characterName)
	}

	sort.Strings(frames)
	return frames[0], nil
}
