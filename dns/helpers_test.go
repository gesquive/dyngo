package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainSplit(t *testing.T) {
	domain, record := SplitDomainRecord("sub.domain.com")
	assert.Equal(t, "sub", record)
	assert.Equal(t, "domain.com", domain)
}

func TestApexDomainSplit(t *testing.T) {
	domain, record := SplitDomainRecord("domain.com")
	assert.Equal(t, "@", record)
	assert.Equal(t, "domain.com", domain)
}

func TestMultiSubDomainSplit(t *testing.T) {
	domain, record := SplitDomainRecord("test.sub.domain.com")
	assert.Equal(t, "test.sub", record)
	assert.Equal(t, "domain.com", domain)
}

func TestBadSplit(t *testing.T) {
	domain, record := SplitDomainRecord(".sub.domain.com")
	assert.Equal(t, ".sub", record)
	assert.Equal(t, "domain.com", domain)
}

func TestBadShortSplit(t *testing.T) {
	domain, record := SplitDomainRecord(".com")
	assert.Equal(t, "@", record)
	assert.Equal(t, ".com", domain)
}
