package dorkalonius_test

import (
  "errors"
  "testing"
)
import . "github.com/sethpollen/dorkalonius"

func TestSuccess(t *testing.T) {
  var calls int = 0
  m := NewMemo(func() (interface{}, error) {
    calls++
    return "foo", nil
  })
  for i := 0; i < 2; i++ {
    r, err := m.Get()
    if err != nil {
      t.Error(err)
    }
    if r.(string) != "foo" {
      t.Error(r)
    }
  }
  if calls != 1 {
    t.Error(calls)
  }
}

func TestFailure(t *testing.T) {
  var calls int = 0
  m := NewMemo(func() (interface{}, error) {
    calls++
    return "", errors.New("foo")
  })
  for i := 0; i < 2; i++ {
    _, err := m.Get()
    if err == nil {
      t.Error("Expected an error")
    }
  }
  if calls != 1 {
    t.Error(calls)
  }
}