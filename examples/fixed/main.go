package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aileron-projects/go-backoff"
	"github.com/aileron-projects/go-backoff/examples"
)

func main() {
	csv, close := examples.NewCsvWriter("fixed.csv") // Output to csv.
	defer close()
	csv.Write([]string{"attempt", "backoff"})

	// Parameters
	interval := 100 * time.Millisecond

	// Instanciate a backoff provider.
	bo, err := backoff.NewFixed(interval)
	if err != nil {
		panic(err)
	}

	for i := range 101 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
		csv.Write([]string{strconv.Itoa(i), strconv.FormatInt(d.Nanoseconds(), 10)})
	}
}
