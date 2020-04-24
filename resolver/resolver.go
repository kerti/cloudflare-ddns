package resolver

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"

	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

// Resolver represents a generic resolver
type Resolver struct {
	Name       string
	Type       string
	URL        string
	JSONPath   string
	HTTPClient *http.Client
}

// Get constructs all generic resolvers
func Get() ([]Resolver, error) {
	resolvers := make([]Resolver, 0)
	viper.UnmarshalKey("resolver.list", &resolvers)

	result := make([]Resolver, 0)
	for _, res := range resolvers {
		res.Init()
		result = append(result, res)
	}

	return result, nil
}

// Init initializes the resolver
func (r *Resolver) Init() {
	logger.Debug("[RESOLVER] Initializing for [%s]", r.URL)
	r.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("resolver.noVerify")},
		},
	}
}

// GetExternalIP invokes the URL and fetches the external IP returned
func (r *Resolver) GetExternalIP() (net.IP, error) {
	response, err := r.HTTPClient.Get(r.URL)
	if err != nil {
		return nil, fmt.Errorf("failed resolving external IP: %v", err.Error())
	}

	if response == nil {
		err = fmt.Errorf("response is nil")
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		switch r.Type {
		case "text":
			return r.readIPText(*response)
		case "json":
			return r.readIPJSON(*response)
		default:
			err = fmt.Errorf("unsupported resolver type: [%v]", r.Type)
			return nil, err
		}
	}

	err = fmt.Errorf("provider responded with HTTP/%v", response.StatusCode)
	return nil, err
}

func (r *Resolver) readIPText(response http.Response) (net.IP, error) {
	bodyBytes, err := r.getResponseBodyBytes(response)
	if err != nil {
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
		err = fmt.Errorf("cannot parse IP: [%v]", string(bodyBytes))
		logger.Error(err.Error())
		return nil, err
	}

	logger.Debug("[RESOLVER] [%s] Detected external IP: %v", r.URL, parsedIP)
	return parsedIP, nil
}

func (r *Resolver) getResponseBodyBytes(response http.Response) ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return bodyBytes, nil
}

func (r *Resolver) readIPJSON(response http.Response) (net.IP, error) {
	kvMap, err := r.getResponseJSONMap(response)
	if err != nil {
		return nil, err
	}

	ipString, err := r.findIPString(kvMap)
	if err != nil {
		return nil, err
	}

	parsedIP := net.ParseIP(ipString)
	if parsedIP == nil {
		err = fmt.Errorf("cannot parse IP: [%v]", ipString)
		logger.Error(err.Error())
		return nil, err
	}

	logger.Debug("[RESOLVER] [%s] Detected external IP: %v", r.URL, parsedIP)
	return parsedIP, nil
}

func (r *Resolver) getResponseJSONMap(response http.Response) (map[string]interface{}, error) {
	bodyBytes, err := r.getResponseBodyBytes(response)
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

func (r *Resolver) findIPString(kvMap map[string]interface{}) (string, error) {
	ipObj := kvMap[r.JSONPath]
	if ipObj == nil {
		err := fmt.Errorf("IP address not found at path [%s]", r.JSONPath)
		logger.Error(err.Error())
		return "", err
	}

	ipString, ok := ipObj.(string)
	if !ok {
		err := fmt.Errorf("cannot convert value at path [%s] to string: %v", r.JSONPath, ipObj)
		logger.Error(err.Error())
		return "", err
	}

	return ipString, nil
}
