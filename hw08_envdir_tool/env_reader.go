package main

import (
	"bufio"
	"errors"
	"io/fs"
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

var ErrIncorrectFileName = errors.New("incorrect file name")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	envs := make(Environment, len(files))
	if err != nil {
		return Environment{}, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		key, env, err := readEnv(dir, file)
		if err != nil {
			return Environment{}, err
		}
		envs[key] = env
	}
	return envs, nil
}

func readEnv(dir string, entry fs.DirEntry) (string, EnvValue, error) {
	_, fileName := path.Split(entry.Name())
	if strings.Contains(fileName, "=") {
		return fileName, EnvValue{}, ErrIncorrectFileName
	}

	info, err := entry.Info()
	if err != nil {
		return fileName, EnvValue{}, err
	}
	if info.Size() == 0 {
		return fileName, EnvValue{
			NeedRemove: true,
		}, nil
	}

	file, err := os.Open(path.Join(dir, fileName))
	if err != nil {
		return fileName, EnvValue{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	return fileName, EnvValue{
		Value:      strings.ReplaceAll(strings.TrimRight(scanner.Text(), " \t"), "\x00", "\n"),
		NeedRemove: false,
	}, nil
}
