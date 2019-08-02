package main

import "testing"

func TestEmpty(t *testing.T) {
    var s string = ""
    e := empty(s)
    if !e {
        t.Error("Expected true, got ", e)
    }
}
