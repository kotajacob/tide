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
	i := -3 // skip first 3 records
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

// print out tidal data
func display(entries map[int][]string) {
	for e := 0; e < len(entries); e++ {
		fmt.Println(entries[e])
	}
}

func main() {
	entries := parseCSV(os.Stdin)
	display(entries)
}
