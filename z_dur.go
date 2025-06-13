package zerocfg

import (
	"encoding/json"
	"fmt"
	"time"
)

type durationValue time.Duration

func newDuration(val time.Duration, p *time.Duration) Value {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(val string) error {
	duration, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	*d = durationValue(duration)
	return nil
}

func (d *durationValue) Type() string {
	return "duration"
}

// Dur registers a time.Duration configuration option and returns a pointer to its value.
//
// Usage:
//
//	timeout := zerocfg.Dur("timeout", 5*time.Second, "timeout for operation")
func Dur(name string, value time.Duration, usage string, opts ...OptNode) *time.Duration {
	return Any(name, value, usage, newDuration, opts...)
}

type durationSliceValue []time.Duration

func newDurationSlice(val []time.Duration, p *[]time.Duration) Value {
	*p = val
	return (*durationSliceValue)(p)
}

func (s *durationSliceValue) Set(val string) error {
	var strRepr []string
	if err := json.Unmarshal([]byte(val), &strRepr); err != nil {
		return err
	}

	ds := make([]time.Duration, 0, len(strRepr))
	for _, str := range strRepr {
		d, err := time.ParseDuration(str)
		if err != nil {
			return fmt.Errorf("duration %q is not a valid duration", str)
		}

		ds = append(ds, d)
	}

	*s = ds
	return nil
}

func (s *durationSliceValue) Type() string {
	return "durations"
}

// Durs registers a slice of time.Duration configuration options and returns a pointer to its value.
//
// Usage:
//
//	intervals := zerocfg.Durs("intervals", []time.Duration{time.Second, 2 * time.Second}, "interval durations")
func Durs(name string, defValue []time.Duration, desc string, opts ...OptNode) *[]time.Duration {
	return Any(name, defValue, desc, newDurationSlice, opts...)
}
