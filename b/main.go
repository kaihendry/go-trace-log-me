package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	tracing "github.com/kaihendry/go-trace-log-me"
	"github.com/tj/go/http/response"

	"github.com/apex/gateway/v2"
)

var Version string

func main() {
	log.SetHandler(jsonhandler.Default)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// the trace ID header probably has a different name depending on the system
		// i'm assuming it's AWS and I'm not not sure if I should be checking for others
		traceID := r.Header.Get(tracing.HeaderKey)

		// log using ctx so we have the traceID to search with
		// TODO: shouldn't http req path et al be logged here too?
		ctx := log.WithFields(log.Fields{"traceID": traceID})

		if traceID == "" {
			ctx.Error("missing trace ID")
			http.Error(w, "missing trace ID", http.StatusBadRequest)
			return
		} else {
			response.OK(w, struct {
				Name    string
				TraceID string
				Env     map[string]string
				Header  http.Header
			}{
				Name:    os.Getenv("AWS_LAMBDA_FUNCTION_NAME") + Version,
				TraceID: traceID,
				Env:     envMap(),
				Header:  r.Header,
			})
			ctx.Info("responded with traceID")
		}
	})

	port := os.Getenv("_LAMBDA_SERVER_PORT")
	var err error
	if port == "" {
		err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	} else {
		err = gateway.ListenAndServe("", nil)
	}
	log.Fatalf("failed to start server: %v", err)
}

func envMap() map[string]string {
	envmap := make(map[string]string)
	for _, e := range os.Environ() {
		ep := strings.SplitN(e, "=", 2)
		if strings.Contains(ep[0], "SEC") || strings.Contains(ep[0], "TOKEN") {
			continue
		}
		envmap[ep[0]] = ep[1]
	}
	return envmap
}
