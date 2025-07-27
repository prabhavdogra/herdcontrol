package heardcontrol

import (
	"sync"
	"testing"
	"time"
)

func TestGroup_Do_CoalescesRequests(t *testing.T) {
	const goroutines = 10000
	var (
		g         = NewGroup()
		callCount int
	)
	fn := func() (any, error) {
		callCount++
		time.Sleep(10 * time.Millisecond)
		return "result", nil
	}

	var wg sync.WaitGroup
	results := sync.Map{}
	errs := sync.Map{}

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			res, err := g.Do("key", fn)
			results.Store(idx, res)
			errs.Store(idx, err)
		}(i)
	}
	wg.Wait()

	if callCount != 1 {
		t.Errorf("fn was called %d times, want 1", callCount)
	}
	for i := 0; i < goroutines; i++ {
		res, ok := results.Load(i)
		if !ok || res != "result" {
			t.Errorf("got result %v, want 'result'", res)
		}
		err, ok := errs.Load(i)
		if !ok || err != nil {
			t.Errorf("got error %v, want nil", err)
		}
	}
}
