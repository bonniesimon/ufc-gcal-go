package main

import "time"

func UNUSED(x ...interface{}) {}

func main() {
	events := GetMMAEvents()

	calendar := NewGoogleCalendar()
	calendar.ShowCalendarEvents()

	for i := 0; i < 1; i++ {
		event := events[i]
		PrettyPrintEvent(event)

		fourAndHalfHoursDuration, _ := time.ParseDuration("4h30m")

		calendar.AddCalendarEvent(
			event.Title,
			GetEventAsString(event),
			event.DateTime,
			event.DateTime.Add(fourAndHalfHoursDuration),
		)
	}
}
