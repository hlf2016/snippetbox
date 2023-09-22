package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2023, 9, 22, 14, 18, 0, 0, time.UTC)
	hd := humanDate(tm)
	if hd != "22 Sep 2023 at 14:18" {
		t.Errorf("got %q; want %q", hd, "22 Sep 2023 at 14:18")
	}
}
