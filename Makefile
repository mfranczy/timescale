build:
	go get -t ./...
	go install cmd/cli/timescale-cli.go
docker-build:
	hack/dockerized.sh
run:
	hack/run.sh
test:
	go test ./...
