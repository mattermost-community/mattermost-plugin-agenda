package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	meetingDateFormatRegex = regexp.MustCompile(`(?m)^(?P<prefix>.*)?(?:{{\s*(?P<dateformat>.*)\s*}})(?P<postfix>.*)?$`)
)

// Meeting represents a meeting agenda
type Meeting struct {
	ChannelID     string         `json:"channelId"`
	Schedule      []time.Weekday `json:"schedule"`
	HashtagFormat string         `json:"hashtagFormat"` // Default: {ChannelName}-Jan-2
}

// GetMeeting returns a meeting
func (p *Plugin) GetMeeting(channelID string) (*Meeting, error) {
	meetingBytes, appErr := p.API.KVGet(channelID)
	if appErr != nil {
		return nil, appErr
	}

	var meeting *Meeting
	if meetingBytes != nil {
		if err := json.Unmarshal(meetingBytes, &meeting); err != nil {
			return nil, err
		}
	} else {
		// Return a default value
		channel, err := p.API.GetChannel(channelID)
		if err != nil {
			return nil, err
		}
		paddedChannelName := strings.ReplaceAll(channel.Name, "-", "_")
		meeting = &Meeting{
			Schedule:      []time.Weekday{time.Thursday},
			HashtagFormat: strings.Join([]string{fmt.Sprintf("%.15s", paddedChannelName), "{{ Jan 2 }}"}, "_"),
			ChannelID:     channelID,
		}
	}

	return meeting, nil
}

// SaveMeeting saves a meeting
func (p *Plugin) SaveMeeting(meeting *Meeting) error {
	jsonMeeting, err := json.Marshal(meeting)
	if err != nil {
		return err
	}

	if appErr := p.API.KVSet(meeting.ChannelID, jsonMeeting); appErr != nil {
		return appErr
	}

	return nil
}

// GenerateHashtag returns a meeting hashtag
func (p *Plugin) GenerateHashtag(channelID string, nextWeek bool, weekday int, requeue bool, assignedDay time.Weekday) (string, error) {
	meeting, err := p.GetMeeting(channelID)
	if err != nil {
		return "", err
	}

	var meetingDate *time.Time
	if weekday > -1 {
		// Get date for given day
		if meetingDate, err = nextWeekdayDate(time.Weekday(weekday), nextWeek); err != nil {
			return "", err
		}
	} else {
		// user didn't provide any specific date, Get date for the list of days of the week
		if !requeue {
			if meetingDate, err = nextWeekdayDateInWeek(meeting.Schedule, nextWeek); err != nil {
				return "", err
			}
		} else {
			if len(meeting.Schedule) == 1 && meeting.Schedule[0] == assignedDay { // if this day is the only day selected in settings
				nextWeek = true
			}
			if meetingDate, err = nextWeekdayDateInWeekSkippingDay(meeting.Schedule, nextWeek, assignedDay); err != nil {
				return "", err
			}
		}
		//---- requeue Logic
	}

	var hashtag string

	if matchGroups := meetingDateFormatRegex.FindStringSubmatch(meeting.HashtagFormat); len(matchGroups) == 4 {
		var (
			prefix        string
			hashtagFormat string
			postfix       string
		)
		prefix = matchGroups[1]
		hashtagFormat = strings.TrimSpace(matchGroups[2])
		postfix = matchGroups[3]
		formattedDate := meetingDate.Format(hashtagFormat)
		formattedDate = strings.ReplaceAll(formattedDate, " ", "_")

		hashtag = fmt.Sprintf("#%s%v%s", prefix, formattedDate, postfix)
	} else {
		hashtag = fmt.Sprintf("#%s", meeting.HashtagFormat)
	}

	return hashtag, nil
}
