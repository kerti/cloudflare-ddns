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

// IPAPICo is the IPAPICo implementation
type IPAPICo struct {
	HTTPClient *http.Client
}

func (i *IPAPICo) Init() error {
	i.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("ipResolver.noVerify")},
		},
	}
	return nil
}

func (i *IPAPICo) GetExternalIP() (net.IP, error) {
	r, err := i.HTTPClient.Get("https://ipapi.co/ip")
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

		logger.Debug("[IPAPICo] Detected external IP: %v", parsedIP)
		return parsedIP, nil
	}

	return nil, err
}
