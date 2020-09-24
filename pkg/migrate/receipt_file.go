package migrate

import (
	"log"
	"os"
	"path/filepath"
)

func mustHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return h
}

func WalkDir(ch chan<- string, k NamingKind) error {
	defer close(ch)

	dir := filepath.Join(mustHomeDir(), "receipt", k.String())

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		log.Println(path)

		ch <- path

		return nil
	})

	return err
}
