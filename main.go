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
	CONN_PORT = ":8000"
)

//Defines structs
type Node struct {
	Name      string   `json:"name`
	TimeSlice float32  `json:"timeslice`
	Cpu       float32  `json:"cpu"`
	Mem       float32  `json:"mem"`
	Process   *Process `json:"process,omitempty"`
}

type Process struct {
	Name      string  `json:"name`
	TimeSlice float32 `json:"timeslice"`
	CpuUsed   float32 `json:"cpuused"`
	MemUsed   float32 `json:"memused"`
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
			nodes[idx].Process = &process
		}
	}
	json.NewEncoder(w).Encode(nodes)
}

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/nodes/", nodesHandler).Methods("GET")
	r.HandleFunc("/analytics/nodes/average/", todoHandler).Methods("GET")
	r.HandleFunc("/analytics/processes/", todoHandler).Methods("GET")
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
	nodes = append(nodes, Node{Name: "kuryr", TimeSlice: 1.0, Cpu: 2.0, Mem: 3.0})
	// Running server!
	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(CONN_PORT, router))
}
