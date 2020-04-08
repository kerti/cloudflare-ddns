package resolver

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

// NewJSONResolver instantiates a new JSON resolver given a URL and JSON path
func NewJSONResolver(url string, jsonPath string) Resolver {
	return &JSON{URL: url, JSONPath: jsonPath}
}

// JSON is the JSON implementation
type JSON struct {
	URL        string
	HTTPClient *http.Client
	JSONPath   string
}

// Init initializes the resolver
func (j *JSON) Init() {
	j.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("ipResolver.noVerify")},
		},
	}
}

// GetExternalIP invokes the URL and fetches the external IP returned
func (j *JSON) GetExternalIP() (net.IP, error) {
	r, err := j.HTTPClient.Get(j.URL)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		return j.readIP(r)
	}

	err = fmt.Errorf("provider responded with HTTP/%v", r.StatusCode)
	logger.Error(err.Error())
	return nil, err
}

func (j *JSON) readIP(r *http.Response) (net.IP, error) {
	kvMap, err := j.getResponseJSONMap(r)
	if err != nil {
		return nil, err
	}

	ipString, err := j.findIPString(kvMap)
	if err != nil {
		return nil, err
	}

	parsedIP := net.ParseIP(ipString)
	if parsedIP == nil {
		err = fmt.Errorf("cannot parse IP: %v", ipString)
		logger.Error(err.Error())
		return nil, err
	}

	logger.Debug("[RSLV-JSON] [%s] Detected external IP: %v", j.URL, parsedIP)
	return parsedIP, nil
}

func (j *JSON) getResponseJSONMap(r *http.Response) (map[string]interface{}, error) {
	if r == nil {
		err := fmt.Errorf("response is nil")
		logger.Error(err.Error())
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var kvMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &kvMap)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return kvMap, nil
}

func (j *JSON) findIPString(kvMap map[string]interface{}) (string, error) {
	ipObj := kvMap[j.JSONPath]
	if ipObj == nil {
		err := fmt.Errorf("IP address not found at path [%s]", j.JSONPath)
		logger.Error(err.Error())
		return "", err
	}

	ipString, ok := ipObj.(string)
	if !ok {
		err := fmt.Errorf("cannot convert value at path [%s] to string: %v", j.JSONPath, ipObj)
		logger.Error(err.Error())
		return "", err
	}

	return ipString, nil
}
