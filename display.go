// tides - Print tidal data information for Aotearoa
// Copyright (C) 2021 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"fmt"
	"math"
	"time"
)

func displayTerm(index int, tides *[]Tide, now time.Time) {
	prevTide := (*tides)[index-1]
	nextTide := (*tides)[index]
	height := getCurrentHeight(prevTide, nextTide, now)
	// nearest := getNearestTide(prevTide, nextTide, now)
	rising := getRising(prevTide, nextTide)
	fmt.Printf("%.2fm", height)
	if rising {
		fmt.Printf("⬆\n")
	} else {
		fmt.Printf("⬇\n")
	}
}

func displaySimple(index int, tides *[]Tide, now time.Time) {
	prevTide := (*tides)[index-1]
	nextTide := (*tides)[index]
	height := getCurrentHeight(prevTide, nextTide, now)
	rising := getRising(prevTide, nextTide)
	fmt.Printf("%.2fm", height)
	if rising {
		fmt.Printf("⬆\n")
	} else {
		fmt.Printf("⬇\n")
	}
}

// getCurrentHeight calculates the current tide height using the previous and
// future tide heights. The forumla comes from
// https://www.linz.govt.nz/sea/tides/tide-predictions/how-calculate-tide-times-heights
func getCurrentHeight(prev, next Tide, now time.Time) float64 {
	tf := getFloatTime(now)
	pf := getFloatTime(prev.Time)
	nf := getFloatTime(next.Time)
	ph := prev.Height
	nh := next.Height
	a := float64(math.Pi) * (((tf - pf) / (nf - pf)) + 1)
	h := ph + (nh-ph)*((math.Cos(a)+1)/2)
	return h
}

// getNearestTide returns either the previous or next tide for a given time
// depending on which is closer in time.
func getNearestTide(prev, next Tide, now time.Time) Tide {
	p := now.Sub(prev.Time)
	n := next.Time.Sub(now)
	if p < n {
		return prev
	} else {
		return next
	}
}

// getRising returns true if tide is rising based on the previous and next
// tides.
func getRising(prev, next Tide) bool {
	if next.Height > prev.Height {
		return true
	} else {
		return false
	}
}

// getFloatTime takes a time variable and returns a "decimal hour" format which
// is hour.minutes_as_percentage
func getFloatTime(t time.Time) float64 {
	h := float64(t.Hour())
	m := float64(t.Minute()) / 60
	return h + m
}
