package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// RunService runs as a service
func RunService(dns dnsProvidersList, syncInterval time.Duration) {
	log.Infof("service: run as service every %s", syncInterval)
	for {
		go SyncDomain(dns)
		time.Sleep(syncInterval)
	}
}

// RunSync syncs your public IP with the given domain
func RunSync(dns dnsProvidersList) {
	log.Infof("update: Updating record for %d providers", len(dns))
	SyncDomain(dns)
}

//SyncDomain sets a domain record point to our public IP address
func SyncDomain(dnsProviders dnsProvidersList) {
	setIPv4 := viper.GetBool("ip_check.ipv4")
	setIPv6 := viper.GetBool("ip_check.ipv6")
	if !setIPv4 && !setIPv6 {
		log.Warnf("All IP checks are turned off, no sync")
	}
	if setIPv4 {
		// First get our public IP
		currentIP, err := getPublicIPv4Address()
		if err != nil {
			log.Errorf("sync: could not get public ipv4 address")
			log.Errorf("sync: err=%s", err)
			return
		}

		// Second, update all DNS providers
		for _, provider := range dnsProviders {
			provider.SyncARecord(currentIP)
		}
	}

	if setIPv6 {
		// First get our public IP
		currentIP, err := getPublicIPv6Address()
		if err != nil {
			log.Errorf("sync: could not get public ipv6 address")
			log.Errorf("sync: err=%s", err)
			return
		}

		// Second, update all DNS providers
		for _, provider := range dnsProviders {
			provider.SyncAAAARecord(currentIP)
		}
	}
}

func getPublicIPv4Address() (ipAddress string, err error) {
	maxAttempts := 3
	ipCheckServices := viper.GetStringSlice("ip_check.ipv4_urls")
	rand.Seed(time.Now().Unix())
	gotIP := false

	for i := 0; i < maxAttempts && !gotIP; i++ {
		victim := rand.Intn(len(ipCheckServices))
		url := ipCheckServices[victim]
		log.Infof("ipchk4: using '%s' for ip check", url)

		response, herr := http.Get(url)
		if herr != nil {
			log.Errorf("ipchk4: Failed to get ip from '%s'", url)
			log.Errorf("ipchk4: err=%s", herr)
			continue
		}
		defer response.Body.Close()
		body, berr := ioutil.ReadAll(response.Body)
		if berr != nil {
			log.Errorf("ipchk4: Could not read response from '%s'", url)
			log.Errorf("ipchk4: err=%s", berr)
			continue
		}
		ipAddress = strings.TrimSpace(string(body))
		if ip := net.ParseIP(ipAddress); ip == nil || ip.To4() == nil {
			log.Errorf("ipchk4: response is not a valid IPv4 address. response='%s'",
				ipAddress)
			ipAddress = ""
			continue
		}
		gotIP = true
	}
	if !gotIP {
		err = fmt.Errorf("ipchk4: ran out of attempts to get IP address")
	}

	log.Infof("ipchk: got public IP address=%s", ipAddress)
	return ipAddress, err
}

func getPublicIPv6Address() (ipAddress string, err error) {
	maxAttempts := 3
	ipCheckServices := viper.GetStringSlice("ip_check.ipv6_urls")
	rand.Seed(time.Now().Unix())
	gotIP := false

	for i := 0; i < maxAttempts && !gotIP; i++ {
		victim := rand.Intn(len(ipCheckServices))
		url := ipCheckServices[victim]
		log.Infof("ipchk6: using '%s' for ip check", url)

		response, herr := http.Get(url)
		if herr != nil {
			log.Errorf("ipchk6: Failed to get ip from '%s'", url)
			log.Errorf("ipchk6: err=%s", herr)
			continue
		}
		defer response.Body.Close()
		body, berr := ioutil.ReadAll(response.Body)
		if berr != nil {
			log.Errorf("ipchk6: Could not read response from '%s'", url)
			log.Errorf("ipchk6: err=%s", berr)
			continue
		}
		ipAddress = strings.TrimSpace(string(body))
		if ip := net.ParseIP(ipAddress); ip == nil || ip.To16() == nil {
			log.Errorf("ipchk6: response is not a valid IPv6 address. response='%s'",
				ipAddress)
			ipAddress = ""
			continue
		}
		gotIP = true
	}
	if !gotIP {
		err = fmt.Errorf("ipchk6: ran out of attempts to get IP address")
	}

	log.Infof("ipchk6: got public IP address=%s", ipAddress)
	return ipAddress, err
}
