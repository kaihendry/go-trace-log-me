VERSION := 0.2

run:
	PORT=3500 go run main.go

test:
	curl -H "x-amzn-trace-id: foobar" http://localhost:3500
