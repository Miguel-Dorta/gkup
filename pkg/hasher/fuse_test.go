package hasher

import (
	"runtime"
	"sync"
	"testing"
)

func TestFuse(t *testing.T) {
	var f fuse
	var wg sync.WaitGroup
	var triggerFuseTimesN = func(f *fuse, n int, wg *sync.WaitGroup) {
		for i:=0; i<n; i++{
			f.trigger()
		}
		wg.Done()
	}

	for i:=0; i<runtime.NumCPU(); i++ {
		wg.Add(1)
		go triggerFuseTimesN(&f, 10 * 1000 * 1000 * 1000, &wg)
	}
	wg.Wait()

	if !f {
		t.Fatalf("fuse was not triggered!")
	}
}
