package env

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Env interface {
	Load() error
	Get(key string) string
	GetEnvName() string
}

type env struct {
	FileName string
}

func New(fileName string) Env {
	return &env{FileName: fileName}
}

func (e *env) Load() error {
	rootDir, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	envPath := filepath.Join(rootDir, e.FileName)
	file, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func (e *env) Get(key string) string {
	return os.Getenv(key)
}

func (e *env) GetEnvName() string {
	return strings.ReplaceAll(e.FileName, ".", "")
}
