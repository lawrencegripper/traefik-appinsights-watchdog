package health

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/lawrencegripper/sfTraefikWatchdog/types"
)

// StartCheck poll health endpoint
func StartCheck(config types.Configuration) {
	aiClient := appinsights.NewTelemetryClient(config.AppInsightsKey)
	healthChannel := make(chan types.StatsEvent)
	go pollHealthEndpoint(healthChannel, config.TraefikHealthEndpoint, config.PollIntervalSec)
	for healthMsg := range healthChannel {
		fmt.Println(healthMsg)
		publishToAppInsights(aiClient, healthMsg, config.InstanceID)
	}
}
func pollHealthEndpoint(healthChannel chan<- types.StatsEvent, endpoint string, intervalInSec int) {
	intervalDuration := time.Second * time.Duration(intervalInSec)
	for {
		ev := getStatsEvent(endpoint)
		healthChannel <- ev
		time.Sleep(intervalDuration)
	}
}

func getStatsEvent(endpoint string) types.StatsEvent {
	event := types.StatsEvent{
		SourceTime: time.Now(),
		Data:       make(map[string]interface{}),
		IsSuccess:  false,
	}
	start := time.Now()
	resp, err := http.Get(endpoint)
	elapsed := time.Since(start)
	event.RequestDuration = elapsed
	if err != nil || resp.StatusCode != http.StatusOK {
		event.IsSuccess = false
		event.ErrorDetails = err.Error()
		return event
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if err != nil {
		event.IsSuccess = false
		event.ErrorDetails = readErr.Error()
		return event
	}
	var data map[string]interface{}
	jsonErr := json.Unmarshal(body, &data)
	if err != nil {
		event.IsSuccess = false
		event.ErrorDetails = jsonErr.Error()
		return event
	}
	event.IsSuccess = true
	event.Data = data
	return event
}
func publishToAppInsights(client appinsights.TelemetryClient, event types.StatsEvent, telemetryID string) {
	telemetry := appinsights.NewEventTelemetry(telemetryID)
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
