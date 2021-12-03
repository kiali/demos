package service

import (
	"bufio"
	"os"
	"strings"
)

func GetVersion(fileName string) string {
	version := "v1"
	labelsFile, err := os.Open(fileName)
	if err != nil {
		return version
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), "=")
		if values[0] == "version" {
			return values[1]
		}
	}
	return version
}
