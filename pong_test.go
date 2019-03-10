package pong

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var defaultTestFallback = func(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
}

func TestPong(t *testing.T) {

	// given

	input := `
-   status: 400
    match:
        path: /health
    headers:
        Content-Type: application/json
        Authorization: Bearer abc
    body: '{ "msg": "Oops!" }'
`
	req, _ := http.NewRequest("GET", "/health", nil)

	// when

	recs := pong(input, req)

	// then

	rec := recs[0]
	if rec.Code != 400 {
		t.Error("invalid code", rec.Code)
	}
	if rec.Body.String() != `{ "msg": "Oops!" }` {
		t.Error("invalid body", rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Error("invalid header", rec.Header())
	}
	if rec.Header().Get("Authorization") != "Bearer abc" {
		t.Error("invalid header", rec.Header())
	}
}

func TestPong_After(t *testing.T) {

	// given

	input := `
-   status: 400
    match:
        path: /health
        after: 2
    headers:
        Content-Type: application/json
    body: '{ "msg": "Oops!" }'
`
	req1, _ := http.NewRequest("GET", "/health", nil)
	req2, _ := http.NewRequest("GET", "/health", nil)

	// when

	recs := pong(input, req1, req2)

	// then

	if recs[0].Code != 200 {
		t.Error("invalid code", recs[0].Code)
	}
	if recs[1].Code != 400 {
		t.Error("invalid code", recs[1].Code)
	}
}

func TestPong_Method(t *testing.T) {

	// given

	input := `
-   status: 400
    match:
        path: /health
        method: POST
`
	req1, _ := http.NewRequest("GET", "/health", nil)
	req2, _ := http.NewRequest("POST", "/health", nil)

	// when

	recs := pong(input, req1, req2)

	// then

	if recs[0].Code != 200 {
		t.Error("invalid code", recs[0].Code)
	}
	if recs[1].Code != 400 {
		t.Error("invalid code", recs[1].Code)
	}
}

func pong(input string, reqs ...*http.Request) []*httptest.ResponseRecorder {
	var recs []*httptest.ResponseRecorder
	handler := From(strings.NewReader(input), defaultTestFallback)
	for _, r := range reqs {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, r)
		recs = append(recs, rec)
	}
	return recs
}
