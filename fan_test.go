package fan_test

import (
	"sync"
	"testing"

	"github.com/DeedleFake/fan"
)

func TestFan(t *testing.T) {
	var f fan.Fan

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			r := f.Receiver()

			t.Logf("%v got %q", i, r.Get())
		}(i)
	}

	const str = "test"
	for i := range str {
		f.Send(str[:i])
	}
	wg.Wait()
}
