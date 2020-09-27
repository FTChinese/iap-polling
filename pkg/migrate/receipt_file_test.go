package migrate

import (
	"runtime"
	"testing"
)

func TestWalkDir(t *testing.T) {
	ch := make(chan string)

	go func() {
		err := WalkDir(ch, "iap_receipts")
		if err != nil {
			t.Error(err)
			return
		}
	}()

	for p := range ch {
		t.Logf("%s\n", p)
	}
}

func TestMaxWorkers(t *testing.T) {
	t.Logf("%d", runtime.GOMAXPROCS(0))
}
