package snapshotter

import (
	"time"

	"github.com/hatchify/errors"
)

// NewConfig will return a new default Config with the given name and extension
func NewConfig(name, ext string) (cfg Config) {
	cfg.Name = name
	cfg.Extension = ext
	cfg.Interval = Minute
	cfg.Truncate = Hour
	cfg.TTL = Month
	return
}

// Config are the basic configuration settings for snapshotter
type Config struct {
	Name      string
	Extension string
	Interval  time.Duration
	Truncate  time.Duration
	TTL       time.Duration
}

// Validate will validate a Config
func (c *Config) Validate() (err error) {
	var errs errors.ErrorList
	// Ensure name is not empty
	if len(c.Name) == 0 {
		errs.Push(ErrInvalidName)
	}

	// Ensure extension is not empty
	if len(c.Extension) == 0 {
		errs.Push(ErrInvalidExtension)
	}

	// Check if truncate value is correct
	if !isValidTruncate(c.Truncate) {
		errs.Push(ErrInvalidTruncate)
	}

	// Ensure interval value is at least one second
	if c.Interval < Second {
		errs.Push(ErrInvalidInterval)
	}

	return errs.Err()
}
