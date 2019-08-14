package dns

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

const libVersion = "0.1.0"

var log = logrus.New()

// Name is the provider name
type Name string

// ProviderConfig is the generic config format
type ProviderConfig map[string]string

// Provider generic interface
type Provider interface {
	SyncARecord(ipv4Address string) error
	SyncAAAARecord(ipv6Address string) error
	GetName() Name
}

// GetDNSProvider returns a provider from a given config
func GetDNSProvider(config ProviderConfig) (dns Provider, err error) {
	name, ok := config["name"]
	if !ok {
		err = errors.New("config missing provider name")
		return
	}
	switch strings.ToLower(strings.TrimSpace(name)) {
	case digitalOceanName:
		dns, err = NewDigitalOceanDNS(config)
	default:
		err = errors.New("dns provider name not recognized")
	}
	return
}

// IntializeLogging sets the logger to use in this library
func IntializeLogging(logger *logrus.Logger) {
	log = logger
}
