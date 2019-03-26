package args

import (
	"flag"
)

// TODO:
// -- worker num validation (check if number is correct and if it didn't reach maximum)
// -- validate non existing option

// Args represents arguments provided by user
type Args struct {
	CSVFile    string
	WorkersNum int
}

// Parse arguments
func Parse() Args {
	a := Args{}
	flag.StringVar(&a.CSVFile, "csv-file", "", "Path to CSV file")
	flag.IntVar(&a.WorkersNum, "workers-num", 1, "Number of running workers")
	flag.Parse()
	return a
}
