package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrEmptyPath                 = errors.New("empty path")
	ErrFileNameContainsEqualSign = errors.New("file name contains equal sign")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}

	if dir == "" {
		return nil, ErrEmptyPath
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			return nil, ErrFileNameContainsEqualSign
		}

		filePath := filepath.Join(dir, file.Name())

		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		if fileInfo.Size() == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		line, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		value := string(bytes.ReplaceAll(line, []byte{0x00}, []byte("\n")))
		value = strings.TrimRight(value, " \t")
		env[file.Name()] = EnvValue{Value: value}
	}

	return env, nil
}
