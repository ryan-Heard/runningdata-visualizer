package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

type RunningData struct {
	Date       string  `csv:"Date,omitempty"`
	Distance   float64 `csv:"Distance,omitempty"`
	AvgHR      int64   `csv:"Avg HR,omitempty"`
	AvgPaceRaw string  `csv:"Avg Pace,omitempty"`
	AvgPace    time.Duration
	AvgStride  float64 `csv:"Avg Stride Length,omitempty"`
}

func ReadCSVtoRunData(s string) []*RunningData {
	raw := []*RunningData{}
	dataset := make([]*RunningData, 0)

	f, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if err := gocsv.UnmarshalFile(f, &raw); err != nil {
		panic(err)
	}

	// I want series runs aand for me thats over 1.1 mile
	for _, d := range raw {
		if d.Distance > 1.1 {
			t, err := ParseTime(d.AvgPaceRaw)
			if err != nil {
				log.Println("Setting Pace to 0. Invalid data on ", d.Date)
				d.AvgPace = 0
			}

			d.AvgPace = t
			dataset = append(dataset, d)
		}
	}

	return dataset
}

func ParseTime(t string) (time.Duration, error) {
	var hours, mins, seconds int
	var err error

	parts := strings.SplitN(t, ":", 3)

	// TODO streamline time more
	switch len(parts) {
	case 1:
		seconds, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
	case 2:
		mins, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		seconds, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
	case 3:
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		mins, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}

		seconds, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}

	default:
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	if seconds > 59 || seconds < 0 || mins > 59 || mins < 0 || hours > 23 || hours < 0 {
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	return time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute + time.Duration(seconds)*time.Second, nil
}
