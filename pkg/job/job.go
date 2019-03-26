package job

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	recordLength = 3
	buffer       = 100
)

// Job represents CSV record
type Job struct {
	Hostname  string
	StartTime string
	EndTime   string
}

// Jobs represents all parsed CSV records
type Jobs struct {
	jobs chan Job
	Num  int
	r    *csv.Reader
}

// New jobs from csv file
func New(csvFile string) (*Jobs, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	return &Jobs{Num: 0, r: csv.NewReader(f)}, nil
}

// Stream jobs
func (j *Jobs) Stream() (ch chan Job) {
	ch = make(chan Job, buffer)
	skipHeader := true

	go func() {
		defer close(ch)
		for {
			j.Num++
			row, err := j.r.Read()
			if err != nil {
				if err == io.EOF {
					j.Num--
					break
				}
				log.Println("Error: ", err)
				continue
			}

			if skipHeader {
				j.Num--
				skipHeader = false
				continue
			}

			if err := validateDates(row); err != nil {
				log.Println("Error: ", err)
				continue
			}
			ch <- Job{row[0], row[1], row[2]}
		}
	}()
	return
}

func validateDates(row []string) error {
	var err error
	var startTime time.Time
	var endTime time.Time
	layout := "2006-01-02 15:04:05"

	if startTime, err = time.Parse(layout, row[1]); err != nil {
		return fmt.Errorf("invalid start_time date format, acceptable is 'yyyy:mm:dd hh:mm:ss'")
	}
	if endTime, err = time.Parse(layout, row[2]); err != nil {
		return fmt.Errorf("invalid end_time date format, acceptable is 'yyyy:mm:dd hh:mm:ss'")
	}
	if startTime.UnixNano() > endTime.UnixNano() {
		return fmt.Errorf("start_time %s is bigger than end_time %s", endTime, startTime)
	}
	return nil
}
