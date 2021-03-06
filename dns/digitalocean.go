package dns

import (
	"context"
	"errors"

	"github.com/digitalocean/godo"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const digitalOceanName = "digitalocean"

// DigitalOceanDNS instance
type DigitalOceanDNS struct {
	name   Name
	token  string
	auth   doAuth
	record string
	log    *logrus.Entry
}

// NewDigitalOceanDNS is DigitalOceanDNS constructor
func NewDigitalOceanDNS(config ProviderConfig) (*DigitalOceanDNS, error) {
	d := &DigitalOceanDNS{}
	d.name = digitalOceanName
	var ok bool
	d.token, ok = config["token"]
	if !ok {
		return d, errors.New("token missing from DigitalOcean provider")
	}
	d.record, ok = config["record"]
	if !ok {
		return d, errors.New(("record missing from DigitalOcean provider"))
	}

	d.log = log.WithFields(logrus.Fields{"dns": "do"})
	return d, nil
}

// GetName returns name identifier
func (d *DigitalOceanDNS) GetName() Name {
	return d.name
}

// SyncARecord sets an A record to the given IPv4 address
func (d *DigitalOceanDNS) SyncARecord(ipv4Address string) error {
	return d.SyncRecord("A", ipv4Address)
}

// SyncAAAARecord sets an AAAA record to the given IPv6 address
func (d *DigitalOceanDNS) SyncAAAARecord(ipv6Address string) error {
	return d.SyncRecord("AAAA", ipv6Address)
}

// SyncRecord sets the given record to match ipAddress
func (d *DigitalOceanDNS) SyncRecord(recordType string, ipAddress string) error {
	// Authenticate with DigitalOcean
	d.auth = newDoAuth(d.token)
	domainName, recordName := SplitDomainRecord(d.record)
	d.log.Debugf("do: searching for domain=%s record=%s", domainName, recordName)

	// First get a list of domain records
	records, err := d.getDomainRecords(domainName)
	if err != nil {
		d.log.Errorf("do: could not get list of domain records")
		d.log.Errorf("do: err=%s", err)
		return err
	}

	// Now we need to find which domain record matches ours
	d.log.Debugf("do: %d records found", len(records))
	matchingIdx := -1
	for idx, record := range records {
		if record.Type == recordType {
			d.log.Debugf("do: record=%s", record)
			if record.Name == recordName {
				matchingIdx = idx
				break
			}
		}
	}
	if matchingIdx < 0 {
		d.log.Infof("do: no matching record found, will attempt to create")
		_, err = d.createDomainRecord(domainName, recordName, recordType, ipAddress)
		if err != nil {
			d.log.Errorf("do: could not create a new domain record")
			d.log.Errorf("do: err=%s", err)
			return err
		}
		d.log.Infof("do: new record successfully created")
		return nil
	}
	d.log.Debugf("do: found matching record id=%d ip=%s",
		records[matchingIdx].ID, records[matchingIdx].Data)
	if ipAddress == records[matchingIdx].Data {
		d.log.Infof("do: record does not need to be updated")
		return nil
	}

	// Else, we need to update the domain record
	editRequest := &godo.DomainRecordEditRequest{
		Type: recordType,
		Data: ipAddress,
	}
	_, _, err = d.auth.Client.Domains.EditRecord(d.auth.Ctx, domainName, records[matchingIdx].ID, editRequest)
	if err != nil {
		d.log.Errorf("do: could not update domain record domain=%s id=%d",
			domainName, records[matchingIdx].ID)
		d.log.Errorf("do: err=%s", err)
		return err
	}
	d.log.Infof("do: record successfully updated")

	return nil
}

func (d *DigitalOceanDNS) getDomainRecords(domain string) ([]godo.DomainRecord, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 1000,
	}

	records, _, err := d.auth.Client.Domains.Records(d.auth.Ctx, domain, opt)
	return records, err
}

func (d *DigitalOceanDNS) createDomainRecord(domainName string, recordName string,
	recordType string, ipAddress string) (*godo.DomainRecord, error) {
	createRequest := &godo.DomainRecordEditRequest{
		Type: recordType,
		Name: recordName,
		Data: ipAddress,
	}

	record, _, err := d.auth.Client.Domains.CreateRecord(d.auth.Ctx, domainName, createRequest)
	return record, err
}

//doTokenSource is a oauth2 helper struct
type doTokenSource struct {
	AccessToken string
}

//Token returns an oauth2 token
func (t *doTokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type doAuth struct {
	Client *godo.Client
	Ctx    context.Context
}

func newDoAuth(apiToken string) doAuth {
	token := &doTokenSource{
		AccessToken: apiToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, token)

	var auth = doAuth{
		Client: godo.NewClient(oauthClient),
		Ctx:    context.TODO(),
	}
	return auth
}
