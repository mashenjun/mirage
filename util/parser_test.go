package util

import (
	"testing"
	"time"
)

func TestCalSeconds(t *testing.T) {
	if secs, err := CalSeconds("07:30:00"); err != nil {
		t.Fatal(err)
	} else if secs != int64(7*3600+30*60) {
		t.Fatal("wrong seconds")
	}
	if secs, err := CalSeconds("15:30:00"); err != nil {
		t.Fatal(err)
	} else if secs != int64(15*3600+30*60) {
		t.Fatal("wrong seconds")
	}
	if secs, err := CalSeconds("23:30:00"); err != nil {
		t.Fatal(err)
	} else if secs != int64(23*3600+30*60) {
		t.Fatal("wrong seconds")
	}
}

func TestIsToday(t *testing.T) {
	today := time.Now()
	others := []time.Time{
		today.AddDate(1, 0, 0),
		today.AddDate(0, 1, 0),
		today.AddDate(0, 0, 1),
		today.AddDate(0, 0, -1),
	}

	if !IsToday(today) {
		t.Errorf("%+v is today", today)
	}
	for _, other := range others {
		if IsToday(other) {
			t.Errorf("%+v is not today", other)
		}
	}
}
