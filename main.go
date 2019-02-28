/*
 * Copyright 2019 Daniel Mellado
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	// ConnPort defines the port to use for the api server
	ConnPort = ":8000"
)

// Node struct representing the node object
type Node struct {
	Name      string     `json:"name`
	TimeSlice float32    `json:"timeslice`
	Cpu       float32    `json:"cpu"`
	Mem       float32    `json:"mem"`
	Process   []*Process `json:"process,omitempty"`
}

// Process struct representing the processes on the node
type Process struct {
	Name      string  `json:"name`
	TimeSlice float32 `json:"timeslice"`
	CpuUsed   float32 `json:"cpu_used"`
	MemUsed   float32 `json:"mem_used"`
}

var nodes []Node

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(nodes)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("TODO")
}

func createNodeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	node := Node{}
	err := json.NewDecoder(r.Body).Decode(&node)
	node.Name = params["nodename"]
	if err != nil {
		log.Print("error occurred while decoding node :: ", err)
		return
	}
	nodes = append(nodes, Node{Name: node.Name, TimeSlice: node.TimeSlice, Cpu: node.Cpu, Mem: node.Mem, Process: node.Process})
	json.NewEncoder(w).Encode(nodes)
}

func createProcessHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	process := Process{}
	err := json.NewDecoder(r.Body).Decode(&process)
	process.Name = params["processname"]
	if err != nil {
		log.Print("error occurred while decoding process :: ", err)
		return
	}
	for idx, node := range nodes {
		if node.Name == params["nodename"] {
			nodes[idx].Process = append(nodes[idx].Process, &process)
		}
	}
	json.NewEncoder(w).Encode(nodes)
}

func getProcessesHandler(w http.ResponseWriter, r *http.Request) {
	for _, node := range nodes {
		json.NewEncoder(w).Encode(node.Process)
	}
}

// AddRoutes creates the mux routes for the rest api
func AddRoutes(r *mux.Router) {
	r.HandleFunc("/nodes/", nodesHandler).Methods("GET")
	r.HandleFunc("/analytics/nodes/average/", todoHandler).Methods("GET")
	r.HandleFunc("/analytics/processes/", getProcessesHandler).Methods("GET")
	r.HandleFunc("/analytics/processes/{processname}/", todoHandler).Methods("GET")
	r.HandleFunc("/metrics/node/{nodename}/", createNodeHandler).Methods("POST")
	r.HandleFunc("/metrics/nodes/{nodename}/process/{processname}/", createProcessHandler).Methods("POST")
}

func main() {
	router := mux.NewRouter()
	// Adds a route prefix for v1 requests.
	AddRoutes(router.PathPrefix("/v1").Subrouter())
	// Print available routes to the terminal
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tmpl, _ := route.GetPathTemplate()
		fmt.Printf("route: %s, handler: %v\n", tmpl, route.GetHandler())
		return nil
	})
	// Initial dummy data
	nodes = append(nodes, Node{Name: "node1", TimeSlice: 60.0, Cpu: 60.0, Mem: 40.0})
	nodes = append(nodes, Node{Name: "node2", TimeSlice: 60.0, Cpu: 40.0, Mem: 60.0})
	// Running server!
	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(ConnPort, router))
}
