package migrate

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func mustHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return h
}

func WalkDir(ch chan<- string, dir string) error {
	defer close(ch)

	dir = filepath.Join(mustHomeDir(), dir)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Ignore sandbox receipt
		if strings.Contains(path, "Sandbox") {
			return nil
		}

		log.Println(path)

		ch <- path

		return nil
	})

	return err
}
