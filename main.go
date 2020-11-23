package main

import (
	"flag"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/migrate"
	"log"
)

var (
	production bool // Determine whether hit production api.
	dir        string
)

func init() {
	flag.BoolVar(&production, "production", false, "Send verification request to production api or localhost")
	flag.StringVar(&dir, "dir", "", "Which directory to read receipts")

	flag.Parse()

	config.MustSetupViper()
}

func main() {
	worker := migrate.NewWorker(production)

	log.Printf("Migrating receipts in %s", dir)

	err := worker.Start(dir)
	if err != nil {
		log.Fatal(err)
	}
}
