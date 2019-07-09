package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gesquive/dyngo/dns"
	"github.com/spf13/viper"
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
	currentIP, err := getPublicIPv4Address()
	if err != nil {
		log.Errorf("sync: could not get public ip address")
		log.Errorf("sync: err=%s", err)
		return
	}

	// Second, update DigitalOcean record
	dodns := dns.NewDigitalOceanDNS(token, log)
	dodns.SyncRecord(domainRecord, currentIP)
}

func getPublicIPv4Address() (ipAddress string, err error) {
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
		if ip := net.ParseIP(ipAddress); ip == nil || ip.To4() == nil {
			log.Errorf("ipchk: response is not a valid IPv4 address. response='%s'",
				ipAddress)
			ipAddress = ""
			continue
		}
		gotIP = true
	}
	if !gotIP {
		err = fmt.Errorf("ran out of attempts to get IP address")
	}

	log.Infof("ipchk: got public IP address=%s", ipAddress)
	return ipAddress, err
}
