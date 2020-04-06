package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kerti/cloudflare-ddns/cloudflare"
	"github.com/kerti/cloudflare-ddns/config"
	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/kerti/cloudflare-ddns/worker"
)

/**********************************************************************/

/* command-line arguments */
var (
	optConfig = flag.String("config", "./config.yaml", "config file to use")
)

var (
	initialLoglevel uint8 = 6
	client          http.Client
)

/* version vars */
var (
	minVersion string
	majVersion string
	buildNum   string
	verSuffix  string
)

/**********************************************************************/

func initialize() error {
	fmt.Printf("[cloudflare-ddns] v%s.%s-%s build %s\n", majVersion, minVersion, verSuffix, buildNum)

	err := logger.InitLogger(&initialLoglevel)
	if err != nil {
		log.Printf("error setting up logger: %s", err.Error())
		return err
	}

	err = config.Load(optConfig)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = cloudflare.CheckConfig()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = worker.CheckConfig()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func main() {
	err := initialize()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	w := new(worker.Worker)
	err = w.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
