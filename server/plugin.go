package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"

	fbClient "github.com/mattermost/focalboard/server/client"
	pluginapi "github.com/mattermost/mattermost-plugin-api"
)

const (
	// BotTokenKey is the KV store key for the bot's access token, used for the Focalboard API
	BotTokenKey = "bot_token"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	pluginAPI *pluginapi.Client

	fbStore FocalboardStore

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// BotId of the created bot account.
	botID string
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	case "/api/v1/queuedItems":
		p.httpQueuedItems(w, r)
	case "/api/v1/settings":
		p.httpMeetingSettings(w, r)
	case "/api/v1/meeting-days-autocomplete":
		p.httpMeetingDaysAutocomplete(w, r, false)
	case "/api/v1/list-meeting-days-autocomplete":
		p.httpMeetingDaysAutocomplete(w, r, true)
	default:
		http.NotFound(w, r)
	}
}

// OnActivate is invoked when the plugin is activated
func (p *Plugin) OnActivate() error {

	pluginAPIClient := pluginapi.NewClient(p.API, p.Driver)
	p.pluginAPI = pluginAPIClient

	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	botID, err := pluginAPIClient.Bot.EnsureBot(&model.Bot{
		Username:    "agenda",
		DisplayName: "Agenda Plugin Bot",
		Description: "Created by the Agenda plugin.",
		OwnerId:     "agenda",
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure agenda bot")
	}
	p.botID = botID

	token := ""
	rawToken, appErr := p.API.KVGet(BotTokenKey)
	if appErr != nil {
		return errors.Wrap(appErr, "failed to get stored bot access token")
	}
	if rawToken == nil {

		accessToken, appErr := p.API.CreateUserAccessToken(&model.UserAccessToken{UserId: botID, Description: "For agenda plugin access to focalboard REST API"})
		if appErr != nil {
			return errors.Wrap(appErr, "failed to create access token for bot")
		}
		token = accessToken.Token
		appErr = p.API.KVSet(BotTokenKey, []byte(token))
		if appErr != nil {
			return errors.Wrap(appErr, "failed to store bot access token")
		}
	} else {
		token = string(rawToken)
	}

	client := fbClient.NewClient("http://localhost:8065/plugins/focalboard", token)
	p.fbStore = NewFocalboardStore(p.API, client)

	return nil
}

func (p *Plugin) httpQueuedItems(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get("Mattermost-User-Id")
	if mattermostUserID == "" {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Request: "+r.Method+" is not allowed.", http.StatusMethodNotAllowed)
		return
	}

	channelID, ok := r.URL.Query()["channelId"]

	if !ok || len(channelID[0]) < 1 {
		http.Error(w, "Missing channelId parameter", http.StatusBadRequest)
		return
	}

	upNextCards, err := p.fbStore.GetUpnextCards(channelID[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.writeJSON(w, upNextCards)
}

func (p *Plugin) httpMeetingSettings(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get("Mattermost-User-Id")
	if mattermostUserID == "" {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}

	switch r.Method {
	case http.MethodPost:
		p.httpMeetingSaveSettings(w, r, mattermostUserID)
	case http.MethodGet:
		p.httpMeetingGetSettings(w, r, mattermostUserID)
	default:
		http.Error(w, "Request: "+r.Method+" is not allowed.", http.StatusMethodNotAllowed)
	}
}

func (p *Plugin) httpMeetingSaveSettings(w http.ResponseWriter, r *http.Request, mmUserID string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var meeting *Meeting
	if err = json.Unmarshal(body, &meeting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = p.SaveMeeting(meeting); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Status string
	}{"OK"}

	p.writeJSON(w, resp)
}

func (p *Plugin) httpMeetingGetSettings(w http.ResponseWriter, r *http.Request, mmUserID string) {
	channelID, ok := r.URL.Query()["channelId"]

	if !ok || len(channelID[0]) < 1 {
		http.Error(w, "Missing channelId parameter", http.StatusBadRequest)
		return
	}

	meeting, err := p.GetMeeting(channelID[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.writeJSON(w, meeting)
}

func (p *Plugin) writeJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		p.API.LogWarn("Failed to marshal JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		p.API.LogWarn("Failed to write JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) httpMeetingDaysAutocomplete(w http.ResponseWriter, r *http.Request, listCommand bool) {
	query := r.URL.Query()
	meeting, err := p.GetMeeting(query.Get("channel_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting meeting days: %s", err.Error()), http.StatusInternalServerError)
		p.API.LogDebug("Failed to find meeting for autocomplete", "error", err.Error(), "listCommand", listCommand)
		return
	}

	ret := make([]model.AutocompleteListItem, 0)

	helpText := "Queue this item "
	if listCommand {
		helpText = "List items "
	}

	for _, meetingDay := range meeting.Schedule {
		ret = append(ret, model.AutocompleteListItem{
			Item:     meetingDay.String(),
			HelpText: fmt.Sprintf(helpText+"for %s's meeting", meetingDay.String()),
			Hint:     "(optional)",
		})
	}
	ret = append(ret, model.AutocompleteListItem{
		Item:     "next-week",
		HelpText: fmt.Sprintf(helpText + "for the first meeting next week"),
		Hint:     "(optional)",
	})

	jsonBytes, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting meeting days: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(jsonBytes); err != nil {
		http.Error(w, fmt.Sprintf("Error getting meeting days: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
