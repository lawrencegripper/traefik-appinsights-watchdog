package main

import (
	"fmt"

	"github.com/lawrencegripper/sfTraefikWatchdog/health"
	"github.com/lawrencegripper/sfTraefikWatchdog/routing"
	"github.com/lawrencegripper/sfTraefikWatchdog/types"
)

func main() {
	fmt.Println("Hello")

	watchdogPort := 8988

	healthChan := make(chan types.StatsEvent)

	context := routing.RequestContext{Port: watchdogPort}
	context.FabricURI = "fabric:/TraefikType/Watchdog"
	context.TraefikServiceURL = fmt.Sprintf("http://localhost:%v", context.Port)
	context.BackendURL = fmt.Sprintf("http://localhost:%v", context.Port)

	//Todo:
	// Add configuration to health
	// Add configuration to routing
	config := types.Configuration{}

	go routing.StartCheck(context, healthChan)
	go stats.StartCheck()

	for {
		event := <-healthChan
		fmt.Println(event.IsSuccess)
	}
}
