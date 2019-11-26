package main

import (
	"time"
)

func nextWeekdayDate(day time.Weekday, nextWeek bool) time.Time {

	daysTill := daysTillNextWeekday(time.Now().Weekday(), day, nextWeek)

	return time.Now().AddDate(0, 0, daysTill)
}

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
