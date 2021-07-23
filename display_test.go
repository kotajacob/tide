// tides - Print tidal data information for Aotearoa
// Copyright (C) 2021 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"bytes"
	"testing"
	"time"
)

var (
	referenceTime = time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)
)

func TestDisplaySimple(t *testing.T) {
	tides := []Tide{
		{
			Height: 2.0,
			Time:   referenceTime,
		},
		{
			Height: 0.0,
			Time:   referenceTime.Add(time.Hour * 8),
		},
		{
			Height: 2.0,
			Time:   referenceTime.Add(time.Hour * 16),
		},
	}
	var tests = []struct {
		time time.Time
		want string
	}{
		{referenceTime, "2.00m⬇\n"},
		{referenceTime.Add(time.Hour * 2), "1.71m⬇\n"},
		{referenceTime.Add(time.Hour * 4), "1.00m⬇\n"},
		{referenceTime.Add(time.Hour * 8), "0.00m⬇\n"},
	}
	for _, test := range tests {
		out = new(bytes.Buffer) // captured output
		displaySimple(1, &tides, test.time)
		got := out.(*bytes.Buffer).String()
		if got != test.want {
			t.Errorf("got = %q\nwant = %q", got, test.want)
		}
	}
}
