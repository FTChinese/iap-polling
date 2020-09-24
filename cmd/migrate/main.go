package main

import (
	"flag"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/migrate"
	"log"
)

var (
	production bool
	dirKind    string
)

func init() {
	flag.BoolVar(&production, "production", false, "Send verification request to production api or localhost")
	flag.StringVar(&dirKind, "dirKind", "user-id", "Which directory to read: u, w, d")

	config.MustSetupViper()
}

func main() {
	worker := migrate.NewWorker(production)

	var kind migrate.NamingKind
	switch dirKind {
	case "u":
		kind = migrate.NamingKindUUID
	case "w":
		kind = migrate.NamingKindWxID
	case "d":
		kind = migrate.NamingKindDevice
	}
	err := worker.Start(kind)
	if err != nil {
		log.Fatal(err)
	}
}
