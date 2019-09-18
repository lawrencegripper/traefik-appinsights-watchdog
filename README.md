Traefik Application Insights Watchdog
=============

[![Build Status](https://travis-ci.org/lawrencegripper/traefik-appinsights-watchdog.svg?branch=master)](https://travis-ci.org/lawrencegripper/traefik-appinsights-watchdog)

## Warning: This project is no longer maintained. It remains available as a reference for those currently using solution. Please fork and re-use or update the code as needed.

## What is it?

[Application Insights](https://azure.microsoft.com/en-us/services/application-insights/) is a managed logging solution in Azure. This watchdog service polls an instance of [Traefik](https://docs.traefik.io/) - reporting its health metrics and checking it's routing correctly.

## How do I use it?

The simplest way to use it is to deploy it along side Traefik on Service Fabric, using the [deployment guide here](https://aka.ms/traefikonsf).

If you would like to test or deploy separately here is a guide to launching the watchdog.

> *WARNING*: No error is shown if the Application Insights key provided is incorrect. If you don't see events surfaced check the key is correct.

``` text
$ ./traefik-appinsights-watchdog start -h
Starts the watchdog, checking both the /health endpoint and request routing

Usage: start [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: start [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
    --allowinvalidcert       Allow invalid certificates when performing routing checks on localhost        (default "false")
    --apiendpointpassword    Stores password required to call APIs including healthcheck
    --apiendpointusername    Stores username required to call APIs including healthcheck
    --appinsightskey         The application insights instrumentation key
    --debug                  Set to true for additional output in the console                              (default "false")
    --instanceid             The name to report for the instance                                           (default "nodename")
    --pollintervalsec        The time waited between requests to the health endpoint                       (default "120")
    --traefikbackendname     This is the Traefik backend name of the watchdog test server. In SF this      (default "fabric:/TraefikType/Watchdog")
                             will be the fabricURI
    --traefikhealthendpoint  The traeifk health endpoint http://localhost:port/health                      (default "http://localhost:8080/health")
    --watchdogtestserverport Port which the simulated backend runs on                                      (default "40001")
    --watchdogtraefikurl     The url traefik will use to route requests to the watchdog                    (default "http://localhost:80/TraefikType/Watchdog")
-h, --help                   Print Help (this message) and exit
```

Events will then be added to your Application Insights instance as `CustomEvents`. You can query these using the [Analytics portal](https://docs.microsoft.com/en-us/azure/application-insights/app-insights-analytics). Metrics from traefik will appear under `customdimensions` on the events. 

Here is an example query to graph `failures` vs `success` over the last `30mins` in the Analytics portal:

```
    customEvents 
    | where timestamp > ago(30m)  
    | summarize count() by tostring(customDimensions.isSuccess), bin(timestamp, 10s)
    | render timechart 
```

Here is a query to show full tabular data for the last `30mins`:

```
    customEvents 
    | where timestamp > ago(30m) 
    | order by timestamp desc 
```

To see what will be logged and for debugging output set `--debug=true` and the watchdog will output the events as `json` into the console. For example you would see:

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

**Note:** The events logged to stdout may differ to the actual events you see in Application Insights due to some post processing.

### Basic Auth in Traefik

If Traefik's API endpoint (including `/health`) is protected with Basic Authentication, watchdog's health-check calls will fail by default (will receive 401, naturally). Hence, 2 new optional parameters were added through which one could pass the required credentials (`apiendpointusername` and `apiendpointpassword`).

## Why was it built?

It was primarily designed to run inside Service Fabric to support the use of Traefik on Service Fabric. Having the watchdog reporting stats from an independent service allows it to log when Traefik is not responding or has crashed. This is preferable to a gap in reporting, which would signal a failure if the stats were reported by the Traefik service in process.

However, it can be run independently inside other orchestrator's such a Kubernetes. *Please Note* Use outside of Service Fabric will require some additional testing and adjustment of the default values to ensure expected behavior.

## What is it's status?

> Currently under development.

This project is a simple watchdog service provided 'as is' and we currently have no plans to expands it's feature set.

We would welcome contributions. If you identify issues please log them and, if possible, develop a fix in as a PR.

## Building

Run `build.sh` this has a dependency on docker.

The build with run a set of checks, execute tests and then output. Once completed you can use `docker run --rm traefik-appinsights-watchdog:latest --appinsightskey=XXXXXXXXX --debug` to test your build.
