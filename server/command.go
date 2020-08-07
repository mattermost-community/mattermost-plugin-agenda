package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	commandTriggerAgenda = "agenda"

	wsEventList = "list"
)

const helpCommandText = "###### Mattermost Agenda Plugin - Slash Command Help\n" +
	"\n* `/agenda queue [weekday (optional)] message` - Queue `message` as a topic on the next meeting. If `weekday` is provided, it will queue for the meeting for. \n" +
	"* `/agenda list [weekday(optional)]` - Show a list of items queued for the next meeting.  If `next-week` is provided, it will list the agenda for the next calendar week. \n" +
	"* `/agenda setting <field> <value>` - Update the setting with the given value. Field can be one of `schedule` or `hashtag` \n"

func (p *Plugin) registerCommands() error {
	if err := p.API.RegisterCommand(createAgendaCommand()); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerAgenda)
	}

	return nil
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)

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

	case "help":
		return p.executeCommandHelp(args), nil
	}

	return responsef("Unknown action: %s", action), nil
}

func (p *Plugin) executeCommandList(args *model.CommandArgs) *model.CommandResponse {
	split := strings.Fields(args.Command)
	nextWeek := len(split) > 2 && split[2] == "next-week"

	weekday := -1
	if !nextWeek && len(split) > 2 {
		parsedWeekday, _ := parseSchedule(split[2])
		weekday = int(parsedWeekday)
	}

	hashtag, err := p.GenerateHashtag(args.ChannelId, nextWeek, weekday)
	if err != nil {
		return responsef("Error calculating hashtags")
	}

	// Send a websocket event to the web app that will open the RHS
	p.API.PublishWebSocketEvent(
		wsEventList,
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

	switch field {
	case "schedule":
		//set schedule
		weekdayInt, err := parseSchedule(value)
		if err != nil {
			return responsef(err.Error())
		}
		meeting.Schedule = []time.Weekday{weekdayInt}

	case "hashtag":
		//set hashtag
		meeting.HashtagFormat = value
	default:
		return responsef("Unknown setting %s", field)
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
	weekday := -1
	message := strings.Join(split[2:], " ")

	if split[2] == "next-week" {
		nextWeek = true
	} else {
		parsedWeekday, _ := parseSchedule(split[2])
		weekday = int(parsedWeekday)
	}

	if nextWeek || weekday > -1 {
		message = strings.Join(split[3:], " ")
	}

	hashtag, error := p.GenerateHashtag(args.ChannelId, nextWeek, weekday)
	if error != nil {
		return responsef("Error calculating hashtags. Check the meeting settings for this channel.")
	}

	searchResults, appErr := p.API.SearchPostsInTeamForUser(args.TeamId, args.UserId, model.SearchParameter{Terms: &hashtag})

	if appErr != nil {
		return responsef("Error calculating list number")
	}

	postList := *searchResults.PostList
	numQueueItems := len(postList.Posts)

	_, appErr = p.API.CreatePost(&model.Post{
		UserId:    args.UserId,
		ChannelId: args.ChannelId,
		Message:   fmt.Sprintf("#### %v %v) %v", hashtag, numQueueItems+1, message),
	})
	if appErr != nil {
		return responsef("Error creating post: %s", appErr.Message)
	}

	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandHelp(args *model.CommandArgs) *model.CommandResponse {
	return responsef(helpCommandText)
}

func responsef(format string, args ...interface{}) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf(format, args...),
		Type:         model.POST_DEFAULT,
	}
}

func createAgendaCommand() *model.Command {
	agenda := model.NewAutocompleteData(commandTriggerAgenda, "[command]", "Available commands: list, queue, setting, help")

	list := model.NewAutocompleteData("list", "", "Show a list of items queued for the next meeting")
	list.AddDynamicListArgument("Day of the week for when to queue the meeting", "/api/v1/list-meeting-days-autocomplete", false)
	agenda.AddCommand(list)

	queue := model.NewAutocompleteData("queue", "", "Queue `message` as a topic on the next meeting.")
	queue.AddDynamicListArgument("Day of the week for when to queue the meeting", "/api/v1/meeting-days-autocomplete", false)
	queue.AddTextArgument("Message for the next meeting date.", "[message]", "")
	agenda.AddCommand(queue)

	setting := model.NewAutocompleteData("setting", "", "Update the setting.")
	schedule := model.NewAutocompleteData("schedule", "", "Update schedule.")
	schedule.AddStaticListArgument("weekday", true, []model.AutocompleteListItem{
		{Item: "Monday"},
		{Item: "Tuesday"},
		{Item: "Wednesday"},
		{Item: "Thursday"},
		{Item: "Friday"},
		{Item: "Saturday"},
		{Item: "Sunday"},
	})
	setting.AddCommand(schedule)
	hashtag := model.NewAutocompleteData("hashtag", "", "Update hastag.")
	hashtag.AddTextArgument("input hashtag", "Default: Jan02", "")
	setting.AddCommand(hashtag)
	agenda.AddCommand(setting)

	help := model.NewAutocompleteData("help", "", "Mattermost Agenda plugin slash command help")
	agenda.AddCommand(help)
	return &model.Command{
		Trigger:          commandTriggerAgenda,
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: list, queue, setting, help",
		AutoCompleteHint: "[command]",
		AutocompleteData: agenda,
	}
}
