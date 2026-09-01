package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aileron-projects/go-backoff"
	"github.com/aileron-projects/go-backoff/examples"
)

func main() {
	csv, close := examples.NewCsvWriter("polynomial.csv") // Output to csv.
	defer close()
	csv.Write([]string{"attempt", "backoff"})

	polynomial := &backoff.PolynomialConfig{
		Offset:   0 * time.Millisecond,
		Limit:    600 * time.Second,
		Coeff:    10 * time.Millisecond,
		Exponent: 2.0,
		Jitter:   backoff.NoJitter,
	}

	bo, err := polynomial.New()
	if err != nil {
		panic(err)
	}

	for i := range 101 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
		csv.Write([]string{strconv.Itoa(i), strconv.FormatInt(d.Nanoseconds(), 10)})
	}
}
