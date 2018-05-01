package main

import (
	"fmt"
	"strconv"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/lawrencegripper/traefik-appinsights-watchdog/types"
)

// NewTelemetryClient initializes an Application Insights Client
func NewTelemetryClient(config types.Configuration) appinsights.TelemetryClient {
	telemetryClient := appinsights.NewTelemetryClient(config.AppInsightsKey)
	telemetryClient.Context().Cloud().SetRoleName("traefik-appinsights-watchdog")
	telemetryClient.Context().Cloud().SetRoleInstance(config.InstanceID)
	return telemetryClient
}

// PublishToAppInsights sends an event
func PublishToAppInsights(client appinsights.TelemetryClient, event types.StatsEvent, config types.Configuration) {
	telemetry := appinsights.NewEventTelemetry(config.InstanceID)
	telemetry.SetProperty("sourceTime", event.SourceTime.String())
	telemetry.SetProperty("source", event.Source)
	telemetry.SetProperty("instanceID", config.InstanceID) //Duplicated for discoverability
	telemetry.SetProperty("isSuccess", strconv.FormatBool(event.IsSuccess))
	telemetry.SetProperty("requestDurationInNs", strconv.FormatInt(event.RequestDuration.Nanoseconds(), 10))
	if !event.IsSuccess {
		telemetry.SetProperty("errorDetails", event.ErrorDetails)
	}

	for key, value := range event.Data {
		subMap, ok := value.(map[string]interface{})
		if !ok {
			s := fmt.Sprint(value)
			telemetry.SetProperty(key, s)
			continue
		}
		if subMap == nil {
			continue
		}
		// If "value" is a map, then generate composite property keys for each of the sub values
		// i.e.
		//   key.subKey1 = value1
		//   key.subKey2 = value2
		// ...
		for subKey, subValue := range subMap {
			subValueStr := fmt.Sprint(subValue)
			compositeKey := key + "." + subKey
			telemetry.SetProperty(compositeKey, subValueStr)
		}
	}
	client.TrackEventTelemetry(telemetry)
}
