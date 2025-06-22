package helper

import (
	"os"
)

func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, 0755)
	return err
}
