package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese.com/iap-polling/pkg/apple"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

var (
	version    string
	build      string
	production bool
	run        bool
)

func init() {
	flag.BoolVar(&production, "production", false, "Connect to production MySQL database if present. Default to localhost.")
	flag.BoolVar(&run, "run", false, "Run immediately")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	config.MustSetupViper()
}

func main() {
	logger := config.MustGetLogger(production)

	verifier := apple.NewVerifier(production, logger)
	defer verifier.Close()

	if run {
		err := verifier.Start()
		if err != nil {
			log.Printf("Running verifier error %v", err)
		}
		return
	}

	c := cron.New(cron.WithLocation(chrono.TZShanghai))

	_, err := c.AddFunc("@daily", func() {
		err := verifier.Start()
		if err != nil {
			log.Printf("Starting verifier error %v", err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	for {
		select {}
	}
}
