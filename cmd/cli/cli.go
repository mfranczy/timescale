package main

import (
	"log"

	// "github.com/pkg/profile"

	"timescale/pkg/args"
	"timescale/pkg/collector"
	"timescale/pkg/db"
	"timescale/pkg/job"
	"timescale/pkg/pool"
)

func main() {
	a := args.Parse()
	pg, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	stmt, err := db.GetCpuQuery(pg)
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
