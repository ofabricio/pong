package pong

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func Default(next http.HandlerFunc) http.HandlerFunc {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(filepath.Dir(ex) + "/pong.yml")
	if err != nil {
		panic(err)
	}
	return From(bytes.NewReader(data), next)
}

func From(r io.Reader, next http.HandlerFunc) http.HandlerFunc {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	var routes []*Route
	if err := yaml.Unmarshal(data, &routes); err != nil {
		panic(err)
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		for _, route := range routes {
			if route.match(req) {
				route.ServeHTTP(rw, req)
				return
			}
		}
		next.ServeHTTP(rw, req)
	}
}

type Route struct {
	Status int
	Match  struct {
		Path   string
		Method string
		After  int
	}
	Headers map[string]string
	Body    string
}

func (r *Route) match(req *http.Request) bool {
	if r.matchMethod(req) && r.matchPath(req) && r.matchAfter(req) {
		return true
	}
	return false
}

func (r *Route) matchMethod(req *http.Request) bool {
	return r.Match.Method == "" || r.Match.Method == req.Method
}

func (r *Route) matchPath(req *http.Request) bool {
	return strings.HasSuffix(req.URL.String(), r.Match.Path)
}

func (r *Route) matchAfter(req *http.Request) bool {
	// TODO: we probably need a lock here
	if r.Match.After > 1 {
		r.Match.After--
		return false
	}
	return true
}

func (r *Route) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	for k, v := range r.Headers {
		rw.Header().Set(k, v)
	}
	rw.WriteHeader(r.Status)
	rw.Write([]byte(r.Body))
}
