package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetStdout() *bufio.Writer {
	f := bufio.NewWriter(os.Stdout)
	return f
}

func GetWriter(filename string, force bool) (*os.File, error) {
	var permissions = os.O_WRONLY | os.O_CREATE
	if !force {
		permissions |= os.O_EXCL
	}

	f, err := os.OpenFile(filename, permissions, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func getReplacements(replacementsFilename string) (map[string]string, error) {
	replacementsHandle, err := os.Open(replacementsFilename)
	if err != nil {
		return nil, err
	}

	replacementsMap := make(map[string]string)

	replacementsScanner := bufio.NewScanner(replacementsHandle)
	replacementsScanner.Split(bufio.ScanLines)

	for replacementsScanner.Scan() {
		line := replacementsScanner.Text()
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, "=>")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid line, \"=>\" not found.")
		}

		replacementsMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return replacementsMap, nil
}

func getTransformedInput(inputFile string, replaceFile string) ([]byte, error) {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	replacementsMap, err := getReplacements(replaceFile)
	if err != nil {
		return nil, err
	}

	var transformedContent = string(content)
	for k, v := range replacementsMap {
		transformedContent = strings.ReplaceAll(transformedContent, k, v)
	}

	return []byte(transformedContent), nil
}

func GetContents(inputFile string, replaceFile string) ([]byte, error) {
	if replaceFile == "" {
		return os.ReadFile(inputFile)
	} else {
		return getTransformedInput(inputFile, replaceFile)
	}
}
