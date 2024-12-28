.PHONY: test fmt

test:
	go test -v ./...

gen:
	go generate

fmt:
	go mod tidy
	gofmt -s -w .
	goarrange run -r .
	golangci-lint run ./...
