package dns

import (
	"strings"

	"github.com/pkg/errors"
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
	cleanName := strings.ToLower(strings.TrimSpace(name))
	switch cleanName {
	case cloudflareName:
		dns, err = NewCloudflareDNS(config)
	case digitalOceanName:
		dns, err = NewDigitalOceanDNS(config)
	case customScriptName:
		dns, err = NewCustomScriptDNS(config)
	default:
		err = errors.Errorf("dns provider name '%s' not recognized", cleanName)
	}
	return
}

// IntializeLogging sets the logger to use in this library
func IntializeLogging(logger *logrus.Logger) {
	log = logger
}
