package loadbalancer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
)

const (
	defaultMaxRetries     = 5
	defaultInitialBackoff = 1 * time.Second
	defaultMaxBackoff     = 32 * time.Second
	defaultRatePerSec     = 10
	defaultRateBurst      = 5
)

// do apply a central rate limit before performing the given operation.
func (a *Adapter) do(ctx context.Context, op func() error) error {
	if a.limiter != nil {
		if err := a.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	return op()
}

// retryOnRateLimit retries the provided operation when OCI responds with HTTP 429 rate limited.
// It applies exponential backoff between retries and preserves the original behavior and error messages.
func retryOnRateLimit(ctx context.Context, maxRetries int, initialBackoff, maxBackoff time.Duration, op func() error) error {
	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := op()
		if err == nil {
			return nil
		}

		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			if attempt == maxRetries-1 {
				return fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}
			// add jitter (up to 25% of backoff)
			var sleepDur = backoff
			jitter := time.Duration(time.Now().UnixNano() % int64(backoff/4))
			sleepDur = backoff + jitter
			t := time.NewTimer(sleepDur)
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return ctx.Err()
			}
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		return err
	}
	return nil
}
