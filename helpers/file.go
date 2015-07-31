package helpers

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		//fmt.Println("Error Reading dir "+dir+":", err)
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

func downloadFromURL(url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}
