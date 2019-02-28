# Golang rest-api example

## Installation & Run
```bash
# Download this project
go get github.com/danielmellado/rest-api (the repo is private as of now, so it might not work out of the blue)
```

```bash
go build
./rest-api
route: /v1, handler: <nil>
route: /v1/nodes/, handler: 0x65bc70
route: /v1/analytics/nodes/average/, handler: 0x65bd70
route: /v1/analytics/processes/, handler: 0x65c7e0
route: /v1/analytics/processes/{processname}/, handler: 0x65bd70
route: /v1/metrics/node/{nodename}/, handler: 0x65bdd0
route: /v1/metrics/nodes/{nodename}/process/{processname}/, handler: 0x65c260
Running server! (on localhost:8000 by default)
```

## Structure
```
├── helm                 // Helm Charts
├── Dockerfile           // Dockerfile used for deploying
├── .travis.yml          // Golang Travis CI configuration
├── rest-api_test.go     // Sample HTTP handler test for nodes
└── rest-api.go          // Main Rest API server, implemented with Gorilla
```

## API

#### /v1/nodes/
* `GET` : Get all nodes
```bash
$ curl -L http://localhost:8000/v1/nodes/
[{"Name":"node1","TimeSlice":60,"cpu":60,"mem":40}]
```

#### /metrics/nodes/{nodename}/process/{processname}/
* `POST` : Create a process on a node
```bash
$ curl -H "Content-Type: application/json" -X POST -d '{"TimeSlice":45.0, "cpu_used":80.0, "mem_used":80.0}' http://localhost:8000/v1/metrics/nodes/test/process/ps/
[{"Name":"test","TimeSlice":2,"cpu":3,"mem":11,"process":[{"Name":"ps","timeslice":45,"cpu_used":80,"mem_used":80}]}]
```

#### /metrics/node/{nodename}/
* `POST` : Create a node
```bash
$ curl -H "Content-Type: application/json" -X POST -d '{"TimeSlice":2.0, "Cpu":3.0, "Mem":11.0}' http://localhost:8000/v1/metrics/node/test/
[{"Name":"test","TimeSlice":2,"cpu":3,"mem":11}]
```

#### /analytics/processes/
* `GET` : Get a detailed output from the processes

### Run tests
```bash
$ go test -v
=== RUN   TestNodesHandler
--- PASS: TestNodesHandler (0.00s)
PASS
ok      github.com/danielmellado/rest-api       0.003s
```
