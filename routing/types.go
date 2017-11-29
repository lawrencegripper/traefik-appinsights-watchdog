package routing

import "time"

type RequestContext struct {
	Nonce             string
	StartTime         time.Time
	Port              int
	TraefikServiceURL string
	FabricURI         string
	BackendURL        string
}