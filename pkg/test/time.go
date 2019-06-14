package test

import (
	"sync/atomic"
	"testing"
	"time"
)

func WithTimeout(t *testing.T, timeout time.Duration, test func(t *testing.T)) {
	var done int32 = 0
	go func() {
		test(t)
		atomic.AddInt32(&done, 1)
	}()

	time.Sleep(timeout)
	if atomic.LoadInt32(&done) == 1 {
		return
	}

	t.Errorf("%s: timeout waiting for test to complete", t.Name())
}
