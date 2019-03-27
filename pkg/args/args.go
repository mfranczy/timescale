package args

import (
	"flag"
	"fmt"
	"strings"
)

// Args represents arguments provided by user
type Args struct {
	CSVFile    string
	WorkersNum int
}

// Parse arguments
func Parse() (Args, error) {
	a := Args{}
	flag.StringVar(&a.CSVFile, "csv-file", "", "Path to CSV file")
	flag.IntVar(&a.WorkersNum, "workers-num", 1, "Number of running workers between <1-100> range")
	flag.Parse()

	if a.CSVFile == "" {
		return a, fmt.Errorf("csv-file parameter is required")
	}
	if a.WorkersNum < 1 || a.WorkersNum > 100 {
		return a, fmt.Errorf("workers-num parameter has to be in <1-100> range")
	}
	if args := flag.Args(); len(args) > 0 {
		return a, fmt.Errorf("Invalid arguments: %s", strings.Join(args, " "))
	}

	return a, nil
}
