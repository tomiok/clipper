package clipper

import (
	"errors"
	"net/http"
	"testing"
)

func BenchmarkDo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Do("my-circuit", func() error {
			_, err := http.Get("http://www.google.com/robots")
			return err
		}, nil)
	}
	b.ReportAllocs()
}

func TestDo(t *testing.T) {
	res := make(chan int, 1)

	ch := Do("my-circuit", func() error {

		res <- 1
		return nil
	}, nil)

	select {
	case v := <-ch:
		if v == 1 {
			t.Error("should be 0")
		}
		return
	}
}

func TestDo_failing(t *testing.T) {
	res := make(chan int, 1)

	ch := Do("my-circuit", func() error {

		res <- 1
		return errors.New("some error here")
	}, nil)

	select {
	case v := <-ch:
		if v != 1 {
			t.Error("should be 1")
		}
		return
	}
}
