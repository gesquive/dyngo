package cmd

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Run runs
func Run() {
	log.Infoln("Running")
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
	client := doAuth(token)

	// Next get a list of domain records
	domain := getDomain(domainRecord)
	records, err := getDomainRecords(client, domain)
	if err != nil {
		log.Errorf("sync: could not get list of domain records")
		log.Errorf("sync: err=%s", err)
		return
	}

	// Now we need to find which domain record matches ours
	log.Debugf("sync: %d records found", len(records))
	matchingIdx := -1
	for idx, record := range records {
		log.Debugf("sync: record=%s", record)
		if record.Name == domainRecord {
			matchingIdx = idx
			break
		}
	}
	if matchingIdx < 0 {
		log.Infof("sync: no matching record found, will attempt to create")
		//TODO: create a new record
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
	_, _, err = client.Domains.EditRecord(domain, records[matchingIdx].ID,
		editRequest)
	if err != nil {
		log.Errorf("sync: could not update domain record domain=%s id=%d",
			domain, records[matchingIdx].ID)
		log.Errorf("sync: err=%s", err)
		return
	}
	log.Infof("sync: record successfully updated")
}

func getDomain(domainRecord string) (domain string) {
	domainParts := strings.Split(domainRecord, ".")
	if len(domainParts) > 2 {
		// sub.domain.net => domain.net
		domain = strings.Join(domainParts[len(domainParts)-2:], ".")
	} else {
		domain = domainRecord
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

func doAuth(apiToken string) (client *godo.Client) {
	token := &TokenSource{
		AccessToken: apiToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, token)
	client = godo.NewClient(oauthClient)
	return
}

func getDomainRecords(client *godo.Client, domain string) ([]godo.DomainRecord, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 1000,
	}

	records, _, err := client.Domains.Records(domain, opt)
	return records, err
}
