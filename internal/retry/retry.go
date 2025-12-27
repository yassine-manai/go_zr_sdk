package retry

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/yassine-manai/go_zr_sdk/config"
	"github.com/yassine-manai/go_zr_sdk/errors"
	"github.com/yassine-manai/go_zr_sdk/logger"
)

// Retryer handles retry logic with exponential backoff
type Retryer struct {
	config config.RetryConfig
	logger logger.Logger
}

// New creates a new retryer
func New(cfg config.RetryConfig, log logger.Logger) *Retryer {
	return &Retryer{
		config: cfg,
		logger: log,
	}
}

// Do executes a function with retry logic
func (r *Retryer) Do(ctx context.Context, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		// Execute operation
		err := operation()

		// Success - return immediately
		if err == nil {
			if attempt > 0 {
				r.logger.Info("operation succeeded after retry",
					logger.Int("attempts", attempt+1),
				)
			}
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !errors.IsRetryable(err) {
			r.logger.Debug("error is not retryable, stopping",
				logger.Error(err),
			)
			return err
		}

		// Last attempt - don't sleep
		if attempt == r.config.MaxRetries {
			r.logger.Warn("max retries reached",
				logger.Int("attempts", attempt+1),
				logger.Error(err),
			)
			break
		}

		// Calculate backoff
		backoff := r.calculateBackoff(attempt, err)

		r.logger.Warn("operation failed, retrying",
			logger.Int("attempt", attempt+1),
			logger.Int("max_retries", r.config.MaxRetries),
			logger.Duration("backoff", backoff),
			logger.Error(err),
		)

		// Wait with context support
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	return lastErr
}

// calculateBackoff calculates backoff duration with exponential backoff + jitter
func (r *Retryer) calculateBackoff(attempt int, err error) time.Duration {
	// Check if error has specific retry-after
	retryAfter := errors.GetRetryAfter(err)
	if retryAfter > 0 {
		return time.Duration(retryAfter) * time.Second
	}

	// Exponential backoff: initialBackoff * (multiplier ^ attempt)
	backoff := float64(r.config.InitialBackoff) * math.Pow(r.config.Multiplier, float64(attempt))

	// Cap at max backoff
	if backoff > float64(r.config.MaxBackoff) {
		backoff = float64(r.config.MaxBackoff)
	}

	// Add jitter (±25%)
	jitter := backoff * 0.25
	backoff = backoff - jitter + (rand.Float64() * jitter * 2)

	return time.Duration(backoff)
}

// DoWithResult executes a function with retry and returns result
func DoWithResult[T any](ctx context.Context, r *Retryer, operation func() (T, error)) (T, error) {
	var result T
	var lastErr error

	err := r.Do(ctx, func() error {
		var err error
		result, err = operation()
		lastErr = err
		return err
	})

	if err != nil {
		return result, lastErr
	}

	return result, nil
}
