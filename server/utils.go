package main

import (
	"time"
)

// nextWeekdayDate calculates the date of the next given weekday
// based on today's date.
// If nextWeek is true, it will be based on the next calendar week.
func nextWeekdayDate(day time.Weekday, nextWeek bool) time.Time {

	daysTill := daysTillNextWeekday(time.Now().Weekday(), day, nextWeek)

	return time.Now().AddDate(0, 0, daysTill)
}

/// daysTillNextWeekday calculates the amount of days between two weekdays.
/// If nexWeek is true, the nextDay will be based on the next calendar week.
func daysTillNextWeekday(today time.Weekday, nextDay time.Weekday, nextWeek bool) int {

	if today > nextDay {
		return int((7 - today) + nextDay)
	}

	daysTillNextWeekday := int(nextDay - today)

	if nextWeek {
		daysTillNextWeekday += 7
	}

	return daysTillNextWeekday
}
