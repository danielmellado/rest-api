// rest-api_test.go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNodesHandler(t *testing.T) {
	// This is lacking the mock response to the API
	// but it would've checked responses with actual mocked one
	req, err := http.NewRequest("GET", "/v1/nodes", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	nodesHandler(res, req)
}
