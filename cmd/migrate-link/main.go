package main

import (
	"github.com/FTChinese.com/iap-polling/pkg/apple"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"go.uber.org/zap"
	"log"
)

func main() {
	config.MustSetupViper()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	migrator := apple.NewLinkMigration(logger)

	err = migrator.Start()
	if err != nil {
		log.Fatal(err)
	}
}
