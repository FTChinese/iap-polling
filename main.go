package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese.com/iap-polling/pkg/config"
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

func main() {

}
