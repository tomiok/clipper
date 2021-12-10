package clipper

import (
	"testing"
	"time"
)

func Test_rand(t *testing.T) {
	s1 := randStr()
	time.Sleep(1 * time.Second)
	s2 := randStr()

	if s1 == s2 {
		t.Fail()
	}
}
