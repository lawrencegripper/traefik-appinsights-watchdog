package types

// Configuration provides the configuration for starting the watchdog
type Configuration struct {
	Debug                  bool   `description:"Set to true for additional output in the console"`
	InstanceID             string `description:"The name to report for the instance"`
	AppInsightsKey         string `description:"The application insights instrumentation key"`
	TraefikBackendName     string `description:"This is the name Traefik backend name of the watchdog test server. In SF this will be the fabricURI"`
	WatchdogTestServerPort int    `description:"Port which the simulated backend runs on"`
	WatchdogTraefikURL     string `description:"The url traefik will use to route requests to the watchdog"`
	TraefikHealthEndpoint  string `description:"The traeifk health endpoint http://localhost:port/health"`
	PollIntervalSec        int    `description:"The time waited between requests to the health endpoint"`
	AllowInvalidCert       bool   `description:"Allow invalid certificates when performing routing checks on localhost"`
}
