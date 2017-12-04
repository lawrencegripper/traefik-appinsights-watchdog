Traefik Application Insights Watchdog
=============

> *Please note*: This project is in development. The current readme is a draft. 

## What is it?

[Application Insights](https://azure.microsoft.com/en-us/services/application-insights/) is a managed logging solution in Azure. This watchdog service polls an instance of [Traefik](https://docs.traefik.io/) - reporting it's health metrics and checking it's routing.

## How do I use it?

The simplest way to use it is to deploy Traefik on Service Fabric, using the [deployment guide here](https://aka.ms/traefikonsf).

If you would like to test or deploy separately here is a guide to launching the watchdog.

> *WARNING*: No error is shown if the Application Insights key provided is incorrect. If you don't see events surfaced check the key is correct.

``` text
    11:13 $ ./traefik-appinsights-watchdog start -h
Starts the watchdog, checking both the /health endpoint and request routing

Usage: start [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: start [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --appinsightskey         The application insights instrumentation key

    --debug                  Set to true for additional output in the console                              (default "false")
    --instanceid             The name to report for the instance                                           (default "nodename")
    --pollintervalsec        The time waited between requests to the health endpoint                       (default "120")
    --traefikbackendname     This is the Traefik backend name of the watchdog test server. In SF this will be the fabricURI (default "fabric:/TraefikType/Watchdog")
    --traefikhealthendpoint  The traeifk health endpoint                                                   (default "http://localhost:8080/health")
    --watchdogtestserverport Port which the simulated backend runs on                                      (default "40001")
    --watchdogtraefikurl     The url traefik will use to route requests to the watchdog                    (default "http://localhost:80/TraefikType/Watchdog")
-h, --help                   Print Help (this message) and exit d
```

Events will then be added to your Application Insights instance as `CustomEvents`. You can query these using the [Analytics portal](https://docs.microsoft.com/en-us/azure/application-insights/app-insights-analytics). Metrics from traefik will appear under `customdimensions` on the events. 

To see what will be logged set `--debug=true` and the watchdog will output the events as `json` into the console. For example you would see:

Run command: `./traefik-appinsights-watchdog --appinsightskey=YourKeyHere`

``` json
 {
 "Debug": false,
 "InstanceID": "Lawrences-Machine",
 "AppInsightsKey": "YourKeyHere",
 "TraefikBackendName": "fabric:/TraefikType/Watchdog",
 "WatchdogTestServerPort": 40001,
 "WatchdogTraefikURL": "http://localhost:80/TraefikType/Watchdog",
 "TraefikHealthEndpoint": "http://localhost:8080/health",
 "PollIntervalSec": 120
}
{
 "Source": "HealthCheck",
 "SourceTime": "2017-12-04T11:25:39.429175Z",
 "RequestDuration": 5539198,
 "IsSuccess": true,
 "ErrorDetails": "",
 "Data": {
  "average_response_time": "0s",
  "average_response_time_sec": 0,
  "count": 0,
  "pid": 10311,
  "recent_errors": [],
  "status_code_count": {},
  "time": "2017-12-04 11:25:39.433441 +0000 GMT m=+22.454186706",
  "total_count": 0,
  "total_response_time": "0s",
  "total_response_time_sec": 0,
  "total_status_code_count": {},
  "unixtime": 1512386739,
  "uptime": "22.291921753s",
  "uptime_sec": 22.291921753
 }
}
{
 "Source": "RoutingCheck",
 "SourceTime": "2017-12-04T11:25:39.430495Z",
 "RequestDuration": 1306232,
 "IsSuccess": false,
 "ErrorDetails": "Get http://localhost:80/TraefikType/Watchdog: dial tcp [::1]:80: getsockopt: connection refused",
 "Data": {}
}
```

## Why was it built?

It was primarily designed to run inside Service Fabric to support the use of Traefik on Service Fabric. Having the watchdog reporting stats from an independent service allows it to log when Traefik is not responding or has crashed. This is preferable to a gap in reporting, which would signal a failure if the stats where reported by the Traefik service in process.

However, it can be run independently inside other orchestrator's such a Kubernetes. *Please Note* Use outside of Service Fabric will require some additional testing and adjustment of the default values to ensure expected behavior.

## What is it's status?

> Currently under development.

This project is a simple watchdog service provided 'as is' and we currently have no plans to expands it's feature set.

We would welcome contributions. If you identify issues please log them and, if possible, develop a fix in as a PR.
