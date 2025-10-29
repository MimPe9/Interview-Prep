package storage

import "testing"

func SimpleTset(t *testing.T) {
	if 2*2 != 4 {
		t.Errorf("Math is broken")
	}
}
