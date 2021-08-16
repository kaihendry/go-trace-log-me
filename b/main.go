package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"

	"github.com/apex/gateway/v2"
)

var Version string

func main() {
	log.SetHandler(jsonhandler.Default)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// the trace ID header probably has a different name depending on the system
		// i'm assuming it's AWS and I'm not not sure if I should be checking for others
		traceID := r.Header.Get("x-amzn-trace-id")

		// log using ctx so we have the traceID to search with
		// TODO: shouldn't http req path et al be logged here too?
		ctx := log.WithFields(log.Fields{"traceID": traceID})

		if traceID == "" {
			ctx.Error("missing trace ID")
			fmt.Fprintf(w, "b has no trace ID passed to via the x-amzn-trace-id header")
		} else {
			fmt.Fprintf(w, "b "+traceID)
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
