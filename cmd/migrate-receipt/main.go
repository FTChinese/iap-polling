package main

import (
	"flag"
	"github.com/FTChinese.com/iap-polling/pkg/apple"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"go.uber.org/zap"
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
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	worker := apple.NewReceiptMigration(production, logger)

	log.Printf("Migrating receipts in %s", dir)

	err = worker.Start(dir)
	if err != nil {
		log.Fatal(err)
	}
}
