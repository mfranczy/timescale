FROM golang:1.11

ENV GOBIN=$GOPATH/bin
ARG db_conf=/etc/timescale-cli/
ARG project=$GOPATH/src/timescale

RUN mkdir ${project} ${db_conf}
ADD . ${project}
WORKDIR ${project}

RUN set -ex \
    && mv artifacts/db.yaml ${db_conf} \
    && go get ./... \
    && go install cmd/cli/timescale-cli.go

ENTRYPOINT [ "timescale-cli" ]