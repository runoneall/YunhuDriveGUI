package structs

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UploadResult struct {
	Name   string   `json:"name"`
	Size   int64    `json:"size"`
	Chunks []string `json:"chunks"`
}

func (u *UploadResult) FormatJson() string {
	jsonData, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func (u *UploadResult) Delete() error {
	return os.Remove(filepath.Join("uploads", u.Name))
}
