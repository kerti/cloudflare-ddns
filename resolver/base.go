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
	// Resolvers is the list of available resolvers
	Resolvers = map[string]string{
		ResolverBigDataCloud:      ResolverBigDataCloud,
		ResolverICanHazIP:         ResolverICanHazIP,
		ResolverIfconfigMe:        ResolverIfconfigMe,
		ResolverIPAPICo:           ResolverIPAPICo,
		ResolverIpify:             ResolverIpify,
		ResolverMyExternalIP:      ResolverMyExternalIP,
		ResolverMyIP:              ResolverMyIP,
		ResolverWhatIsMyIPAddress: ResolverWhatIsMyIPAddress,
		ResolverWtfIsMyIP:         ResolverWtfIsMyIP,
	}
)

// CheckConfig checks the resolver configuration
func CheckConfig() error {
	logger.Debug("Resolver configuration OK.")
	return nil
}

// Resolver is the base service
type Resolver interface {
	Init() error
	GetExternalIP() (net.IP, error)
}

// Get instantiates a new instance of the resolver
func Get(key string) (Resolver, error) {
	var resolver Resolver
	switch key {
	case ResolverBigDataCloud:
		resolver = new(BigDataCloud)
	case ResolverICanHazIP:
		resolver = new(ICanHazIP)
	case ResolverIfconfigMe:
		resolver = new(IfconfigMe)
	case ResolverIPAPICo:
		resolver = new(IPAPICo)
	case ResolverIpify:
		resolver = new(Ipify)
	case ResolverMyExternalIP:
		resolver = new(MyExternalIP)
	case ResolverMyIP:
		resolver = new(MyIP)
	case ResolverWhatIsMyIPAddress:
		resolver = new(WhatIsMyIPAddress)
	case ResolverWtfIsMyIP:
		resolver = new(WTFIsMyIP)
	default:
		return nil, fmt.Errorf("unsupported resolver: %s", key)
	}

	err := resolver.Init()
	if err != nil {
		return nil, err
	}
	logger.Debug("Resolver [%s] initiated...", key)
	return resolver, nil
}

// GetAll instantiate all resolvers
func GetAll() ([]Resolver, error) {
	resolvers := make([]Resolver, 0)
	for resKey := range Resolvers {
		instantiated, err := Get(resKey)
		if err != nil {
			return nil, err
		}
		resolvers = append(resolvers, instantiated)
	}
	return resolvers, nil
}
