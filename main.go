package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
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
				panic("Application insights key is required use '--appinsightskey=key' use '-h' to see help")
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
	client := newTelemetryClient(config)

	go routing.StartCheck(config, healthChan)
	go health.StartCheck(config, healthChan)

	for {
		event := <-healthChan
		publishToAppInsights(client, event, config)
		if config.Debug {
			fmt.Println(prettyPrintStruct(event))
		}
	}
}

func newTelemetryClient(config types.Configuration) appinsights.TelemetryClient {
	telemetryClient := appinsights.NewTelemetryClient(config.AppInsightsKey)
	telemetryClient.Context().Cloud().SetRoleName("traefik-appinsights-watchdog")
	telemetryClient.Context().Cloud().SetRoleInstance(config.InstanceID)
	return telemetryClient
}

func publishToAppInsights(client appinsights.TelemetryClient, event types.StatsEvent, config types.Configuration) {
	telemetry := appinsights.NewEventTelemetry(config.InstanceID)
	telemetry.SetProperty("sourceTime", event.SourceTime.String())
	telemetry.SetProperty("source", event.Source)
	telemetry.SetProperty("instanceID", config.InstanceID) //Duplicated for discoverability
	telemetry.SetProperty("isSuccess", strconv.FormatBool(event.IsSuccess))
	telemetry.SetProperty("requestDurationInNs", strconv.FormatInt(event.RequestDuration.Nanoseconds(), 10))
	if !event.IsSuccess {
		telemetry.SetProperty("errorDetails", event.ErrorDetails)
	}
	for k, v := range event.Data {
		subMap, ok := v.(map[string]interface{})
		if !ok {
			s := fmt.Sprint(v)
			telemetry.SetProperty(k, s)
			continue
		}
		if subMap == nil {
			continue
		}
		for subk, subv := range subMap {
			subs := fmt.Sprint(subv)
			var buffer bytes.Buffer
			buffer.WriteString(k)
			buffer.WriteString(".")
			buffer.WriteString(subk)
			compk := buffer.String()
			telemetry.SetProperty(compk, subs)
		}
	}
	client.TrackEventTelemetry(telemetry)
}
