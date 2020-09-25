package main

import (
	"flag"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/migrate"
	"log"
)

var (
	production bool
	dir        string
)

func init() {
	flag.BoolVar(&production, "production", false, "Send verification request to production api or localhost")
	flag.StringVar(&dir, "dir", "", "Which directory to read: u, w, d")

	flag.Parse()

	config.MustSetupViper()
}

func main() {
	worker := migrate.NewWorker(production)

	var kind migrate.DirKind
	switch dir {
	case "u":
		kind = migrate.DirKindUUID
	case "w":
		kind = migrate.DirKindWxID
	case "d":
		kind = migrate.DirKindDevice
	default:
		kind = migrate.DirKindAll
	}

	log.Printf("Mirating receipts in %s", kind)

	err := worker.Start(kind)
	if err != nil {
		log.Fatal(err)
	}
}
