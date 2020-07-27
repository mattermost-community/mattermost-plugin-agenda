package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	scheduleErrorInvalid       = "Invalid weekday. Must be between 1-5 or Mon-Fri"
	scheduleErrorInvalidNumber = "Invalid weekday. Must be between 1-5"
)

var daysOfWeek = map[string]time.Weekday{}

func init() {
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := strings.ToLower(d.String())
		daysOfWeek[name] = d
		daysOfWeek[name[:3]] = d
	}
}

// parseSchedule will return a given Weekday based on the string being either a number,
// short / full name day of week.
func parseSchedule(val string) (time.Weekday, error) {
	if len(val) < 3 {
		return parseScheduleNumber(val)
	}
	if weekDayName, ok := daysOfWeek[strings.ToLower(val)]; ok {
		return weekDayName, nil
	}
	// try parsing number again in case prefixed by zeros
	weekDayInt, err := parseScheduleNumber(val)
	if err != nil {
		return -1, errors.New(scheduleErrorInvalid)
	}
	return weekDayInt, nil
}

// parseScheduleNumber will return a given Weekday based on the corresponding int val.
func parseScheduleNumber(val string) (time.Weekday, error) {
	weekdayInt, err := strconv.Atoi(val)
	validWeekday := weekdayInt >= 0 && weekdayInt <= 6
	if err != nil || !validWeekday {
		return -1, errors.New(scheduleErrorInvalidNumber)
	}
	return time.Weekday(weekdayInt), nil
}

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
