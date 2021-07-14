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
	nextDuration := fmtDuration(nextTide.Time.Sub(now))
	rising := getRising(prevTide, nextTide)
	height := getHeight(prevTide, nextTide, now)
	fmt.Printf("%.2fm", height)
	if rising {
		fmt.Printf("⬆ - high tide (%.2fm) in %v\n",
			nextTide.Height, nextDuration)
	} else {
		fmt.Printf("⬇ - low tide (%.2fm) in %v\n",
			nextTide.Height, nextDuration)
	}
}

func displaySimple(index int, tides *[]Tide, now time.Time) {
	prevTide := (*tides)[index-1]
	nextTide := (*tides)[index]
	rising := getRising(prevTide, nextTide)
	height := getHeight(prevTide, nextTide, now)
	fmt.Printf("%.2fm", height)
	if rising {
		fmt.Printf("⬆\n")
	} else {
		fmt.Printf("⬇\n")
	}
}

// getHeight calculates the tide height using the previous and future
// tide heights. The forumla comes from
// https://www.linz.govt.nz/sea/tides/tide-predictions/how-calculate-tide-times-heights
func getHeight(prev, next Tide, t time.Time) float64 {
	tf := getFloatTime(t)
	pf := getFloatTime(prev.Time)
	nf := getFloatTime(next.Time)
	ph := prev.Height
	nh := next.Height
	a := float64(math.Pi) * (((tf - pf) / (nf - pf)) + 1)
	h := ph + (nh-ph)*((math.Cos(a)+1)/2)
	return h
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

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh%02dm", h, m)
}
