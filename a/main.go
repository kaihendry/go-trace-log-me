package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/google/uuid"

	"github.com/apex/gateway/v2"
)

// Version is set during the build Makefile
var Version string

func main() {
	log.SetHandler(jsonhandler.Default)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		traceID, ok := os.LookupEnv("x-amzn-trace-id")
		if !ok {
			// generate UUID
			traceID = uuid.New().String()
		}

		ctx := log.WithFields(log.Fields{"traceID": traceID, "service": "a"})

		endpoint, ok := os.LookupEnv("ENDPOINT")
		if !ok {
			http.Error(w, fmt.Errorf("tracing endpoint is unset").Error(), http.StatusInternalServerError)
			return
		}

		ctx.WithField("endpoint", endpoint).Info("tracing endpoint")

		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Forward the header traceID so we can see it in the logs
		req.Header.Set("x-amzn-trace-id", traceID)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := ioutil.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		t, err := template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>{{ .Name }}</title>
</head>
<body>
<h1>Trace ID</h1>
<pre>{{ .TraceID }}</pre>
<h1>Response</h1>
<pre>
{{ .Response }}
</pre>
<dl>
{{range $key, $value := .Env -}}
<dt>{{ $key }}</dt><dd>{{ $value }}</dd>
{{- end}}
</dl>
<p><a href="https://github.com/kaihendry/x-amzn-trace-id">Source code</a></p>
</body>
</html>`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = t.Execute(w, struct {
			Name     string
			TraceID  string
			Response string
			Env      map[string]string
		}{
			Name:     os.Getenv("AWS_LAMBDA_FUNCTION_NAME") + Version,
			TraceID:  traceID,
			Response: string(response),
			Env:      envMap(),
		})

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
