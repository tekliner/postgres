package leader_election

import (
	"io/ioutil"
	"log"
)

func getDirectoryContent(directoryName string) ([]string, error) {
	files, err := ioutil.ReadDir(directoryName)
	if err != nil {
		log.Print(err)
	}
	var filesList []string
	for _, file := range files {
		filesList = append(filesList, file.Name())
	}
	return filesList, err
}
