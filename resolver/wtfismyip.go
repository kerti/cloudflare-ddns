package resolver

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

// WTFIsMyIP is the WTFIsMyIP implementation
type WTFIsMyIP struct {
	HTTPClient *http.Client
}

func (i *WTFIsMyIP) Init() error {
	i.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("ipResolver.noVerify")},
		},
	}
	return nil
}

func (i *WTFIsMyIP) GetExternalIP() (net.IP, error) {
	r, err := i.HTTPClient.Get("https://wtfismyip.com/text")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		ipString := string(bodyBytes[:len(bodyBytes)-1])
		parsedIP := net.ParseIP(ipString)
		if parsedIP == nil {
			return nil, fmt.Errorf("cannot parse IP: %s", ipString)
		}

		logger.Debug("[WTFIsMyIP] Detected external IP: %v", parsedIP)
		return parsedIP, nil
	}

	return nil, err
}
