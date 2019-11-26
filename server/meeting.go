package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Meeting struct {
	Schedule      time.Weekday `json:"schedule"`
	HashtagFormat string       `json:"hashtagFormat"` //Default: Jan02
}

func (p *Plugin) GetMeeting(channelId string) (*Meeting, error) {

	meettingBytes, appErr := p.API.KVGet(channelId)
	if appErr != nil {
		return nil, appErr
	}

	var meeting *Meeting
	if meettingBytes != nil {

		if err := json.Unmarshal(meettingBytes, &meeting); err != nil {
			return nil, err
		}

	} else {
		//Default values
		meeting = &Meeting{
			Schedule:      time.Thursday,
			HashtagFormat: "Jan02",
		}
	}

	return meeting, nil
}

func (p *Plugin) SaveMeeting(channelId string, meeting *Meeting) error {

	jsonMeeting, err := json.Marshal(meeting)
	if err != nil {
		return err
	}

	if err := p.API.KVSet(channelId, jsonMeeting); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) GenerateHashtag(channelId string, nextWeek bool) (string, error) {

	meeting, err := p.GetMeeting(channelId)
	if err != nil {
		return "", err
	}

	meetingDate := nextWeekdayDate(meeting.Schedule, nextWeek)

	hastag := fmt.Sprintf("#%v", meetingDate.Format(meeting.HashtagFormat))

	return hastag, nil
}
