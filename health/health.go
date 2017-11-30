package health

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lawrencegripper/sfTraefikWatchdog/types"
)

// StartCheck poll health endpoint
func StartCheck(config types.Configuration, healthChannel chan<- types.StatsEvent) {
	intervalDuration := time.Second * time.Duration(config.PollIntervalSec)
	for {
		ev := getStatsEvent(config.TraefikHealthEndpoint)
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
