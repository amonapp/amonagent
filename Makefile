build:
	go build -o amonagent -ldflags \
		"-X main.Version=$(VERSION)" \
		./cmd/amonagent.go
