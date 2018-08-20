package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"context"

	"github.com/digitalocean/godo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// RunService runs as a service
func RunService(token string, domainRecord string, syncInterval time.Duration) {
	log.Infof("service: run as service every %s", syncInterval)
	for {
		go SyncDomain(token, domainRecord)
		time.Sleep(syncInterval)
	}
}

// RunSync syncs your public IP with the given domain
func RunSync(token string, domainRecord string) {
	log.Infof("update: Updating domain record=%s", domainRecord)
	SyncDomain(token, domainRecord)
}

//SyncDomain sets a domain record point to our public IP address
func SyncDomain(token string, domainRecord string) {
	// First get our public IP
	currentIP, err := getPublicIPAddress()
	if err != nil {
		log.Errorf("sync: could not get public ip address")
		log.Errorf("sync: err=%s", err)
		return
	}
	log.Infof("sync: got public IP address=%s", currentIP)

	// Second, authenticate with DigitalOcean
	do := NewDoAuth(token)

	// Next get a list of domain records
	domainName, recordName := splitDomainRecord(domainRecord)
	log.Debugf("sync: searching for domain=%s record=%s", domainName, recordName)
	records, err := getDomainRecords(&do, domainName)
	if err != nil {
		log.Errorf("sync: could not get list of domain records")
		log.Errorf("sync: err=%s", err)
		return
	}

	// Now we need to find which domain record matches ours
	log.Debugf("sync: %d records found", len(records))
	matchingIdx := -1
	for idx, record := range records {
		if record.Type == "A" {
			log.Debugf("sync: record=%s", record)
			if record.Name == recordName {
				matchingIdx = idx
				break
			}
		}
	}
	if matchingIdx < 0 {
		log.Infof("sync: no matching record found, will attempt to create")
		_, err = createDomainRecord(&do, domainName, recordName, currentIP)
		if err != nil {
			log.Errorf("sync: could not create a new domain record")
			log.Errorf("sync: err=%s", err)
		} else {
			log.Infof("sync: new record successfully created")
		}
		return
	}
	log.Debugf("sync: found matching record id=%d ip=%s",
		records[matchingIdx].ID, records[matchingIdx].Data)

	if currentIP == records[matchingIdx].Data {
		log.Infof("sync: record does not need to be updated")
		return
	}

	// Else, we need to update the domain record
	editRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Data: currentIP,
	}
	_, _, err = do.Client.Domains.EditRecord(do.Ctx, domainName, records[matchingIdx].ID, editRequest)
	if err != nil {
		log.Errorf("sync: could not update domain record domain=%s id=%d",
			domainName, records[matchingIdx].ID)
		log.Errorf("sync: err=%s", err)
		return
	}
	log.Infof("sync: record successfully updated")
}

func splitDomainRecord(domainRecord string) (domain string, record string) {
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

func getPublicIPAddress() (ipAddress string, err error) {
	maxAttempts := 3
	ipCheckServices := viper.GetStringSlice("url_list")
	rand.Seed(time.Now().Unix())
	gotIP := false

	for i := 0; i < maxAttempts && !gotIP; i++ {
		victim := rand.Intn(len(ipCheckServices))
		url := ipCheckServices[victim]
		log.Infof("ipchk: using '%s' for ip check", url)

		response, herr := http.Get(url)
		if herr != nil {
			log.Errorf("ipchk: Failed to get ip from '%s'", url)
			log.Errorf("ipchk: err=%s", herr)
			continue
		}
		defer response.Body.Close()
		body, berr := ioutil.ReadAll(response.Body)
		if berr != nil {
			log.Errorf("ipchk: Could not read response from '%s'", url)
			log.Errorf("ipchk: err=%s", berr)
			continue
		}
		ipAddress = strings.TrimSpace(string(body))
		if net.ParseIP(ipAddress) == nil {
			log.Errorf("ipchk: response is not a valid IP address. response='%s'",
				ipAddress)
			ipAddress = ""
			continue
		}
		gotIP = true
	}
	if !gotIP {
		err = fmt.Errorf("ran out of attempts to get IP address")
	}
	return ipAddress, err
}

//TokenSource is a oauth2 helper struct
type TokenSource struct {
	AccessToken string
}

//Token returns an oauth2 token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type DoAuth struct {
	Client *godo.Client
	Ctx    context.Context
}

func NewDoAuth(apiToken string) DoAuth {
	token := &TokenSource{
		AccessToken: apiToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, token)

	var auth = DoAuth{
		Client: godo.NewClient(oauthClient),
		Ctx:    context.TODO(),
	}
	return auth
}

func getDomainRecords(do *DoAuth, domain string) ([]godo.DomainRecord, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 1000,
	}

	records, _, err := do.Client.Domains.Records(do.Ctx, domain, opt)
	return records, err
}

func createDomainRecord(do *DoAuth, domainName string, recordName string, ipAddress string) (*godo.DomainRecord, error) {
	createRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: recordName,
		Data: ipAddress,
	}

	record, _, err := do.Client.Domains.CreateRecord(do.Ctx, domainName, createRequest)
	return record, err
}
