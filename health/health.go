package health

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lawrencegripper/traefik-appinsights-watchdog/types"
)

// StartCheck poll health endpoint
func StartCheck(ctx context.Context, config types.Configuration, healthChannel chan<- types.StatsEvent) {
	intervalDuration := time.Second * time.Duration(config.PollIntervalSec)
	tlsConfig := &tls.Config{}
	if config.AllowInvalidCert {
		tlsConfig.InsecureSkipVerify = true
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ev := getStatsEvent(config.TraefikHealthEndpoint, tlsConfig)
			healthChannel <- ev
			time.Sleep(intervalDuration)
		}
	}
}

func getStatsEvent(endpoint string, tlsConfig *tls.Config) types.StatsEvent {
	event := types.StatsEvent{
		Source:     "HealthCheck",
		SourceTime: time.Now(),
		Data:       make(map[string]interface{}),
		IsSuccess:  false,
	}
	start := time.Now()
	client := &http.Client{
		Timeout: time.Second * 3,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get(endpoint)
	elapsed := time.Since(start)
	event.RequestDuration = elapsed
	if err != nil {
		event.IsSuccess = false
		event.ErrorDetails = err.Error()
		return event
	}
	if resp.StatusCode != http.StatusOK {
		event.IsSuccess = false
		event.ErrorDetails = fmt.Sprintf("Health endpoint returned error code: %v", resp.StatusCode)
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
	if jsonErr != nil {
		event.IsSuccess = false
		event.ErrorDetails = jsonErr.Error()
		return event
	}
	event.IsSuccess = true
	event.Data = data
	return event
}
