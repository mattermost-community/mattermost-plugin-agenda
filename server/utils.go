package main

import (
	"time"
)

// nextWeekdayDate calculates the date of the next given weekday
// based on today's date.
// If nextWeek is true, it will be based on the next calendar week.
func nextWeekdayDate(meetingDays []time.Weekday, nextWeek bool) time.Time {
	todayWeekday := time.Now().Weekday()

	// Find which meeting weekday to calculate the date for
	meetingDay := meetingDays[0]
	for _, day := range meetingDays {
		if todayWeekday <= day {
			meetingDay = day
			break
		}
	}

	daysTill := daysTillNextWeekday(todayWeekday, meetingDay, nextWeek)

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
