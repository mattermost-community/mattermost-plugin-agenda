package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	commandTriggerAgenda = "agenda"

	WS_EVENT_LIST = "list"
)

func (p *Plugin) registerCommands() error {
	if err := p.API.RegisterCommand(&model.Command{

		Trigger:          commandTriggerAgenda,
		AutoComplete:     true,
		AutoCompleteHint: "[command]",
		AutoCompleteDesc: "Available commands: list, queue, setting",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerAgenda)
	}

	return nil
}

// ExecuteCommand
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]

	if command != "/agenda" {
		return responsef("Unknown command: " + args.Command), nil
	}

	if len(split) < 2 {
		return responsef("Missing command. You can try queue, list, setting"), nil
	}

	action := split[1]

	switch action {
	case "list":
		return p.executeCommandList(args), nil

	case "queue":
		return p.executeCommandQueue(args), nil

	case "setting":
		return p.executeCommandSetting(args), nil

	}

	return responsef("Unknown action: " + action), nil
}

func (p *Plugin) executeCommandList(args *model.CommandArgs) *model.CommandResponse {

	split := strings.Fields(args.Command)
	nextWeek := len(split) > 2 && split[2] == "next-week"

	hashtag, error := p.GenerateHashtag(args.ChannelId, nextWeek)
	if error != nil {
		return responsef("Error calculating hashtags")
	}

	p.API.PublishWebSocketEvent(
		WS_EVENT_LIST,
		map[string]interface{}{
			"hashtag": hashtag,
		},
		&model.WebsocketBroadcast{UserId: args.UserId},
	)

	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandSetting(args *model.CommandArgs) *model.CommandResponse {

	// settings: hashtag, schedule
	split := strings.Fields(args.Command)

	if len(split) < 4 {
		return responsef("Setting parameters missing")
	}

	field := split[2]
	value := split[3]

	meeting, err := p.GetMeeting(args.ChannelId)
	if err != nil {
		return responsef("Error getting meeting information for this channel")
	}

	if field == "schedule" {
		//set schedule
		weekdayInt, err := strconv.Atoi(value)
		if err != nil {
			return responsef("Invalid weekday. Must be between 1-5")
		}

		meeting.Schedule = time.Weekday(weekdayInt)

	} else if field == "hashtag" {
		//set hashtag
		meeting.HashtagFormat = value
	} else {
		return responsef("Unknow setting " + field)
	}

	if err := p.SaveMeeting(meeting); err != nil {
		return responsef("Error saving setting")
	}

	return responsef("Updated setting %v to %v", field, value)
}

func (p *Plugin) executeCommandQueue(args *model.CommandArgs) *model.CommandResponse {
	split := strings.Fields(args.Command)

	if len(split) <= 2 {
		return responsef("Missing parameters for queue command")
	}

	nextWeek := false
	message := strings.Join(split[2:], " ")

	if split[2] == "next-week" {
		nextWeek = true
		message = strings.Join(split[3:], " ")
	}

	hashtag, error := p.GenerateHashtag(args.ChannelId, nextWeek)
	if error != nil {
		return responsef("Error calculating hashtags")
	}

	itemsQueued, appError := p.API.SearchPostsInTeam(args.TeamId, []*model.SearchParams{{Terms: hashtag, IsHashtag: true}})

	if appError != nil {
		return responsef("Error getting user")
	}

	_, err := p.API.CreatePost(&model.Post{
		UserId:    args.UserId,
		ChannelId: args.ChannelId,
		Message:   fmt.Sprintf("#### %v %v) %v", hashtag, len(itemsQueued)+1, message),
	})
	if err != nil {
		return responsef("Error creating post: " + err.Message)
	}

	return &model.CommandResponse{}
}

func responsef(format string, args ...interface{}) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf(format, args...),
		Type:         model.POST_DEFAULT,
	}
}
