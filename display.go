// tides - Print tidal data information for Aotearoa
// Copyright (C) 2021 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/wzshiming/ctc"
)

const graphWidth int = 36
const graphHeight int = 8

func displayTerm(index int, tides *[]Tide, now time.Time) {
	prevTide := (*tides)[index-1]
	nextTide := (*tides)[index]
	nextDuration := fmtDuration(nextTide.Time.Sub(now))
	rising := getRising(prevTide, nextTide)
	height := getHeight(prevTide, nextTide, now)
	fmt.Fprintf(out, "%.2fm", height)
	if rising {
		fmt.Fprintf(out, "⬆ - high tide (%.2fm) in %v\n",
			nextTide.Height, nextDuration)
	} else {
		fmt.Fprintf(out, "⬇ - low tide (%.2fm) in %v\n",
			nextTide.Height, nextDuration)
	}
	fmt.Fprintf(out, "%s", graph(prevTide, nextTide, now))
}

func displaySimple(index int, tides *[]Tide, now time.Time) {
	prevTide := (*tides)[index-1]
	nextTide := (*tides)[index]
	rising := getRising(prevTide, nextTide)
	height := getHeight(prevTide, nextTide, now)
	fmt.Fprintf(out, "%.2fm", height)
	if rising {
		fmt.Fprintf(out, "⬆\n")
	} else {
		fmt.Fprintf(out, "⬇\n")
	}
}

// graph returns a printable string with a wave graph of the tide heights using
// the previous and next tides, and the current time. The last graph point will
// be the next Tide, the first will be after the previous tide.
func graph(prev, next Tide, now time.Time) string {
	timeInterval := next.Time.Sub(prev.Time) / time.Duration(graphWidth)
	minHeight := math.Min(next.Height, prev.Height)
	maxHeight := math.Max(next.Height, prev.Height)
	var waves [graphWidth][graphHeight]string
	for x, w := range waves {
		d := timeInterval * time.Duration(x+1)
		t := prev.Time.Add(d)
		h := getHeight(prev, next, t)
		scaledHeight := scaleDatum(h, minHeight, maxHeight, graphHeight)
		for y := range w {
			waves[x][y] = " "
		}
		if now.After(t) {
			waves[x][scaledHeight] = fmt.Sprintf("%s%s%s",
				ctc.ForegroundBrightBlack, "█", ctc.Reset)
		} else {
			waves[x][scaledHeight] = "█"
		}
	}
	// build the print string
	var s string
	for y := graphHeight - 1; y >= 0; y-- {
		for x := 0; x < graphWidth; x++ {
			s += waves[x][y]
		}
		s += "\n"
	}
	return s
}

// scaleDatum scales a point of data to the closest 'degrees' space between min
// and max. e.g. if datum = 14, min = 10, and max = 20, and degrees = 5, then
// this function would output '2', meaning 14 is 2/5 of the way between 10 and
// 20.  The actual output is floored and 0 indexed, so it will be a number from
// 0-4. A special case is needed in the datum can equal the max.
func scaleDatum(datum float64, min float64, max float64, degrees int) int {
	if datum == max {
		return degrees - 1
	}
	proportion := (datum - min) / (max - min)
	scaled := proportion * float64(degrees)
	return int(scaled)
}

// getHeight calculates the tide height using the previous and future
// tide heights. The forumla comes from
// https://www.linz.govt.nz/sea/tides/tide-predictions/how-calculate-tide-times-heights
func getHeight(prev, next Tide, t time.Time) float64 {
	tf := float64(t.Unix())
	pf := float64(prev.Time.Unix())
	nf := float64(next.Time.Unix())
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

// fmtDuration returns a formatted string with hours and minutes.
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh%02dm", h, m)
}
