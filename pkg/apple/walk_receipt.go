package apple

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

func WalkReceipts(dir string) <-chan string {

	absDir := filepath.Join(mustHomeDir(), dir)
	log.Printf("Walking receipts under %s", absDir)

	ch := make(chan string)

	go func() {
		defer close(ch)

		err := filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
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

			if path == ".DS_Store" {
				return nil
			}

			log.Println(path)

			ch <- path

			return nil
		})

		if err != nil {
			log.Println(err.Error())
		}
	}()

	return ch
}
