# pong

Pong allows you to mock http responses through a `yml` configuration file 


## Usecase

Sometimes I need to test what happens with a mobile app when it gets a specific
response from an endpoint call. With pong I easily setup such responses.

I use pong along with [capture](https://github.com/ofabricio/capture), as a plugin.

Best pong's feature for me is `After`, because often I have to make an endpoint fail
*after* a successful call.


## Usage

This repo provides three ways of using *pong*:

1. As a mock server
1. As a plugin for [capture](https://github.com/ofabricio/capture) 
1. As a package

But before diving into them, let's see the yml configuration file.
It is a list of routes:

#### Routes configuration file

```yaml
-   status: 400
    match:
        path: /health
        method: GET
        after: 2
    headers:
        Content-Type: application/json
        Connection: close
    body: '{ "msg": "Oops!" }'
```

| Property | Description |
|---|---|
| status | **Required**. The http response status |
| match.path | This route matches if the request URL suffix matches with this path |
| match.method | This route matches if the request method matches |
| match.after | This route matches *after* a certain number of requests happened |
| headers | A map with any `key: value` pair |
| body | A string containing the response body |

### Mock server

This repository also provides a mock server for convenience.

How to build:

    git clone https://github.com/ofabricio/pong.git
    cd pong/cmd/server
    go build -o pong .

Run it and try two requests to `http://localhost:4000/health`.
The first returns 200 and the second returns 400, as told in `pong.yml`.
You can change it later as you like.

### Plugin

This repository also provides a plugin for [capture](https://github.com/ofabricio/capture).

How to build:

    git clone https://github.com/ofabricio/pong.git
    cd pong/cmd/plugin
    go build -buildmode=plugin -o pong.so main.go

Now copy the plugin to [capture](https://github.com/ofabricio/capture)'s directory.

### Package

    go get github.com/ofabricio/pong

```go
import (
    "net/http"
    "github.com/ofabricio/pong"
)

func main() {
    fallback := func(rw http.ResponseWriter, r *http.Request) {
        rw.WriteHeader(200)
    }
    http.HandleFunc("/", pong.Default(fallback))
    http.ListenAndServe(":4000", nil)
}
```

`pong.Default` will look for the `pong.yml` file in the current directory
by default and will fail if not found. So make sure it is there.
Use `pong.From` instead if you need to point somewhere else.
