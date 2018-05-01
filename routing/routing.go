package routing

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lawrencegripper/traefik-appinsights-watchdog/types"
)

// StartCheck begins checking that traefik is routing information successfully by settings up a
// dummy server and pushing requests through traefik back to itself.
func StartCheck(config types.Configuration, healthChannel chan<- types.StatsEvent) {
	context := RequestContext{
		Port:              config.WatchdogTestServerPort,
		FabricURI:         config.TraefikBackendName,
		TraefikServiceURL: config.WatchdogTraefikURL,
		StartTime:         time.Now(),
		InstanceID:        config.InstanceID,
	}
	intervalDuration := time.Second * time.Duration(config.PollIntervalSec)
	go context.runServer()
	for {
		context.StartTime = time.Now()
		context.Nonce = uuid.New().String()
		healthChannel <- context.makeRequest()
		time.Sleep(intervalDuration)
	}
}

func (context *RequestContext) receiveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookie(generateCookieName(context.FabricURI)))
	w.Header().Set("x-response-from", context.InstanceID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(context.Nonce))
}

func (context *RequestContext) runServer() {
	http.HandleFunc("/", context.receiveHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%v", context.Port), nil)
	if err != nil {
		panic(err)
	}
}

func (context *RequestContext) makeRequest() types.StatsEvent {
	event := types.StatsEvent{
		Source:          "RoutingCheck",
		SourceTime:      time.Now(),
		Data:            make(map[string]interface{}),
		IsSuccess:       false,
		RequestDuration: time.Since(context.StartTime),
	}

	client := &http.Client{
		Timeout: time.Second * 3,
	}

	req, err := http.NewRequest("GET", context.TraefikServiceURL, nil)
	if err != nil {
		event.ErrorDetails = err.Error()
		return event
	}

	//Set a cookie to ensure sticky sessions route the request to this service.
	req.AddCookie(&http.Cookie{
		Expires: time.Now().Add(time.Hour),
		Domain:  "localhost",
		Name:    generateCookieName(context.FabricURI),
		Value:   fmt.Sprintf("http://%v:%v/", getOutboundIP(), context.Port),
		Path:    "/",
	})

	result, err := client.Do(req)
	if err != nil {
		event.ErrorDetails = err.Error()
		return event
	}

	event.Data["statusCode"] = result.StatusCode

	if result.StatusCode != http.StatusOK {
		event.ErrorDetails = "Returned failure code"
		return event
	}

	responseFrom := result.Header.Get("x-response-from")
	if responseFrom != context.InstanceID {
		event.ErrorDetails = fmt.Sprintf("Response from wrong instance expected: %s got response from: %s", context.InstanceID, responseFrom)
		return event
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		event.ErrorDetails = "Unable to read request body"
		return event
	}

	if string(body) != context.Nonce {
		event.ErrorDetails = fmt.Sprintf("Returned value doesn't match got: %s expected: %s response was from: %s", string(body), context.Nonce, responseFrom)
		return event
	}

	event.IsSuccess = true
	return event
}

const cookieNameLength = 6

// GenerateName Generate a hashed name
func generateCookieName(backendName string) string {
	data := []byte("_TRAEFIK_BACKEND_" + backendName)

	hash := sha1.New()
	_, err := hash.Write(data)
	if err != nil {
		// Impossible case
		panic(err)
	}

	return fmt.Sprintf("_%x", hash.Sum(nil))[:cookieNameLength]
}

// Get preferred outbound ip of this machine
// no connection is made so no traffic leaves the box.
// conn object only used to find ip
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
