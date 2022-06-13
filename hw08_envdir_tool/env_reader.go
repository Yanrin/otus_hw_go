package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Environment map[string]*EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envValues := make(Environment)

	if dir == "" {
		return envValues, nil
	}

	fileList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[^\w_]`)
	for _, file := range fileList {
		if file.IsDir() {
			continue // skip directory
		}

		fi, err := file.Info()
		if err != nil {
			return nil, err
		}

		name := fi.Name()
		if re.Match([]byte(name)) {
			continue // check incorrect symbols in filename
		}

		if fi.Size() == 0 {
			envValues[name] = &EnvValue{Value: "", NeedRemove: true}
			continue // variable needs to be deleted
		}

		line, err := firstFileLine(filepath.Join(dir, name)) // parse value from file
		if err != nil {
			return nil, err
		}
		envValues[name] = &EnvValue{Value: strings.TrimRight(line, "\t "), NeedRemove: false}
	}

	return envValues, nil
}

func firstFileLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var bline []byte
	for scanner.Scan() {
		bline = scanner.Bytes()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return prepareLine(bline), nil
}

func prepareLine(bline []byte) string {
	bline = bytes.ReplaceAll(bline, []byte{0}, []byte{'\n'})
	line := string(bline)
	line = strings.TrimRight(line, "\t ")
	return line
}
