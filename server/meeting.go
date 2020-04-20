package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Meeting represents a meeting agenda
type Meeting struct {
	ChannelID     string       `json:"channelId"`
	Schedule      time.Weekday `json:"schedule"`
	HashtagFormat string       `json:"hashtagFormat"` //Default: {ChannelName}-Jan02
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
		channel, err := p.API.GetChannel(channelID)
		if err != nil {
			return nil, err
		}
		meeting = &Meeting{
			Schedule:      time.Thursday,
			HashtagFormat: strings.Join([]string{fmt.Sprintf("%.15s", channel.Name), "Jan02"}, "-"),
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
func (p *Plugin) GenerateHashtag(channelID string, nextWeek bool) (string, error) {

	meeting, err := p.GetMeeting(channelID)
	if err != nil {
		return "", err
	}

	meetingDate := nextWeekdayDate(meeting.Schedule, nextWeek)

	hashtag := fmt.Sprintf("#%v", meetingDate.Format(meeting.HashtagFormat))

	return hashtag, nil
}
