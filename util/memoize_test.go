package util_test

import (
	"testing"
)
import . "github.com/sethpollen/dorkalonius/util"

func TestBasic(t *testing.T) {
	var calls int = 0
	m := NewMemo(func() interface{} {
		calls++
		return "foo"
	})
	for i := 0; i < 2; i++ {
		r := m.Get()
		if r.(string) != "foo" {
			t.Error(r)
		}
	}
	if calls != 1 {
		t.Error(calls)
	}
}
