package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aileron-projects/go-backoff"
	"github.com/aileron-projects/go-backoff/examples"
)

func main() {
	csv, close := examples.NewCsvWriter("exponential.csv") // Output to csv.
	defer close()
	csv.Write([]string{"attempt", "backoff"})

	// Parameters
	offset := 0 * time.Millisecond
	limit := 600 * time.Second
	coeff := 1 * time.Millisecond
	base := 1.2
	jitter := backoff.NoJitter

	// Instanciate a backoff provider.
	bo, err := backoff.NewExponential(offset, limit, coeff, base, jitter)
	if err != nil {
		panic(err)
	}

	for i := range 101 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
		csv.Write([]string{strconv.Itoa(i), strconv.FormatInt(d.Nanoseconds(), 10)})
	}
}
