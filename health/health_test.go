package health

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lawrencegripper/traefik-appinsights-watchdog/types"
)

func TestHealthRetreiveMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleHealthSuceed))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := types.Configuration{TraefikHealthEndpoint: server.URL + "/health"}
	channel := make(chan types.StatsEvent)

	go StartCheck(ctx, config, channel)

	timeout := time.After(time.Second * 3)

	select {
	case statEvent := <-channel:
		if !statEvent.IsSuccess {
			t.Error("Stats event was a failure")
		}
		t.Log(statEvent)
		return
	case <-timeout:
		t.Error("Timeout occurred")
		return
	}
}

func TestHealthRetreiveMetrics_Invalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleHealthInvalid))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := types.Configuration{TraefikHealthEndpoint: server.URL + "/health"}
	channel := make(chan types.StatsEvent)

	go StartCheck(ctx, config, channel)

	timeout := time.After(time.Second * 3)

	select {
	case statEvent := <-channel:
		if statEvent.IsSuccess {
			t.Error("Stats expected to fail but suceeded")
		}
		t.Log(statEvent)
		return
	case <-timeout:
		t.Error("Timeout occurred")
		return
	}
}

func TestHealthRetreiveMetrics_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleHealthTimeout))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := types.Configuration{TraefikHealthEndpoint: server.URL + "/health"}
	channel := make(chan types.StatsEvent)

	go StartCheck(ctx, config, channel)

	timeout := time.After(time.Second * 5)

	select {
	case statEvent := <-channel:
		if statEvent.IsSuccess {
			t.Error("Stats expected to fail but suceeded")
		}
		if !strings.Contains(statEvent.ErrorDetails, "net/http: request canceled (Client.Timeout exceeded while awaiting headers)") {
			t.Error("Expected timout error")
		}
		t.Log(statEvent)
		return
	case <-timeout:
		t.Error("Test timeout occurred")
		return
	}
}

func handleHealthSuceed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadFile("testdata/healthresponse_normal.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Fatal(err)
	}
}

func handleHealthInvalid(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadFile("testdata/healthresponse_invalid.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Fatal(err)
	}
}

func handleHealthTimeout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		http.NotFound(w, r)
		return
	}

	time.Sleep(time.Second * 5)
}
