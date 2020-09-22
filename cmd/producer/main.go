package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese.com/iap-polling/pkg/apple"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	version    string
	build      string
	production bool
)

func init() {
	flag.BoolVar(&production, "production", false, "Connect to production MySQL database if present. Default to localhost.")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	config.MustSetupViper()
}

func mustGetLogger() *zap.Logger {
	var logger *zap.Logger
	var err error
	if production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatal(err)
	}

	return logger
}

func main() {
	logger := mustGetLogger()
	dbConn := config.MustDBConn(production)
	kafkaAddr := config.MustKafkaAddress().PickSlice(production)

	p := apple.NewProducer(dbConn, kafkaAddr, logger)
	defer p.Close()

	c := cron.New(cron.WithLocation(chrono.TZShanghai))

	_, err := c.AddFunc("@daily", func() {
		p.Produce()
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	for {
		select {}
	}
}
