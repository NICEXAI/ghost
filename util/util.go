package util

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// CreateIfNotExist creates a file if it is not exists
func CreateIfNotExist(fileName, content string) error {
	pathList := strings.Split(fileName, "/")

	folder := strings.Join(pathList[:len(pathList)-1], "/")
	if _, err := os.Stat(folder); err != nil {
		if err = os.MkdirAll(folder, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	return err
}

// RemoveIfExist deletes the specified file if it is exists
func RemoveIfExist(filename string) error {
	if !IsFileExist(filename) {
		return nil
	}

	return os.Remove(filename)
}

// IsFileExist returns true if the specified file is exists
func IsFileExist(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsFolderExist returns true if the specified folder is exists
func IsFolderExist(folder string) bool {
	info, err := os.Stat(folder)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileInfo file info
type FileInfo struct {
	Name string
	Path string
}

// GetFileListFromFolder get all file from folder
func GetFileListFromFolder(folder string) (files []FileInfo, err error) {
	var fileInfoList []fs.FileInfo

	fileInfoList, err = ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			var childFiles []FileInfo
			childFiles, err = GetFileListFromFolder(path.Join(folder, fileInfo.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, childFiles...)
		} else {
			files = append(files, FileInfo{
				Name: fileInfo.Name(),
				Path: path.Join(folder, fileInfo.Name()),
			})
		}
	}

	return files, nil
}
