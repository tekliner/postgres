package leader_election

import (
	"io"
	"io/ioutil"
	"os"
)

func isDirectoryEmpty(directoryName string) (bool, error) {
	f, err := os.Open(directoryName)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// TODO: finish it
func getDirectoryContent(directoryName string) ([]string, error) {
	ioutil.ReadDir(directoryName)
	return []string{"a", "b"}, nil
}

