# Timescale CLI

A command line tool that can be used to benchmark SELECT query performance across multiple workers/clients against a TimescaleDB instance.

## Launch project

### Dockerized environment

You can build this project as a dockerized environment by simply executing `make docker-build`.

It will create two containers, a **timescaledb container** with sample data provided from *artifacts/cpu_usage.csv* file and **timescale-cli container** which mounts *artifacts* directory into *timescale-cli:/data*.

> *project-path/artifacts* contains all provided files for this assignment

To run **timescale-cli** inside container you can execute `make run` with provided ENVs:

* CSV_FILE - path to CSV file
* WORKERS_NUM - number of running workers between <1-100> (default 1)

as example:

```sh
CSV_FILE=/data/query_params.csv WORKERS_NUM=10 make run
```

### Local environment

Run `make build` to build **timescale-cli** binary locally on your host.

> It requires $GOBIN env to be set and make sure that it has been added to $PATH

Now you should be able to execeute `timescale-cli --help`

```
Usage of timescale-cli:
  -csv-file - path to CSV file
  -db-config - path to database config file (default "/etc/timescale-cli/db.yaml")
  -workers-num - number of running workers between <1-100> (default 1)
```

you can combine this solution with the dockerized environment, first run `make docker-build` which will create a timescaledb container with sample data and then execute:

```
timescale-cli --csv-file artifacts/query_params.csv --workers-num 10 --db-config artifacts/db.yaml
```

## Architecture overview

![Architecture](arch.png)