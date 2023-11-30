package utils

import (
	"errors"
	"os"
	"path"
	"strings"
)

const (
	pathSeparator = "/"
)

// CreateFoldersIfNotExist creates the folders that don't yet exist on the specified file path.
func CreateFoldersIfNotExist(filePath string) error {
	if len(strings.TrimSpace(filePath)) == 0 {
		return errors.New("a file path should be specified")
	}

	var err error
	dir, _ := path.Split(filePath)
	dirs := strings.Split(dir, pathSeparator)
	for i, d := range dirs {
		if len(strings.TrimSpace(d)) == 0 {
			// ignore empty directory
			continue
		}
		dirPath := strings.Join(dirs[0:i+1], pathSeparator)
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
