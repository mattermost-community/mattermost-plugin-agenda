# Agenda Plugin

The Agenda Plugin helps users queue and list items in a channel's meeting agenda. The agenda is identified by a hashtag based on the meeting date.

The plugin will create posts for the user preceding the agenda item with configured hashtag format and can open a search with that hashtag to view the agenda list. 

Initial development as part of [Mattermost Hackathon 2019](https://github.com/mattermost/mattermost-hackathon-nov2019)

## Usage

#### Meeting Settings Configuration

The meeting settings for each channel can be configured in the Channel Header Dropwdown (supported in [this WebApp branch](https://github.com/mattermost/mattermost-webapp/tree/MM-19902))

![channel_header_menu](./assets/channelHeaderDropdown.png)

![settings_dialog](./assets/settingsDialog.png)

Meeting settings include:

- Schedule Day: Day of the week when the meeting is scheduled.
- Hashtag Format: The format of the hashtag for the meeting date. The date format is based on [Go date and time formatting](https://yourbasic.org/golang/format-parse-string-time-date-example/#standard-time-and-date-formats)

#### Slash Commands to manage the meeting agenda:

```
/agenda queue [next-week] message
```
Creates a post for the user with the given `message` for the next meeting date with the configured hashtag format preceding it.
If `next-week` is indicated (optional), it will use the date of the meeting in the next calendar week. 

![post_example](./assets/postExample.png)

```
/agenda list [next-week]
```
Executes a search of the hashtag of the next meeting, opening the RHS with all the posts with that hashtag. 
If `next-week` is indicated (optional), it will use the date of the meeting in the next calendar week. 

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
- Handle multiple meeting days in a week.
- Handle time in meeting schedule. 
