build:
	go get ./...
	go install cmd/cli/timescale-cli.go
docker-build:
	hack/dockerized.sh
run:
	hack/run.sh
