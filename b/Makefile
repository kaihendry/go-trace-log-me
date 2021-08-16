VERSION := 0.2

run:
	PORT=3500 go run main.go

.PHONY: build-BMainFunction
build-BMainFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/main

test:
	curl -H "x-amzn-trace-id: foobar" http://localhost:3500
