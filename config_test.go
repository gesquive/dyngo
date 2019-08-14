package main

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestIpCheck(t *testing.T) {
	str := []byte(
		`ip_check:
    ipv4: true
    ipv4_urls:
    - "http://ipv4-1.net"
    - "http://ipv4-2.net"
    ipv6: true
    ipv6_urls:
    - "http://ipv6-1.net"
    - "http://ipv6-2.net"
    - "http://ipv6-3.net"
`)

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(str))
	assert.NoError(t, err, "error reading conf")

	assert.True(t, viper.GetBool("ip_check.ipv4"))
	assert.Len(t, viper.GetStringSlice("ip_check.ipv4_urls"), 2, "IPv4Url count does not match")
	assert.True(t, viper.GetBool("ip_check.ipv6"))
	assert.Len(t, viper.GetStringSlice("ip_check.ipv6_urls"), 3, "IPv6Url count does not match")
}

func TestDefaultIpCheck(t *testing.T) {
	str := []byte(`ip_check: {}`)

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(str))
	assert.NoError(t, err, "error reading conf")

	assert.True(t, viper.GetBool("ip_check.ipv4"))
	assert.Len(t, viper.GetStringSlice("ip_check.ipv4_urls"), 0, "IPv4Url count does not match")
	assert.True(t, viper.GetBool("ip_check.ipv6"))
	assert.Len(t, viper.GetStringSlice("ip_check.ipv6_urls"), 0, "IPv6Url count does not match")
}
