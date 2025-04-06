package cmd

import (
	"bufio"
	"os"
)

func getStdout() *bufio.Writer {
	f := bufio.NewWriter(os.Stdout)
	return f
}

func getWriter(filename string, force bool) (*os.File, error) {
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
