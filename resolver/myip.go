package resolver

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

type MyIPResponse struct {
	IP          string `json:"ip`
	Country     string `json:"country"`
	CountryCode string `json:"cc"`
}

// MyIP is the MyIP implementation
type MyIP struct {
	HTTPClient *http.Client
}

func (i *MyIP) Init() error {
	i.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("ipResolver.noVerify")},
		},
	}
	return nil
}

func (i *MyIP) GetExternalIP() (net.IP, error) {
	r, err := i.HTTPClient.Get("https://api.myip.com")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		var resp MyIPResponse
		json.NewDecoder(r.Body).Decode(&resp)
		parsedIP := net.ParseIP(resp.IP)
		if parsedIP == nil {
			return nil, fmt.Errorf("cannot parse IP: %s", resp.IP)
		}

		logger.Debug("[MyIP] Detected external IP: %v", parsedIP)
		return parsedIP, nil
	}

	return nil, err
}
