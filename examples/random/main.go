package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aileron-projects/go-backoff"
	"github.com/aileron-projects/go-backoff/examples"
)

func main() {
	csv, close := examples.NewCsvWriter("random.csv") // Output to csv.
	defer close()
	csv.Write([]string{"attempt", "backoff"})

	// Parameters
	offset := 0 * time.Millisecond
	limit := 100 * time.Millisecond

	// Instanciate a backoff provider.
	bo, err := backoff.NewRandom(offset, limit)
	if err != nil {
		panic(err)
	}

	for i := range 101 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
		csv.Write([]string{strconv.Itoa(i), strconv.FormatInt(d.Nanoseconds(), 10)})
	}
}
