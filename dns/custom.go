package dns

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

const customScriptName = "custom"

// CustomScriptDNS instance
type CustomScriptDNS struct {
	name   Name
	path   string
	auth   doAuth
	record string
	args   string
	log    *logrus.Entry
}

// NewCustomScriptDNS is CustomScriptDNS constructor
func NewCustomScriptDNS(config ProviderConfig) (*CustomScriptDNS, error) {
	c := &CustomScriptDNS{}
	c.name = customScriptName
	var ok bool
	c.path, ok = config["path"]
	if !ok {
		return c, errors.New("path missing from Custom Script provider")
	}
	c.record, ok = config["record"]
	if !ok {
		return c, errors.New(("domain missing from Custom Script provider"))
	}
	c.args = config["args"]

	c.log = log.WithFields(logrus.Fields{"dns": "cus"})
	return c, nil
}

// GetName returns name identifier
func (c *CustomScriptDNS) GetName() Name {
	return c.name
}

// SyncARecord sets an A record to the given IPv4 address
func (c *CustomScriptDNS) SyncARecord(ipv4Address string) error {
	return c.SyncRecord("A", ipv4Address)
}

// SyncAAAARecord sets an AAAA record to the given IPv6 address
func (c *CustomScriptDNS) SyncAAAARecord(ipv6Address string) error {
	return c.SyncRecord("AAAA", ipv6Address)
}

// SyncRecord sets the given record to match ipAddress
func (c *CustomScriptDNS) SyncRecord(recordType string, ipAddress string) error {
	// Run the script
	cmd := exec.Command(c.path, c.record, recordType, ipAddress, c.args)
	log.Debugf("cus: running cmd %v", cmd.Args)
	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		log.Errorf("cus: script '%s' returned with errors", c.path)
		if out.Len() > 0 {
			log.Errorf("stderr: %s", strings.TrimSpace(out.String()))
		}
		log.Error(err)
		return err
	}

	return nil
}
