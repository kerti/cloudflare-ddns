package resolver

import (
	"fmt"
	"net"

	"github.com/kerti/cloudflare-ddns/logger"
)

const (
	// ResolverBigDataCloud is the resolver identifier for https://api.bigdatacloud.net/data/client-ip
	ResolverBigDataCloud = "bigdatacloud"
	// ResolverICanHazIP is the resolver identifier for http://icanhazip.com
	ResolverICanHazIP = "icanhazip"
	// ResolverIfconfigMe is the resolver identifier for https://ifconfig.me/ip
	ResolverIfconfigMe = "ifconfigme"
	// ResolverIPAPICo is the resolver identifier for https://ipapi.co/ip
	ResolverIPAPICo = "ipapico"
	// ResolverIpify is the resolver identifier for https://ipify.org
	ResolverIpify = "ipify"
	// ResolverMyExternalIP is the resolver identifier for https://myexternalip.com
	ResolverMyExternalIP = "myexternalip"
	// ResolverMyIP is the resolver identifier for https://api.myip.com
	ResolverMyIP = "myip"
	// ResolverWhatIsMyIPAddress is the resolver identifier for http://ipv4bot.whatismyipaddress.com
	ResolverWhatIsMyIPAddress = "whatismyipaddress"
	// ResolverWtfIsMyIP is the resolver identifier for https://wtfismyip.com
	ResolverWtfIsMyIP = "wtfismyip"
)

var (
	// ResolverURLs is the map of available resolvers and their URLs
	ResolverURLs = map[string]string{
		ResolverBigDataCloud:      "https://api.bigdatacloud.net/data/client-ip",
		ResolverICanHazIP:         "http://icanhazip.com",
		ResolverIfconfigMe:        "https://ifconfig.me/ip",
		ResolverIPAPICo:           "https://ipapi.co/ip",
		ResolverIpify:             "https://api.ipify.org?format=text",
		ResolverMyExternalIP:      "https://myexternalip.com/raw",
		ResolverMyIP:              "https://api.myip.com",
		ResolverWhatIsMyIPAddress: "http://ipv4bot.whatismyipaddress.com",
		ResolverWtfIsMyIP:         "https://wtfismyip.com/text",
	}

	// ResolverJSONPaths is the map of available resolvers and their JSON paths
	ResolverJSONPaths = map[string]string{
		ResolverBigDataCloud: "ipString",
		ResolverMyIP:         "ip",
	}
)

// Resolver is the base service
type Resolver interface {
	Init() error
	GetExternalIP() (net.IP, error)
}

// Get instantiates a new instance of the resolver
func Get(key string) (*Resolver, error) {
	var resolver Resolver
	switch key {
	case ResolverBigDataCloud:
		resolver = NewJSONResolver(
			ResolverURLs[ResolverBigDataCloud],
			ResolverJSONPaths[ResolverBigDataCloud])
	case ResolverICanHazIP:
		resolver = NewSimpleResolver(ResolverURLs[ResolverICanHazIP])
	case ResolverIfconfigMe:
		resolver = NewSimpleResolver(ResolverURLs[ResolverIfconfigMe])
	case ResolverIPAPICo:
		resolver = NewSimpleResolver(ResolverURLs[ResolverIPAPICo])
	case ResolverIpify:
		resolver = NewSimpleResolver(ResolverURLs[ResolverIpify])
	case ResolverMyExternalIP:
		resolver = NewSimpleResolver(ResolverURLs[ResolverMyExternalIP])
	case ResolverMyIP:
		resolver = NewJSONResolver(
			ResolverURLs[ResolverMyIP],
			ResolverJSONPaths[ResolverMyIP])
	case ResolverWhatIsMyIPAddress:
		resolver = NewSimpleResolver(ResolverURLs[ResolverWhatIsMyIPAddress])
	case ResolverWtfIsMyIP:
		resolver = NewSimpleResolver(ResolverURLs[ResolverWtfIsMyIP])
	default:
		err := fmt.Errorf("unsupported resolver: %s", key)
		logger.Error(err.Error())
		return nil, err
	}

	err := resolver.Init()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.Debug("Resolver [%s] initiated...", key)
	return &resolver, nil
}

// GetAll instantiate all resolvers
func GetAll() ([]*Resolver, error) {
	resolvers := make([]*Resolver, 0)
	for resKey := range ResolverURLs {
		instantiated, err := Get(resKey)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		resolvers = append(resolvers, instantiated)
	}
	return resolvers, nil
}
