package herdcontrol

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGroup_Do_CoalescesRequests(t *testing.T) {
	const goroutines = 3000
	const key = 100
	var (
		g         = NewGroup()
		callCount int
	)
	fn := func() (any, error) {
		callCount++
		time.Sleep(300 * time.Microsecond)
		return "result", nil
	}

	var wg sync.WaitGroup
	results := sync.Map{}
	errs := sync.Map{}

	for k := 0; k < key; k++ {
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				res, err := g.Do("key_"+fmt.Sprintf("%d", k), fn)
				results.Store(idx, res)
				errs.Store(idx, err)
			}(i)
		}
		wg.Wait()
	}

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
