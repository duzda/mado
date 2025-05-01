package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
			return nil, fmt.Errorf("invalid line, \"=>\" not found")
		}

		replacementsMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return replacementsMap, nil
}

func getTransformedContents(contents []byte, replaceFile string) ([]byte, error) {
	replacementsMap, err := getReplacements(replaceFile)
	if err != nil {
		return nil, err
	}

	var transformedContent = string(contents)
	for k, v := range replacementsMap {
		transformedContent = strings.ReplaceAll(transformedContent, k, v)
	}

	return []byte(transformedContent), nil
}

func getContentsInteractively() ([]byte, error) {
	file, err := os.CreateTemp("", "mado-*")
	defer func() {
		err = errors.Join(err, os.Remove(file.Name()))
	}()
	if err != nil {
		return nil, err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return nil, fmt.Errorf("$EDITOR is not set, can't get input interactively")
	}

	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	contents, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	return contents, err
}

func GetContents(inputFile string, replaceFile string) (contents []byte, err error) {
	if inputFile == "" {
		contents, err = getContentsInteractively()
	} else {
		contents, err = os.ReadFile(inputFile)
	}

	if err != nil {
		return nil, err
	}

	if replaceFile != "" {
		contents, err = getTransformedContents(contents, replaceFile)
		if err != nil {
			return nil, err
		}
	}

	return contents, nil
}

func JoinErrors(err *error, fn func() error) {
	*err = errors.Join(*err, fn())
}
