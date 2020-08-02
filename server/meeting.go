package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Meeting represents a meeting agenda
type Meeting struct {
	ChannelID     string         `json:"channelId"`
	Schedule      []time.Weekday `json:"schedule"`
	HashtagFormat string         `json:"hashtagFormat"` //Default: Jan02
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
		//Return a default value
		meeting = &Meeting{
			Schedule:      []time.Weekday{time.Thursday},
			HashtagFormat: "Jan02",
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
func (p *Plugin) GenerateHashtag(channelID string, nextWeek bool, weekday int) (string, error) {

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
		// Get date for the list of days of the week
		if meetingDate, err = nextWeekdayDateInWeek(meeting.Schedule, nextWeek); err != nil {
			return "", err
		}
	}

	hashtag := fmt.Sprintf("#%v", meetingDate.Format(meeting.HashtagFormat))

	return hashtag, nil
}
