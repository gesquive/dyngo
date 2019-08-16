package main

import (
	"github.com/gesquive/dyngo/dns"
	"github.com/spf13/viper"
)

// type dnsProvidersList []map[string]string
type dnsProvidersList []dns.Provider

func getDNSProviders() (dnsProvidersList, error) {
	if ! viper.IsSet("dns_providers") {
		var dnsPrv dnsProvidersList
		return dnsPrv, nil
	}

	var dnsConfigs []map[string]string
	err := viper.UnmarshalKey("dns_providers", &dnsConfigs)
	if err != nil {
		return dnsProvidersList{}, err
	}

	dnsPrv := make(dnsProvidersList, len(dnsConfigs))
	for i, providerConfig := range dnsConfigs {
		dnsProvider, err := dns.GetDNSProvider(providerConfig)
		if err != nil {
			return nil, err
		}
		dnsPrv[i] = dnsProvider
	}

	return dnsPrv, nil
}
