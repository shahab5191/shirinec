package repositories

import (
	"log"
	"os"
	"path"
	"strings"

	"shirinec.com/config"
)

func LoadSqlFromFile(fileName string) (string, error) {
	filePath := path.Join(config.AppConfig.SqlFolder, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("[Error] - LoadSqlFromFile - Reading file - filePath: %s\nError: %+v\n", filePath, err)
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
