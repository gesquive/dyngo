package dns

const libVersion = "0.1.0"


// Name is the provider name
type Name string

// A DNS connection
type DNSProvider interface {
	SyncRecord(record string, ipAddress string) error
	GetName() Name
}
