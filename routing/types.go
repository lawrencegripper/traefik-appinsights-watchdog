package routing

import "time"

// RequestContext is used to hold the nonce and other configuration information during a routing check
type RequestContext struct {
	Nonce             string
	StartTime         time.Time
	Port              int
	TraefikServiceURL string
	FabricURI         string
	BackendURL        string
	InstanceID        string
}
