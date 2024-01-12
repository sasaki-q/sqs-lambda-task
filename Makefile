build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./infra/bin/handler ./lambda/main.go
	zip ./infra/bin/handler.zip ./infra/bin/handler
	rm -rf ./infra/bin/handler