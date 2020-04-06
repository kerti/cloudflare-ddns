package worker

import (
	"fmt"
	"net"
	"time"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/kerti/cloudflare-ddns/cloudflare"
	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/kerti/cloudflare-ddns/resolver"
	"github.com/spf13/viper"
)

// CheckConfig checks the worker configuration
func CheckConfig() error {
	checkInterval := viper.GetUint32("worker.checkInterval")
	if checkInterval < 10 {
		return fmt.Errorf("worker check interval cannot be set to under 10 seconds")
	}
	if checkInterval > 4294967295 {
		return fmt.Errorf("worker check interval cannot be set to over 4294967295 seconds")
	}
	logger.Debug("[WORKER] Worker configuration OK.")
	return nil
}

// Worker is the worker class
type Worker struct {
	Cloudflare cloudflare.Cloudflare
	Resolvers  []*resolver.Resolver
	Counter    int
	Hosts      []string
	HostMap    map[string]cf.DNSRecord
	CurrentIP  net.IP
}

func (w *Worker) initProperties() error {
	// initialize the resolvers
	resolvers, err := resolver.GetAll()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	w.Resolvers = resolvers

	// initialize counter
	w.Counter = 0

	// initialize hosts
	w.Hosts = viper.GetStringSlice("cloudflare.hostnames")

	// initialize cloudflare client
	cfClient, err := cloudflare.New()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	w.Cloudflare = *cfClient

	return nil
}

func (w *Worker) initExternal() error {
	// initialize hostmap
	w.getDNSRecords()

	// get current IP
	rslv := *w.Resolvers[len(w.Resolvers)-1]
	currentIP, err := rslv.GetExternalIP()
	if err != nil {
		logger.Error(err.Error())
	}
	w.CurrentIP = currentIP

	// run first check on host list
	err = w.checkHosts()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (w *Worker) init() error {
	// initialize properties
	err := w.initProperties()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// initialize stuff requiring external connections
	err = w.initExternal()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Debug("[WORKER] Worker initialization complete.")

	return nil
}

// Run runs the worker
func (w *Worker) Run() error {
	err := w.init()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	checkInterval := viper.GetUint32("worker.checkInterval")
	t := time.NewTicker(time.Duration(checkInterval) * time.Second)
	for range t.C {
		err := w.check()
		if err != nil {
			logger.Error(err.Error())
		}
	}

	return nil
}

func (w *Worker) check() error {
	err := w.getExternalIP()
	if err != nil {
		logger.Error(err.Error())
	}

	if w.Counter == len(w.Resolvers)-1 {
		w.getDNSRecords()
	}

	err = w.checkHosts()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (w *Worker) getExternalIP() error {
	if w.Counter >= len(w.Resolvers) {
		w.Counter = 0
	}

	rslv := *w.Resolvers[w.Counter]
	externalIP, err := rslv.GetExternalIP()
	w.Counter++

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	w.CurrentIP = externalIP
	return nil
}

func (w *Worker) getDNSRecords() {
	w.HostMap = make(map[string]cf.DNSRecord)
	for _, host := range w.Hosts {
		hostmap, err := w.Cloudflare.FetchA(host)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		for k, v := range hostmap {
			w.HostMap[k] = v
		}
	}
}

func (w *Worker) checkHosts() error {
	for _, host := range w.Hosts {
		rec, ok := w.HostMap[host]

		if !ok {
			logger.Debug("[WORKER] Host [%s] not found, adding...", host)
			res, err := w.Cloudflare.CreateA(host, w.CurrentIP)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			w.HostMap[host] = res
			continue
		}

		parsedContent := net.ParseIP(rec.Content)
		if parsedContent == nil {
			logger.Debug("[WORKER] Host [%s] has invalid IP address, setting...", host)
			err := w.Cloudflare.UpdateA(rec.ID, host, w.CurrentIP)
			if err != nil {
				logger.Error(err.Error())
			}
			continue
		}

		if parsedContent.String() != w.CurrentIP.String() {
			logger.Debug("[WORKER] Host [%s] has different IP address, setting...", host)
			err := w.Cloudflare.UpdateA(rec.ID, host, w.CurrentIP)
			if err != nil {
				logger.Error(err.Error())
			}
			continue
		}

		logger.Debug("[WORKER] Host [%s] has correct IP address set, skipping...", host)
	}

	return nil
}
