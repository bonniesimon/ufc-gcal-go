package main

func main() {
	events := GetMMAEvents()

	for i := 0; i < len(events); i++ {
		PrettyPrintEvent(events[i])
	}

	ShowCalendarEvents()
}
