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
	action := ""
	if len(split) > 1 {
		action = split[1]
	}

	if command != "/agenda" {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}

	switch action {
	case "list":
		return p.executeCommandList(args), nil

	case "queue":
		return p.executeCommandQueue(args), nil

	case "setting":
		return p.executeCommandSetting(args), nil

	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown action: " + action),
	}, nil
}

func (p *Plugin) executeCommandList(args *model.CommandArgs) *model.CommandResponse {

	split := strings.Fields(args.Command)

	nextWeek := false
	if len(split) > 2 && split[2] == "next-week" {
		nextWeek = true
	}

	hashtag, error := p.GenerateHashtag(args.ChannelId, nextWeek)
	if error != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error calculating hashtags"),
		}
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
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Setting parameters missing"),
		}
	}

	field := split[2]
	value := split[3]

	meeting, err := p.GetMeeting(args.ChannelId)

	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error getting meeting information for this channel "),
		}
	}

	if field == "schedule" {
		//set schedule
		weekdayInt, _ := strconv.Atoi(value)
		meeting.Schedule = time.Weekday(weekdayInt)

	} else if field == "hashtag" {
		//set hashtag
		meeting.HashtagFormat = value
	} else {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknow setting " + field),
		}
	}

	if err := p.SaveMeeting(args.ChannelId, meeting); err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error saving setting "),
		}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Updated setting %v to %v", field, value),
	}
}

func (p *Plugin) executeCommandQueue(args *model.CommandArgs) *model.CommandResponse {
	split := strings.Fields(args.Command)

	if len(split) <= 2 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Missing parameters for queue command"),
		}
	}

	nextWeek := false
	message := strings.Join(split[2:], " ")

	if split[2] == "next-week" {
		nextWeek = true
		message = strings.Join(split[3:], " ")
	}

	hashtag, error := p.GenerateHashtag(args.ChannelId, nextWeek)
	if error != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error calculating hashtags"),
		}
	}

	user, appError := p.API.GetUser(args.UserId)

	if appError != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error getting user"),
		}
	}

	_, err := p.API.CreatePost(&model.Post{
		UserId:    p.botID,
		ChannelId: args.ChannelId,
		Message:   fmt.Sprintf("#### %v %v \n _%v @%v_", hashtag, message, "Queued by:", user.Username),
	})

	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error creating post: " + err.Message),
		}
	}

	return &model.CommandResponse{}
}
