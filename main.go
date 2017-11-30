package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/containous/flaeg"
	"github.com/lawrencegripper/sfTraefikWatchdog/health"
	"github.com/lawrencegripper/sfTraefikWatchdog/routing"
	"github.com/lawrencegripper/sfTraefikWatchdog/types"
)

func main() {
	config := &types.Configuration{}

	rootCmd := &flaeg.Command{
		Name:                  "start",
		Description:           `Starts the watchdog, checking both the /health endpoint and request routing`,
		Config:                config,
		DefaultPointersConfig: config,
		Run: func() error {
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

func startWatchdog(config types.Configuration) {
	healthChan := make(chan types.StatsEvent)

	go routing.StartCheck(config, healthChan)
	go health.StartCheck(config, healthChan)

	for {
		event := <-healthChan
		publishToAppInsights(event, config)
		fmt.Println(event.IsSuccess)
	}
}

func publishToAppInsights(event types.StatsEvent, config types.Configuration) {
	client := appinsights.NewTelemetryClient(config.AppInsightsKey)
	telemetry := appinsights.NewEventTelemetry(config.InstanceID)
	telemetry.SetProperty("sourceTime", event.SourceTime.String())
	telemetry.SetProperty("isSuccess", strconv.FormatBool(event.IsSuccess))
	telemetry.SetProperty("requestDuration", event.RequestDuration.String())
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
		if len(subMap) == 0 {
			fmt.Printf("Sub map %v+ is empty\n", k)
		}
		for subk, subv := range subMap {
			subs := fmt.Sprint(subv)
			var buffer bytes.Buffer
			buffer.WriteString(k)
			buffer.WriteString(".")
			buffer.WriteString(subk)
			compk := buffer.String()
			fmt.Println("ss: " + subs)
			fmt.Println("ck: " + compk)
			telemetry.SetProperty(compk, subs)
		}
	}
	client.TrackEventTelemetry(telemetry)
}
