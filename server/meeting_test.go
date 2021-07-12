package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
)

func assertNextWeekdayDate(meetingDay time.Weekday, nextWeek bool) *time.Time {
	weekDay, err := nextWeekdayDate(meetingDay, nextWeek)
	if err != nil {
		panic(err)
	}
	return weekDay
}

func TestPlugin_GenerateHashtag(t *testing.T) {
	tAssert := assert.New(t)
	mPlugin := Plugin{}
	api := &plugintest.API{}
	mPlugin.SetAPI(api)

	type args struct {
		nextWeek bool
		meeting  *Meeting
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "No Date Formatting",
			args: args{
				nextWeek: true,
				meeting: &Meeting{
					ChannelID:     "Developers",
					Schedule:      []time.Weekday{time.Thursday},
					HashtagFormat: "Developers",
				}},
			want:    "#Developers",
			wantErr: false,
		},
		{
			name: "Only Date Formatting",
			args: args{
				nextWeek: true,
				meeting: &Meeting{
					ChannelID:     "QA",
					Schedule:      []time.Weekday{time.Wednesday},
					HashtagFormat: "{{Jan 2}}",
				}},
			want:    "#" + strings.ReplaceAll(assertNextWeekdayDate(time.Wednesday, true).Format("Jan 2"), " ", "_"),
			wantErr: false,
		},
		{
			name: "Date Formatting with Prefix",
			args: args{
				nextWeek: true,
				meeting: &Meeting{
					ChannelID:     "QA Backend",
					Schedule:      []time.Weekday{time.Monday},
					HashtagFormat: "QA_{{January 02 2006}}",
				}},
			want:    "#QA_" + strings.ReplaceAll(assertNextWeekdayDate(time.Monday, true).Format("January 02 2006"), " ", "_"),
			wantErr: false,
		},
		{
			name: "Date Formatting with Postfix",
			args: args{
				nextWeek: false,
				meeting: &Meeting{
					ChannelID:     "QA FrontEnd",
					Schedule:      []time.Weekday{time.Monday},
					HashtagFormat: "{{January 02 2006}}.vue",
				}},
			want:    "#" + strings.ReplaceAll(assertNextWeekdayDate(time.Monday, false).Format("January 02 2006"), " ", "_") + ".vue",
			wantErr: false,
		},
		{
			name: "Date Formatting with Prefix and Postfix",
			args: args{
				nextWeek: false,
				meeting: &Meeting{
					ChannelID:     "QA Middleware",
					Schedule:      []time.Weekday{time.Monday},
					HashtagFormat: "React {{January 02 2006}} Born",
				}},
			want:    "#React " + strings.ReplaceAll(assertNextWeekdayDate(time.Monday, false).Format("January 02 2006"), " ", "_") + " Born",
			wantErr: false,
		},
		{
			name: "Date Formatting while ignoring Golang Time Formatting without brackets",
			args: args{
				nextWeek: false,
				meeting: &Meeting{
					ChannelID:     "Coffee Time",
					Schedule:      []time.Weekday{time.Monday},
					HashtagFormat: "January 02 2006 {{January 02 2006}} January 02 2006",
				}},
			want:    "#January 02 2006 " + strings.ReplaceAll(assertNextWeekdayDate(time.Monday, false).Format("January 02 2006"), " ", "_") + " January 02 2006",
			wantErr: false,
		},
		{
			name: "Date Formatting whitespace",
			args: args{
				nextWeek: false,
				meeting: &Meeting{
					ChannelID: "Dates with Spaces",
					Schedule:  []time.Weekday{time.Monday},
					HashtagFormat: "{{   January 02 2006			}}",
				}},
			want:    "#" + strings.ReplaceAll(assertNextWeekdayDate(time.Monday, false).Format("January 02 2006"), " ", "_"),
			wantErr: false,
		},
		{
			name: "Date Formatting ANSI",
			args: args{
				nextWeek: false,
				meeting: &Meeting{
					ChannelID:     "Dates",
					Schedule:      []time.Weekday{time.Monday},
					HashtagFormat: "{{ Mon Jan _2 }}",
				}},
			want:    "#" + strings.ReplaceAll(assertNextWeekdayDate(time.Monday, false).Format("Mon Jan _2"), " ", "_"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonMeeting, err := json.Marshal(tt.args.meeting)
			tAssert.Nil(err)
			api.On("KVGet", tt.args.meeting.ChannelID).Return(jsonMeeting, nil)
			got, err := mPlugin.GenerateHashtag(tt.args.meeting.ChannelID, tt.args.nextWeek, -1, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateHashtag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateHashtag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlugin_GetMeeting(t *testing.T) {
	tAssert := assert.New(t)
	mPlugin := Plugin{}
	api := &plugintest.API{}
	mPlugin.SetAPI(api)

	type args struct {
		channelID    string
		channelName  string
		storeMeeting *Meeting
	}
	tests := []struct {
		name    string
		args    args
		want    *Meeting
		wantErr bool
	}{
		{
			name: "Test Short Name",
			args: args{
				channelID:    "#short.name.channel",
				channelName:  "Short",
				storeMeeting: nil,
			},
			want: &Meeting{
				Schedule:      []time.Weekday{time.Thursday},
				HashtagFormat: "Short_{{ Jan 2 }}",
				ChannelID:     "#short.name.channel",
			},
			wantErr: false,
		},
		{
			name: "Test Log Name",
			args: args{
				channelID:    "#long.name.channel",
				channelName:  "Very Long Channel Name",
				storeMeeting: nil,
			},
			want: &Meeting{
				Schedule:      []time.Weekday{time.Thursday},
				HashtagFormat: "Very Long Chann_{{ Jan 2 }}",
				ChannelID:     "#long.name.channel",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.storeMeeting != nil {
				jsonMeeting, err := json.Marshal(tt.args.storeMeeting)
				tAssert.Nil(err)
				api.On("KVGet", tt.args.channelID).Return(jsonMeeting, nil)
			} else {
				api.On("KVGet", tt.args.channelID).Return(nil, nil)
			}
			api.On("GetChannel", tt.args.channelID).Return(GenerateFakeChannel(tt.args.channelID, tt.args.channelName))
			got, err := mPlugin.GetMeeting(tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMeeting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMeeting() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func GenerateFakeChannel(channelID, name string) (channel *model.Channel, appError *model.AppError) {
	channel = &model.Channel{
		Id:          channelID,
		DisplayName: strings.ToTitle(name),
		Name:        name,
		CreatorId:   "test",
	}
	return
}
