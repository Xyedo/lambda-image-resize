
build: 
	set GOOS=linux &\
	set GOARCH=amd64 &\
	go build -o app ./cmd/main.go

.PHONY: build