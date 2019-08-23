package dns

import (
	"errors"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

const cloudflareName = "cloudflare"

// CloudflareDNS instance
type CloudflareDNS struct {
	name   Name
	token  string
	api    *cloudflare.API
	record string
	args   string
	log    *logrus.Entry
}

// NewCloudflareDNS is CloudflareDNS constructor
func NewCloudflareDNS(config ProviderConfig) (*CloudflareDNS, error) {
	c := &CloudflareDNS{}
	c.name = cloudflareName
	var ok bool
	c.token, ok = config["token"]
	if !ok {
		return c, errors.New("token missing from Cloudflare provider")
	}
	c.record, ok = config["record"]
	if !ok {
		return c, errors.New(("record missing from Cloudflare provider"))
	}

	c.log = log.WithFields(logrus.Fields{"dns": "cfl"})
	return c, nil
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
	// Authenticate with Cloudflare
	var err error
	c.api, err = cloudflare.NewWithAPIToken(c.token)
	if err != nil {
		c.log.Errorf("cfl: could not log in: %v", err)
		return err
	}
	domainName, recordName := SplitDomainRecord(c.record)
	c.log.Debugf("cfl: searching for domain=%s record=%s", domainName, recordName)

	// First get a list of records that match
	zoneID, err := c.api.ZoneIDByName(domainName)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"domain": domainName,
			"err":    err,
		}).Errorf("cfl: could not find the domain")
		return err
	}
	records, err := c.api.DNSRecords(zoneID, cloudflare.DNSRecord{
		Type: recordType,
		Name: c.record,
	})
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("cfl: could not get a list of records")
		return err
	}
	c.log.Debugf("cfl: %d matching records found", len(records))
	if len(records) > 1 {
		c.log.Errorf("cfl: Found %d matching records, will not update a round robin record set", len(records))
		return errors.New("Found more then one matching record")
	} else if len(records) == 0 {
		c.log.Infof("cfl: no matching record found, will attempt to create")
		_, err := c.createDomainRecord(zoneID, recordType, ipAddress)
		if err != nil {
			c.log.WithFields(logrus.Fields{
				"domain": domainName,
				"err":    err,
			}).Errorf("cfl: could not create a new domain record")
			return err
		} else {
			c.log.Infof("cfl: new record suceessfully created")
		}
		return nil
	}

	// We have one matching record, check if we need to update
	record := records[0]
	c.log.WithFields(logrus.Fields{
		"id": record.ID,
		"ip": record.Content,
	}).Debugf("cfl: found matching record")
	if ipAddress == records[0].Content {
		c.log.Infof("cfl: record does not need to be updated")
		return nil
	}

	// Else, we need to update the domain record
	c.log.WithFields(logrus.Fields{
		"id": record.ID,
		"ip": record.Content,
	}).Infof("cfl: updating record")
	record.Content = ipAddress
	err = c.api.UpdateDNSRecord(zoneID, record.ID, record)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"domain": c.record,
			"id":     record.ID,
			"err":    err,
		}).Errorf("cfl: could not update domain record")
		return err
	}
	c.log.Infof("cfl: record successfully updated")

	return nil
}

func (c *CloudflareDNS) createDomainRecord(zoneID string, recordType string, ipAddress string) (*cloudflare.DNSRecordResponse, error) {

	record := cloudflare.DNSRecord{
		Type:    recordType,
		Name:    c.record,
		Content: ipAddress,
		ZoneID:  zoneID,
	}
	res, err := c.api.CreateDNSRecord(zoneID, record)
	return res, err

}
