package tracing

import (
	"os"
	"strings"
)

// HeaderKey is the key we will forward to chain the request through services together
const HeaderKey string = "x-foobar"

func EnvMap() map[string]string {
	envmap := make(map[string]string)
	for _, e := range os.Environ() {
		ep := strings.SplitN(e, "=", 2)
		// SEC for SECURITY
		if strings.Contains(ep[0], "SEC") || strings.Contains(ep[0], "TOKEN") {
			continue
		}
		envmap[ep[0]] = ep[1]
	}
	return envmap
}
