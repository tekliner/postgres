package leader_election

import (
	"io/ioutil"
	"log"
)

func getDirectoryContent(directoryName string) []string {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	var filesList []string
	for _, file := range files {
		filesList = append(filesList, file.Name())
	}
	return filesList
}
