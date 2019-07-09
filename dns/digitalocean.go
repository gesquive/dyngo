package dns

import (
	"context"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const digitalOceanName = "DigitalOcean"

type DigitalOceanDNS struct {
	name  Name
	token string
	auth  doAuth
	log   *logrus.Entry
}

// NewDigitalOceanDNS is DigitalOceanDNS constructor
func NewDigitalOceanDNS(token string, logger *logrus.Logger) *DigitalOceanDNS {
	d := &DigitalOceanDNS{}
	d.name = digitalOceanName
	d.token = token

	d.log = logger.WithFields(logrus.Fields{"dns": "do"})
	return d
}

// GetName returns name identifier
func (d *DigitalOceanDNS) GetName() Name {
	return d.name
}

// SyncRecord sets the given record to match ipAddress
func (d *DigitalOceanDNS) SyncRecord(record string, ipAddress string) error {
	// Authenticate with DigitalOcean
	d.auth = newDoAuth(d.token)
	domainName, recordName := doSplitDomainRecord(record)
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
		if record.Type == "A" {
			d.log.Debugf("do: record=%s", record)
			if record.Name == recordName {
				matchingIdx = idx
				break
			}
		}
	}
	if matchingIdx < 0 {
		d.log.Infof("do: no matching record found, will attempt to create")
		_, err = d.createDomainRecord(domainName, recordName, ipAddress)
		if err != nil {
			d.log.Errorf("do: could not create a new domain record")
			d.log.Errorf("do: err=%s", err)
			return err
		} else {
			d.log.Infof("do: new record successfully created")
		}
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
		Type: "A",
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

func (d *DigitalOceanDNS) createDomainRecord(domainName string, recordName string, ipAddress string) (*godo.DomainRecord, error) {
	createRequest := &godo.DomainRecordEditRequest{
		Type: "A",
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

// ======== Helpers =============
func doSplitDomainRecord(domainRecord string) (domain string, record string) {
	domainParts := strings.Split(domainRecord, ".")
	if len(domainParts) > 2 {
		// sub.domain.net => domain.net
		domain = strings.Join(domainParts[len(domainParts)-2:], ".")
		record = strings.Join(domainParts[:len(domainParts)-2], ".")
	} else {
		domain = domainRecord
		record = "@"
	}
	return
}
