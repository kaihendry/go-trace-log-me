VERSION := 0.2

run:
	PORT=3000 ENDPOINT=http://localhost:3500 go run main.go

.PHONY: build-AMainFunction
build-AMainFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/main
