package main

import (
	"fmt"
	"time"
)

func convertETtoIST(dateStr string) (time.Time, error) {
	t, err := time.Parse("Monday, January  2,  3:04 PM ET", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %v", err)
	}

	t = t.AddDate(time.Now().Year(), 0, 0)

	et, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading ET location: %v", err)
	}

	tET := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, et)

	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading IST location: %v", err)
	}

	return tET.In(ist), nil
}
