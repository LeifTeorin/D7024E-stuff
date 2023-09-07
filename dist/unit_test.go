package d7024e

import (
	"testing"
)

func test_test(t *testing.T){
	got := 1
	want:= 1

	if got != want {
		t.Errorf("Got %q, wanted %q", got, want);
	}
}
