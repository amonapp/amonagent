VERSION := $(shell sh -c 'git describe --always --tags')

build:
	go build -o amonagent -ldflags \
		"-X main.Version=$(VERSION)" \
		./cmd/amonagent.go
