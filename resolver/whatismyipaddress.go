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

// WhatIsMyIPAddress is the WhatIsMyIPAddress implementation
type WhatIsMyIPAddress struct {
	HTTPClient *http.Client
}

func (i *WhatIsMyIPAddress) Init() error {
	i.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("ipResolver.noVerify")},
		},
	}
	return nil
}

func (i *WhatIsMyIPAddress) GetExternalIP() (net.IP, error) {
	r, err := i.HTTPClient.Get("http://ipv4bot.whatismyipaddress.com/")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		ipString := string(bodyBytes)
		parsedIP := net.ParseIP(ipString)
		if parsedIP == nil {
			return nil, fmt.Errorf("cannot parse IP: %s", ipString)
		}

		logger.Debug("[WhatIsMyIPAddress] Detected external IP: %v", parsedIP)
		return parsedIP, nil
	}

	return nil, err
}
