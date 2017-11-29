package types

type Configuration struct {
	InstanceID string
	AppInsightsKey string
	WatchdogFabricURI string
	SimulatedBackendURL string
	SimulatedBackendPort int
	WatchdogTraefikURL string
	TraefikHealthEndpoint string
	PollIntervalSec int
}