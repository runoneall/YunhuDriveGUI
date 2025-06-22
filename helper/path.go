package helper

import (
	"os"
	"strings"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func ListFiles(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []string{}
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			files = append(files, entry.Name())
		}
	}
	return files
}
