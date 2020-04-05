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

type BigDataCloudResponse struct {
	IPString      string `json:"ipString`
	IPNumeric     int    `json:"ipNumeric"`
	IPType        string `json:"ipType"`
	IsBehindProxy string `json:"isBehindProxy"`
}

// BigDataCloud is the BigDataCloud implementation
type BigDataCloud struct {
	HTTPClient *http.Client
}

func (i *BigDataCloud) Init() error {
	i.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("ipResolver.noVerify")},
		},
	}
	return nil
}

func (i *BigDataCloud) GetExternalIP() (net.IP, error) {
	r, err := i.HTTPClient.Get("https://api.bigdatacloud.net/data/client-ip")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		var resp BigDataCloudResponse
		json.NewDecoder(r.Body).Decode(&resp)
		parsedIP := net.ParseIP(resp.IPString)
		if parsedIP == nil {
			return nil, fmt.Errorf("cannot parse IP: %s", resp.IPString)
		}

		logger.Debug("[BigDataCloud] Detected external IP: %v", parsedIP)
		return parsedIP, nil
	}

	return nil, err
}
