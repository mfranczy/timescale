#!/bin/bash

docker run -it -v $PWD/artifacts:/data:ro --net=host --rm timescale-cli --csv-file $CSV_FILE --workers-num $WORKERS_NUM