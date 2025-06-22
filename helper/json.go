package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"yunhudrive/structs"
)

func SaveJson(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadUploadResult(filename string) (structs.UploadResult, error) {
	var result structs.UploadResult
	jsonData, err := os.ReadFile(filepath.Join("uploads", filename))
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
