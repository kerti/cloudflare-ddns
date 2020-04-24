package cloudflare

import (
	"fmt"
	"net"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

// CheckConfig checks the cloudflare configuration
func CheckConfig() error {
	email := viper.GetString("cloudflare.email")
	apiKey := viper.GetString("cloudflare.apiKey")
	zoneID := viper.GetString("cloudflare.zoneID")
	if len(email) == 0 || len(apiKey) == 0 || len(zoneID) == 0 {
		return fmt.Errorf("Cloudflare is not properly set up, check your config file")
	}
	logger.Debug("[CLOUDFLARE] Cloudflare configuration OK.")
	return nil
}

// Cloudflare is the cloudflare client
type Cloudflare struct {
	cf     *cf.API
	zoneID string
}

// New instantiates a new cloudflare client
func New() (*Cloudflare, error) {
	email := viper.GetString("cloudflare.email")
	apiKey := viper.GetString("cloudflare.apiKey")
	zoneID := viper.GetString("cloudflare.zoneID")

	api, err := cf.New(apiKey, email)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &Cloudflare{
		cf:     api,
		zoneID: zoneID,
	}, nil
}

// FetchA fetches an A-record
func (c *Cloudflare) FetchA(host string) (result map[string]cf.DNSRecord, err error) {
	logger.Debug("[CLOUDFLARE] Fetching A record for host [%v]", host)
	records, err := c.cf.DNSRecords(c.zoneID, cf.DNSRecord{Name: host, Type: "A"})
	if err != nil {
		return nil, fmt.Errorf("failed fetching A record: %v", err.Error())
	}

	result = make(map[string]cf.DNSRecord)
	for _, r := range records {
		logger.Debug("[CLOUDFLARE] IP Address for hostname [%s] is [%s]", r.Name, r.Content)
		result[r.Name] = r
	}

	return
}

// CreateA creates an A-record
func (c *Cloudflare) CreateA(name string, ip net.IP) (result cf.DNSRecord, err error) {
	logger.Debug("[CLOUDFLARE] Creating A record for host [%v]", name)
	rr := cf.DNSRecord{
		Type:    "A",
		Name:    name,
		Content: ip.String(),
	}

	cfResponse, err := c.cf.CreateDNSRecord(c.zoneID, rr)
	if err != nil {
		return result, fmt.Errorf("failed creating A record: %v", err.Error())
	}

	result = cfResponse.Result
	return
}

// UpdateA updates an A-record
func (c *Cloudflare) UpdateA(recordID string, name string, ip net.IP) (err error) {
	logger.Debug("[CLOUDFLARE] Updating A record for host [%v]", name)
	rr := cf.DNSRecord{
		Type:    "A",
		Name:    name,
		Content: ip.String(),
	}

	return c.cf.UpdateDNSRecord(c.zoneID, recordID, rr)
}
