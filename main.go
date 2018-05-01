package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/containous/flaeg"
	"github.com/lawrencegripper/traefik-appinsights-watchdog/health"
	"github.com/lawrencegripper/traefik-appinsights-watchdog/routing"
	"github.com/lawrencegripper/traefik-appinsights-watchdog/types"
)

func main() {
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Println("Unable to automatically set instanceid to hostname")
	}
	config := &types.Configuration{
		Debug:                  false,
		PollIntervalSec:        120,
		InstanceID:             hostName,
		WatchdogTestServerPort: 40001,
		TraefikHealthEndpoint:  "http://localhost:8080/health",
		TraefikBackendName:     "fabric:/TraefikType/Watchdog",
		WatchdogTraefikURL:     "http://localhost:80/TraefikType/Watchdog",
	}

	rootCmd := &flaeg.Command{
		Name:                  "start",
		Description:           `Starts the watchdog, checking both the /health endpoint and request routing`,
		Config:                config,
		DefaultPointersConfig: config,
		Run: func() error {
			if config.AppInsightsKey == "" {
				fmt.Println("Application insights key is required use '--appinsightskey=key' use '-h' to see help")
				os.Exit(1)
			}

			fmt.Printf("Running watchdog with config :\n %+v\n", prettyPrintStruct(config))
			startWatchdog(*config)
			return nil
		},
	}

	//init flaeg
	flaeg := flaeg.New(rootCmd, os.Args[1:])

	//run test
	if err := flaeg.Run(); err != nil {
		fmt.Printf("Error %s \n", err.Error())
	}
}

func prettyPrintStruct(item interface{}) string {
	b, _ := json.MarshalIndent(item, "", " ")
	return string(b)
}

func startWatchdog(config types.Configuration) {
	healthChan := make(chan types.StatsEvent)
	client := NewTelemetryClient(config)

	go routing.StartCheck(config, healthChan)
	go health.StartCheck(config, healthChan)

	for {
		event := <-healthChan
		PublishToAppInsights(client, event, config)
		if config.Debug {
			fmt.Println(prettyPrintStruct(event))
		}
	}
}
