package routing

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/google/uuid"
	"github.com/lawrencegripper/sfTraefikWatchdog/types"
)

// StartCheck begins checking that traefik is routing information successfully by settings up a
// dummy server and pushing requests through traefik back to itself.
func StartCheck(config types.Configuration, healthChannel chan<- types.StatsEvent) {
	context := RequestContext{
		Port:              config.SimulatedBackendPort,
		FabricURI:         config.WatchdogFabricURI,
		TraefikServiceURL: config.WatchdogTraefikURL,
		StartTime:         time.Now(),
	}
	go context.runServer()
	for {
		context.StartTime = time.Now()
		nonceUUID, _ := uuid.NewUUID()
		context.Nonce = nonceUUID.String()
		healthChannel <- context.makeRequest()
		time.Sleep(3 * time.Second)
	}
}

func (context *RequestContext) receiveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookie(generateCookieName(context.FabricURI)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(context.Nonce))
}

func (context *RequestContext) runServer() {
	http.HandleFunc("/", context.receiveHandler)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%v", context.Port), nil)
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
		RequestDuration: time.Now().Sub(context.StartTime),
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
		Value:   fmt.Sprintf("http://%v:%v", getOutboundIP(), context.Port),
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

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		event.ErrorDetails = "Unable to read request body"
		return event
	}

	if string(body) != context.Nonce {
		event.ErrorDetails = fmt.Sprintf("Returned value doesn't match expected %s", string(body))
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
		log.Errorf("Fail to create cookie name: %v", err)
	}

	return fmt.Sprintf("_%x", hash.Sum(nil))[:cookieNameLength]
}

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "127.0.0.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
