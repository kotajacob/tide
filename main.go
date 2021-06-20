// tides - Print tidal data information for Aotearoa
// Copyright (C) 2021 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/guptarohit/asciigraph"
)

type Tides map[time.Time]float64

func (t *Tides) Set(time time.Time, height float64) *Tides {
	(*t)[time] = height
	return t
}

// getTides reads and parses a csv file from LINZ with tidal data into
// [][]string and skips the first 3 metadata lines.
func getRecords(f *os.File) ([][]string, error) {
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // allows for variable number of fields
	// skip the first 3 lines
	for i := 0; i < 3; i++ {
		reader.Read()
	}
	records, err := reader.ReadAll()
	return records, err
}

// parseRecord reads a []string representing a line in the csv file
func parseRecord(tides *Tides, times *[]time.Time, record []string) error {
	// Each record represents a single date, but contains multiple tides at
	// different times.
	date, err := getDate(record[3], record[2], record[0])
	if err != nil {
		return err
	}
	for r := 4; r < len(record); r += 2 {
		// some days have less tides
		if record[r] == "" {
			break
		}
		duration, err := getDuration(record[r])
		t := date.Add(duration)
		if err != nil {
			return err
		}
		height, err := strconv.ParseFloat(record[r+1], 64)
		if err != nil {
			return err
		}
		tides.Set(t, height)
		*times = append(*times, t)
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

// getDuration takes a string in the hh:mm format and returns a time.Duration
func getDuration(s string) (time.Duration, error) {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	t := strings.FieldsFunc(s, f)
	duration, err := time.ParseDuration(fmt.Sprintf("%vh%vm", t[0], t[1]))
	return duration, err
}

func display(index int, tides *Tides, times *[]time.Time) {
	var heights []float64
	heights = append(heights, (*tides)[(*times)[index-1]])
	heights = append(heights, (*tides)[(*times)[index]])
	heights = append(heights, (*tides)[(*times)[index+1]])
	heights = append(heights, (*tides)[(*times)[index+2]])
	heights = append(heights, (*tides)[(*times)[index+3]])
	heights = append(heights, (*tides)[(*times)[index+4]])
	heights = append(heights, (*tides)[(*times)[index+5]])
	graph := asciigraph.Plot(heights, asciigraph.Height(8))
	fmt.Println(graph)
	fmt.Println(heights)
	fmt.Printf("%v - %v\n", (*times)[index-1], (*tides)[(*times)[index-1]])
	fmt.Printf("%v - %v\n", (*times)[index], (*tides)[(*times)[index]])
	fmt.Printf("%v - %v\n", (*times)[index+1], (*tides)[(*times)[index+1]])
	fmt.Printf("%v - %v\n", (*times)[index+2], (*tides)[(*times)[index+2]])
	fmt.Printf("%v - %v\n", (*times)[index+3], (*tides)[(*times)[index+3]])
	fmt.Printf("%v - %v\n", (*times)[index+4], (*tides)[(*times)[index+4]])
	fmt.Printf("%v - %v\n", (*times)[index+5], (*tides)[(*times)[index+5]])
}

func main() {
	now := time.Now()
	tides := &Tides{}
	var times []time.Time
	records, err := getRecords(os.Stdin)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	for _, record := range records {
		err := parseRecord(tides, &times, record)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}

	for k, v := range times {
		if v.After(now) {
			display(k, tides, &times)
			break
		}
	}
}
