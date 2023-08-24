# Agenda Plugin

[![CircleCI](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-agenda/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-agenda)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattermost/mattermost-plugin-agenda)](https://goreportcard.com/report/github.com/mattermost/mattermost-plugin-agenda)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-agenda/master.svg)](https://codecov.io/gh/mattermost/mattermost-plugin-agenda)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-agenda?include_prereleases)](https://github.com/mattermost/mattermost-plugin-agenda/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-agenda/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-agenda/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

**Maintainer:** [@jfrerich](https://github.com/jfrerich)
**Co-Maintainer:** [@iomodo](https://github.com/iomodo)

The Agenda Plugin helps users queue and list items in a channel's meeting agenda. The agenda is identified by a hashtag based on the meeting date.

The plugin will create posts for the user preceding the agenda item with configured hashtag format and can open a search with that hashtag to view the agenda list. 

Initial development as part of [Mattermost Hackathon 2019](https://github.com/mattermost/mattermost-hackathon-nov2019) which was demoed [here](https://www.youtube.com/watch?v=Tl08dt7TheI&feature=youtu.be&t=821).

## Usage

### Enable the plugin

Once this plugin is installed, a Mattermost admin can enable it in the Mattermost System Console by going to **Plugins > Plugin Management**, and selecting **Enable**.

### Configure meeting settings

The meeting settings for each channel can be configured in the Channel Header Dropdown.

![channel_header_menu](./assets/channelHeaderDropdown.png)

![settings_dialog](./assets/settingsDialog.png)

Meeting settings include:

- Schedule Day: Day of the week when the meeting is scheduled.
- Hashtag Format: The format of the hashtag for the meeting date. The date format is based on [Go date and time formatting](https://yourbasic.org/golang/format-parse-string-time-date-example/#standard-time-and-date-formats).
  The date format must be wrapped in double Braces ( {{ }} ).
  A default is generated from the first 15 characters of the channel's name with the short name of the month and day (i.e. Dev-{{ Jan02 }}).

#### Slash Commands to manage the meeting agenda

```
/agenda queue [meetingDay] message
```
Creates a post for the user with the given `message` for the next meeting date or the specified `meetingDay` (optional). The configured hashtag will precede the `message`.
The meeting day supports long (Monday, Tuesday), short name (Mon Tue), number (0-6) or `next-week`. If `next-week` is indicated, it will use the date of the first meeting in the next calendar week. 

![post_example](./assets/postExample.png)

```
/agenda list [meetingDay]
```
Executes a search of the hashtag of the next meeting or the specified `meetingDay` (optional), opening the RHS with all the posts with that hashtag. 
The meeting day supports long (Monday, Tuesday), short name (Mon Tue), number (0-6) or `next-week`. If `next-week` is indicated, it will use the date of the first meeting in the next calendar week. 

```
/agenda setting field value
```
Updates the given setting with the provided value for the meeting settings of that channel. 

`Field` can be one of:

- `schedule`: Day of the week of the meeting. It is an int based on [`time.Weekday`](https://golang.org/pkg/time/#Weekday)
- `hashtag`: Format of the hashtag for the meeting date. It is based on the format used in [`time.Format`](https://golang.org/pkg/time/#Time.Format)

## Future Improvements

- Mark items as resolved or queue for next week. 
- Queue a post using a menu option in the post dot menu. 
- Handle time in meeting schedule. 

## Contributing

If you would like to make contributions to this plugin, please checkout the open issues labeled [`Help Wanted` and `Up For Grabs`](https://github.com/mattermost/mattermost-plugin-agenda/issues?q=is%3Aopen+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)
