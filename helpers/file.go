package helpers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// SearchFile ищет в заданной директории файл, по части его имени и возвращает полное название файла
func SearchFile(filePartName string, dir string) (string, error) {
	var err error
	var filesInfo []os.FileInfo
	var fullFileName string

	filesInfo, err = ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error Reading dir "+dir+":", err)
		return "", err
	}

	for i := range filesInfo {
		if strings.Contains(strings.ToLower(filesInfo[i].Name()), filePartName) {
			fullFileName = filesInfo[i].Name()
			break
		}
	}

	if fullFileName == "" {
		return "", errors.New("File not found")
	}

	return fullFileName, nil
}
