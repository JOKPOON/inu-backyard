package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type FileInfoWithPath struct {
	Path string
	Info os.FileInfo
}

// deleteOldFiles deletes oldest files if file count exceeds limit
func DeleteOldFiles(folderPath string, maxFiles int) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read folder: %w", err)
	}

	// Filter out directories and prepare list with full paths
	var fileList []FileInfoWithPath
	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(folderPath, file.Name())
			info, err := file.Info()
			if err != nil {
				fmt.Printf("Failed to get file info for %s: %v\n", fullPath, err)
				continue
			}
			fileList = append(fileList, FileInfoWithPath{Path: fullPath, Info: info})
		}
	}

	// No need to delete if within limit
	if len(fileList) <= maxFiles {
		return nil
	}

	// Sort by modification time (oldest first)
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].Info.ModTime().Before(fileList[j].Info.ModTime())
	})

	// Calculate number of files to delete
	numToDelete := len(fileList) - maxFiles
	for i := range numToDelete {
		err := os.Remove(fileList[i].Path)
		if err != nil {
			fmt.Printf("Failed to delete %s: %v\n", fileList[i].Path, err)
		} else {
			fmt.Printf("Deleted: %s\n", fileList[i].Path)
		}
	}

	return nil
}
