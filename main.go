// tides - Print tidal data information for Aotearoa
// Copyright (C) 2021 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// read csv file and return map of string arrays
// skips the first 3 lines as they are comments in the LINZ csv files
func parseCSV(f *os.File) map[int][]string {
	entries := make(map[int][]string)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Failed reading file: %v", err)
	}
	r := csv.NewReader(strings.NewReader(string(data)))
	r.FieldsPerRecord = -1 // allows for variable number of fields
	// Skip the first 3 lines
	for i := 0; i < 3; i++ {
		r.Read()
	}
	i := 0
	for {
		entry, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		entries[i] = entry
		i++
	}
	return entries
}

func entryTime(year, month, day string) time.Time {
	loc, _ := time.LoadLocation("NZ") // Timezone isn't included in the CSV
	month = fmt.Sprintf("%02s", month)
	day = fmt.Sprintf("%02s", day)
	t, err := time.ParseInLocation("20060102", year+month+day, loc)
	if err != nil {
		fmt.Printf("failed parse date: %v\n", err)
		os.Exit(1)
	}
	return t
}

// print out tidal data
func display(entries map[int][]string) {
	nowTime := time.Now()
	fmt.Println(nowTime.Local())
	fmt.Println("---")
	for e := 0; e < len(entries); e++ {
		// fmt.Printf("%v ", e)
		// fmt.Println(entries[e])
		eTime := entryTime(entries[e][3], entries[e][2], entries[e][0])
		fmt.Println(eTime, entries[e][4], entries[e][5], entries[e][6], entries[e][7], entries[e][8], entries[e][9], entries[e][10], entries[e][11])
	}
}

func main() {
	entries := parseCSV(os.Stdin)
	display(entries)
}
