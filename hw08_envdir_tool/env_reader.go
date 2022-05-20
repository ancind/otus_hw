package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrUnsupportedInput = errors.New("input path isn't a directory")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrUnsupportedInput
	}

	env := make(Environment)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key, value := pair[0], pair[1]
		v := EnvValue{Value: value}
		env[key] = v
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range files {
		name := fileInfo.Name()
		filePath := path.Join(dir, name)

		if fileInfo.Size() == 0 {
			delete(env, name)
			continue
		}
		if strings.Contains(name, "=") {
			continue
		}

		lines, err := readLines(filePath)
		if err != nil {
			continue
		}

		value := cleanStringValue(lines[0])
		v := EnvValue{Value: value}
		env[name] = v
	}

	return env, nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func cleanStringValue(value string) string {
	value = strings.TrimRight(value, " \t\n")
	value = string(bytes.ReplaceAll([]byte(value), []byte("\x00"), []byte("\n")))
	return value
}
