package resolver

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

// NewSimpleResolver instantiates a new simple resolver given a URL
func NewSimpleResolver(url string) Resolver {
	return &Simple{URL: url}
}

// Simple is the simplest implementation of external IP resolver, relying on text-based respnse
type Simple struct {
	URL        string
	HTTPClient *http.Client
}

// Init initializes the resolver
func (s *Simple) Init() error {
	s.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("ipResolver.noVerify")},
		},
	}
	return nil
}

// GetExternalIP invokes the URL and fetches the external IP returned
func (s *Simple) GetExternalIP() (net.IP, error) {
	r, err := s.HTTPClient.Get(s.URL)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		return s.readIP(r)
	}

	err = fmt.Errorf("provider responded with HTTP/%v", r.StatusCode)
	logger.Error(err.Error())
	return nil, err
}

func (s *Simple) readIP(r *http.Response) (net.IP, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	reg, err := regexp.Compile("[^0-9\\.]+")
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	ipString := reg.ReplaceAllString(string(bodyBytes), "")
	parsedIP := net.ParseIP(ipString)
	if parsedIP == nil {
		err = fmt.Errorf("cannot parse IP: %v", ipString)
		logger.Error(err.Error())
		return nil, err
	}

	logger.Debug("[RSLV-SIMP] [%s] Detected external IP: %v", s.URL, parsedIP)
	return parsedIP, nil
}
