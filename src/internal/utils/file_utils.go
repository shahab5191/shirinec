package utils

import (
	"os"
	"path"

	"shirinec.com/config"
)

func RemoveMedia(filename string) error {
    mediaFolder := config.AppConfig.UploadFolder
    fileAddress := path.Join(mediaFolder, filename)

    err := os.Remove(fileAddress)
    return err
}
