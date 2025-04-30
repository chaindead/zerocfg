package zerocfg

import (
	"encoding/json"
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

func (d *durationValue) String() string {
	return time.Duration(*d).String()
}

func Dur(name string, value time.Duration, usage string, opts ...OptNode) *time.Duration {
	return Any(name, value, usage, newDuration, opts...)
}

type durationSliceValue []time.Duration

func newDurationSlice(val []time.Duration, p *[]time.Duration) Value {
	*p = val
	return (*durationSliceValue)(p)
}

func (s *durationSliceValue) Set(val string) error {
	var durations []time.Duration
	if err := json.Unmarshal([]byte(val), &durations); err != nil {
		return err
	}

	*s = durations
	return nil
}

func (s *durationSliceValue) Type() string {
	return "durations"
}

func (s *durationSliceValue) String() string {
	ds := make([]string, 0, len(*s))

	for _, d := range *s {
		ds = append(ds, d.String())
	}

	data, _ := json.Marshal(ds)
	return string(data)
}

func Durs(name string, defValue []time.Duration, desc string, opts ...OptNode) *[]time.Duration {
	return Any(name, defValue, desc, newDurationSlice, opts...)
}
