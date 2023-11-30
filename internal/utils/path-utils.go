package utils

import (
	"errors"
	"os"
	"path"
	"strings"
)

// CreateFoldersIfNotExist creates the folders that don't yet exist on the specified file path.
func CreateFoldersIfNotExist(filePath string) error {
	if len(strings.TrimSpace(filePath)) == 0 {
		return errors.New("a file path should be specified")
	}

	var err error
	pathSeparator := "/"
	directory, _ := path.Split(filePath)
	directories := strings.Split(directory, pathSeparator)
	for i, dir := range directories {
		if len(strings.TrimSpace(dir)) == 0 {
			// ignore empty directory
			continue
		}
		dirPath := strings.Join(directories[0:i+1], pathSeparator)
		err = createFolderIfNotExist(dirPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
