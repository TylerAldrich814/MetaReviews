package grpcutil

import (
	"errors"
	"log"
	"math"
	"time"
)

// DoRequestWithBackoff Makes 'maxRetries' attemps at making a successful network request.
func DoRequestWithBackoff[Resp any](
  maxRetries int,
  baseDelay  time.Duration,
  req func()(*Resp, error),
  shouldRetry func(error) bool,
)( *Resp, error ){
  var err error
  var resp *Resp

  for i := 0; i < maxRetries; i++ {
    if resp, err = req(); err == nil {
      return resp, err
    }
    if !shouldRetry(err) {
      return nil, err
    }

    if i < maxRetries-1 {
      sleepDuration := time.Duration(
        math.Pow(2, float64(i)),
      ) * baseDelay
      log.Printf(
        "attempt %d failed: %v. Retrying in %s...\n",
        i+1,err, sleepDuration,
      )
      time.Sleep(sleepDuration)
    }
  }

  return nil, errors.New("all retry attemps failed")
}
