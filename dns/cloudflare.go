package dns

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

const cloudflareName = "cloudflare"

// CloudflareDNS instance
type CloudflareDNS struct {
	name   Name
	path   string
	auth   doAuth
	record string
	args   string
	log    *logrus.Entry
}

// NewCloudflareDNS is CloudflareDNS constructor
func NewCloudflareDNS(config ProviderConfig) (*CloudflareDNS, error) {
	d := &CloudflareDNS{}
	d.name = cloudflareName
	var ok bool
	d.path, ok = config["path"]
	if !ok {
		return d, errors.New("path missing from Custom Script provider")
	}
	d.record, ok = config["domain"]
	if !ok {
		return d, errors.New(("domain missing from Custom Script provider"))
	}
	d.args = config["args"]

	d.log = log.WithFields(logrus.Fields{"dns": "cfl"})
	return d, nil
}

// GetName returns name identifier
func (c *CloudflareDNS) GetName() Name {
	return c.name
}

// SyncARecord sets an A record to the given IPv4 address
func (c *CloudflareDNS) SyncARecord(ipv4Address string) error {
	return c.SyncRecord("A", ipv4Address)
}

// SyncAAAARecord sets an AAAA record to the given IPv6 address
func (c *CloudflareDNS) SyncAAAARecord(ipv6Address string) error {
	return c.SyncRecord("AAAA", ipv6Address)
}

// SyncRecord sets the given record to match ipAddress
func (c *CloudflareDNS) SyncRecord(recordType string, ipAddress string) error {
	// Run the script
	cmd := exec.Command(c.path, recordType, ipAddress, c.args)
	log.Debugf("Running cmd %v", cmd.Args)
	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		log.Errorf("custom script '%s' returned with errors", c.path)
		if out.Len() > 0 {
			log.Errorf("stderr: %s", strings.TrimSpace(out.String()))
		}
		log.Error(err)
		return err
	}

	return nil
}
