package apple

import (
	"runtime"
	"testing"
)

func TestWalkDir(t *testing.T) {
	ch := WalkReceipts("iap_receipts")

	for p := range ch {
		t.Logf("%s\n", p)
	}
}

func TestMaxWorkers(t *testing.T) {
	t.Logf("%d", runtime.GOMAXPROCS(0))
}
