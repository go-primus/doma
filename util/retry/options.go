package retry

import "time"

type config struct {
	attempts     uint
	sleep        time.Duration
	maxSleepTime time.Duration
	isRetryErr   func(err error) bool
}

func newDefaultConfig() *config {
	return &config{
		attempts:     uint(10),
		sleep:        200 * time.Millisecond,
		maxSleepTime: 3 * time.Second,
	}
}

// Option is used to config the retry function.
type Option func(*config)

func AttemptAlways() Option {
	return func(c *config) {
		c.attempts = 0
	}
}

// Attempts is used to config the max retry times.
func Attempts(attempts uint) Option {
	return func(c *config) {
		c.attempts = attempts
	}
}

// Sleep is used to config the initial interval time of each execution.
func Sleep(sleep time.Duration) Option {
	return func(c *config) {
		c.sleep = sleep
		// ensure max retry interval is always larger than retry interval
		if c.sleep*2 > c.maxSleepTime {
			c.maxSleepTime = 2 * c.sleep
		}
	}
}

// MaxSleepTime is used to config the max interval time of each execution.
func MaxSleepTime(maxSleepTime time.Duration) Option {
	return func(c *config) {
		// ensure max retry interval is always larger than retry interval
		if c.sleep*2 > maxSleepTime {
			c.maxSleepTime = 2 * c.sleep
		} else {
			c.maxSleepTime = maxSleepTime
		}
	}
}

func RetryErr(isRetryErr func(err error) bool) Option {
	return func(c *config) {
		c.isRetryErr = isRetryErr
	}
}
