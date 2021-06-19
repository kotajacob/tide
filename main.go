// tides - Print tidal data information for Aotearoa
// Copyright (C) 2021 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"
)

// getTides reads csv file from LINZ with tidal data and returns a map of times
// with tide heights. The first 3 lines are metadata and are thus skipped.
func getTides(f *os.File, tides *map[time.Time]float32) error {
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // allows for variable number of fields
	// skip the first 3 lines
	for i := 0; i < 3; i++ {
		r.Read()
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		addTide(tides, record)
	}
	return nil
}

// addTide takes a csv record and appends tide entries to a map.
func addTide(tides *map[time.Time]float32, records []string) error {
	// each record represents a single date, but contains multiple tides
	date, err := getDate(records[3], records[2], records[0])
	if err != nil {
		return err
	}
	for r := 4; r < len(records); r += 2 {
		// some days have less tides
		if records[r] == "" {
			break
		}
		// reformat time into time.Duration
		f := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
		h := strings.FieldsFunc(records[r], f)
		duration, err := time.ParseDuration(fmt.Sprintf("%vh%vm", h[0], h[1]))
		time := date.Add(duration)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return err
		}
		fmt.Println(time, records[r+1])
	}
	return nil
}

// getDate takes the year, month, day strings from the CSV file and returns a
// time.Time value with the correct timezone.
func getDate(year, month, day string) (time.Time, error) {
	loc, _ := time.LoadLocation("NZ") // Timezone isn't included in the CSV
	month = fmt.Sprintf("%02s", month)
	day = fmt.Sprintf("%02s", day)
	t, err := time.ParseInLocation("20060102", year+month+day, loc)
	return t, err
}

func main() {
	tides := make(map[time.Time]float32)
	err := getTides(os.Stdin, &tides)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
