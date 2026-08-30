package main

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand/v2"
	"time"

	"github.com/aileron-projects/go-backoff"
)

func main() {
	// Prepare backoff strategy.
	bo, err := backoff.NewFixed(time.Second)
	if err != nil {
		panic(err)
	}

	retryer := backoff.NewRetryer(bo)
	retryer.MaxRetry = 20                     // Change me!
	retryer.MaxElapsedTime = 60 * time.Second // Change me!
	ctxTimeout := 10 * time.Second            // Change me!
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	// Run retryable function.
	counter := 0
	err = retryer.RunContext(ctx, func(notify *backoff.Notify) error {
		counter++
		log.Printf("Hello! It's %d-th call.\n", counter)
		random := rand.Float64()
		switch {
		case random > 0.9: // 10% suceess
			notify.Success()
			log.Println("NOTIFY success !!")
		case random > 0.8: // 10% stop with error
			notify.Stop()
			log.Println("NOTIFY stop !!")
		default: // 80% continue with error
			return io.ErrUnexpectedEOF
		}
		return nil
	})

	// Handle errors if any.
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("ERROR Canceled:", err)
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("ERROR DeadlineExceeded:", err)
		case errors.Is(err, backoff.ErrMaxRetry):
			log.Println("ERROR ErrMaxRetry:", err)
		case errors.Is(err, backoff.ErrTimeout):
			log.Println("ERROR ErrTimeout:", err)
		default:
			log.Println("ERROR others:", err)
		}
	} else {
		log.Println("Finished with success !!")
	}
}
