package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Meeting struct {
	ChannelId     string       `json:"channelId"`
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
		//Return a default value
		meeting = &Meeting{
			Schedule:      time.Thursday,
			HashtagFormat: "Jan02",
			ChannelId:     channelId,
		}
	}

	return meeting, nil
}

func (p *Plugin) SaveMeeting(meeting *Meeting) error {

	jsonMeeting, err := json.Marshal(meeting)
	if err != nil {
		return err
	}

	if err := p.API.KVSet(meeting.ChannelId, jsonMeeting); err != nil {
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
