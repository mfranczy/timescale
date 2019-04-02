package main

import (
	"log"

	"github.com/mfranczy/timescale/pkg/args"
	"github.com/mfranczy/timescale/pkg/collector"
	"github.com/mfranczy/timescale/pkg/db"
	"github.com/mfranczy/timescale/pkg/job"
	"github.com/mfranczy/timescale/pkg/pool"
)

func main() {
	a, err := args.Parse()
	if err != nil {
		log.Fatal(err)
	}

	pg, err := db.Connect(a.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	stmt, err := db.GetCPUQuery(pg)
	if err != nil {
		log.Fatal(err)
	}

	j, err := job.New(a.CSVFile)
	if err != nil {
		log.Fatal(err)
	}

	p := pool.New(a.WorkersNum)
	rCh := collector.Process(p.OutCh)
	p.Run(j.Stream(), stmt)

	select {
	case r := <-rCh:
		r.PrintData(j.Num)
	}
}
