package worker

import (
	"math"
	"net"
	"os"
	"strconv"
	"time"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/kerti/cloudflare-ddns/cloudflare"
	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/kerti/cloudflare-ddns/notifier"
	"github.com/kerti/cloudflare-ddns/resolver"
	"github.com/spf13/viper"
)

// Worker is the worker class
type Worker struct {
	// Interval is the interval between checks
	Interval int

	// Cloudflare is the Cloudflare client
	Cloudflare cloudflare.Cloudflare

	// CloudflareErrorCount is the counter for cloudflare errors
	CloudflareErrorCount int

	// Resolvers is the collection of resolvers
	Resolvers []resolver.Resolver

	// ResolverErrorCount is the counter for resolver errors
	ResolverErrorCount int

	// Counter is the counter for the round-robin
	Counter int

	// Hosts is a list of hosts to be updated
	Hosts []string

	// HostMap is a map of existing hosts configured on Cloudflare
	HostMap map[string]cf.DNSRecord

	// CurrentIP is the current IP address as resolved by one of the resolvers
	CurrentIP net.IP
}

func (w *Worker) initInterval() {
	rslvLength := len(w.Resolvers)
	if rslvLength <= 0 {
		rslvLength = 1
	}

	checkIntervalStr := viper.GetString("worker.checkInterval")
	checkIntervalInt, err := strconv.Atoi(checkIntervalStr)

	if err != nil {
		if checkIntervalStr != "auto" {
			logger.Warn("Invalid check interval: [%s]", checkIntervalStr)
			logger.Warn("Reverting to automatic check interval.")
		}

		w.Interval = int(math.Round(float64(300) / float64(rslvLength)))
	} else {
		w.Interval = checkIntervalInt

		// prevent invoking a provider more than once per minute
		if rslvLength*checkIntervalInt < 300 {
			w.Interval = int(math.Round(float64(300) / float64(rslvLength)))
		}
	}

	// prevent checking more than twice per minute
	if w.Interval < 30 {
		w.Interval = 30
	}

	logger.Info("Check interval set at %d seconds.", w.Interval)
}

func (w *Worker) initProperties() error {
	// initialize the resolvers
	resolvers, err := resolver.Get()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	w.Resolvers = resolvers

	// initialize counter
	w.Counter = 0

	// initialize the interval
	w.initInterval()

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

// Init initializes the worker
func (w *Worker) Init() error {
	// initialize properties
	err := w.initProperties()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Debug("[WORKER] Worker initialization complete.")

	return nil
}

// Run runs the worker
func (w *Worker) Run() error {
	w.check()

	t := time.NewTicker(time.Duration(w.Interval) * time.Second)
	defer t.Stop()

	for range t.C {
		w.check()
	}

	return nil
}

// Stop stops the worker
func (w *Worker) Stop() {
	logger.Info("Interrupt detected, shutting down gracefully...")
	// Do whatever you need to do here before exiting
	logger.Info("Shutdown completed, exiting...")
	os.Exit(0)
}

func (w *Worker) check() {
	err := w.getExternalIP()
	if err != nil {
		logger.Warn("failed resolving external IP, skipping check: %v", err.Error())
		return
	}

	err = w.getDNSRecords()
	if err != nil {
		logger.Warn("failed retrieving existing DNS records, skipping check: %v", err.Error())
		return
	}

	err = w.checkHosts()
	if err != nil {
		logger.Warn("failed checking hosts, skipping check: %v", err.Error())
		return
	}
}

func (w *Worker) getExternalIP() error {
	if w.Counter >= len(w.Resolvers) {
		w.Counter = 0
	}

	rslv := w.Resolvers[w.Counter]
	externalIP, err := rslv.GetExternalIP()
	w.Counter++

	if err != nil {
		return err
	}

	w.CurrentIP = externalIP
	return nil
}

func (w *Worker) getDNSRecords() (err error) {
	w.HostMap = make(map[string]cf.DNSRecord)
	for _, host := range w.Hosts {
		hostmap, err := w.Cloudflare.FetchA(host)
		if err != nil {
			logger.Warn(err.Error())
			return err
		}
		for k, v := range hostmap {
			w.HostMap[k] = v
		}
	}

	return
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
		if parsedContent == nil || parsedContent.String() != w.CurrentIP.String() {
			logger.Debug("[WORKER] Host [%s] has different IP address, setting...", host)
			err := w.Cloudflare.UpdateA(rec.ID, host, w.CurrentIP)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			go w.setIP(host, w.CurrentIP)
			go w.notify(host, parsedContent.String(), w.CurrentIP.String())
			continue
		}

		logger.Debug("[WORKER] Host [%s] has correct IP address set, skipping...", host)
	}

	return nil
}

func (w *Worker) setIP(host string, newIP net.IP) {
	rec, ok := w.HostMap[host]
	if !ok {
		w.getDNSRecords()

		rec, ok = w.HostMap[host]
		if !ok {
			return
		}
	}

	rec.Content = newIP.String()
	w.HostMap[host] = rec
}

func (w *Worker) notify(host string, oldIP string, newIP string) {
	ifttt := notifier.IFTTT{V1: host, V2: oldIP, V3: newIP}
	err := ifttt.Notify()
	if err != nil {
		logger.Error(err.Error())
	}
}
