package types

// Configuration provides the configuration for starting the watchdog
type Configuration struct {
	Debug                 bool   `description:"Set to true for additional output in the console"`
	InstanceID            string `description:"The name to report for the instance"`
	AppInsightsKey        string `description:"The application insights instrumentation key"`
	WatchdogFabricURI     string `description:"Fabric URI of the watchdog service will run under"`
	SimulatedBackendPort  int    `description:"Port which the simulated backend runs on"`
	WatchdogTraefikURL    string `description:"The url traefik will use to route requests to the watchdog"`
	TraefikHealthEndpoint string `description:"The traeifk health endpoint http://localhost:port/health"`
	PollIntervalSec       int    `description:"The time waited between requests to the health endpoint"`
}
