package migrate

import "testing"

func TestWalkDir(t *testing.T) {
	ch := make(chan string)

	go func() {
		err := WalkDir(ch, DirKindUUID)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	for p := range ch {
		t.Logf("%s\n", p)
	}
}
