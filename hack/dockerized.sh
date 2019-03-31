#!/bin/bash

docker rm -f timescaledb
set -ex \
&& docker run -d --name timescaledb -e POSTGRES_PASSWORD=password -p 5432:5432 timescale/timescaledb:1.2.2-pg10 \
&& sleep 5 \
&& docker run -i --net=host --env PGPASSWORD=password --rm timescale/timescaledb:1.2.2-pg10 psql -h localhost -U postgres < artifacts/cpu_usage.sql \
&& docker run -i --net=host --env PGPASSWORD=password -v $PWD/artifacts:/data:ro --rm timescale/timescaledb:1.2.2-pg10 psql -d homework -h localhost -U postgres -c "\COPY cpu_usage FROM '/data/cpu_usage.csv' CSV HEADER" \
&& docker build -t timescale-cli:latest $PWD