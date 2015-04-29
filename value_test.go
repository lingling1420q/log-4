package log_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/peterbourgon/gokit/log"
)

func TestValueBinding(t *testing.T) {
	var output []interface{}

	logger := log.Logger(log.LoggerFunc(func(keyvals ...interface{}) error {
		output = keyvals
		return nil
	}))

	start := time.Date(2015, time.April, 25, 0, 0, 0, 0, time.UTC)
	now := start
	mocktime := func() time.Time {
		now = now.Add(time.Second)
		return now
	}

	logger = log.With(logger, "ts", log.Timestamp(mocktime), "caller", log.DefaultCaller)

	logger.Log("foo", "bar")
	timestamp, ok := output[1].(time.Time)
	if !ok {
		t.Fatalf("want time.Time, have %T", output[1])
	}
	if want, have := start.Add(time.Second), timestamp; want != have {
		t.Errorf("output[1]: want %v, have %v", want, have)
	}
	if want, have := "value_test.go:28", fmt.Sprint(output[3]); want != have {
		t.Errorf("output[3]: want %s, have %s", want, have)
	}

	// A second attempt to confirm the bindings are truly dynamic.
	logger.Log("foo", "bar")
	timestamp, ok = output[1].(time.Time)
	if !ok {
		t.Fatalf("want time.Time, have %T", output[1])
	}
	if want, have := start.Add(2*time.Second), timestamp; want != have {
		t.Errorf("output[1]: want %v, have %v", want, have)
	}
	if want, have := "value_test.go:41", fmt.Sprint(output[3]); want != have {
		t.Errorf("output[3]: want %s, have %s", want, have)
	}
}

func BenchmarkValueBindingTimestamp(b *testing.B) {
	logger := discard
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Log("k", "v")
	}
}

func BenchmarkValueBindingCaller(b *testing.B) {
	logger := discard
	logger = log.With(logger, "caller", log.DefaultCaller)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Log("k", "v")
	}
}
