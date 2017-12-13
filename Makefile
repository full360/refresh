VERSION=0.1.0
GO_LDFLAGS="-X main.version=$(VERSION)"

default: test

install:
	@go get -u ./...

test:
	@go test -v -race -cover .

bin:
	@mkdir -p bin/{darwin,linux}
	@rm -rf bin/*

release: release-darwin \
	release-linux

release-darwin: bin
	GOOS=darwin GOARCH=amd64 go build -ldflags=$(GO_LDFLAGS) -o bin/darwin/refresh ./cmd/refresh
	cd bin/darwin && tar -cvzf refresh.$(VERSION).darwin-amd64.tgz refresh

release-linux: bin
	GOOS=linux GOARCH=amd64 go build -ldflags=$(GO_LDFLAGS) -o bin/linux/refresh ./cmd/refresh
	cd bin/linux && tar -cvzf refresh.$(VERSION).linux-amd64.tgz refresh

.PHONY: default install test bin release
