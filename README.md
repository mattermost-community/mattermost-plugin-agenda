# Agenda Plugin

The Agenda Plugin helps manage a channel's meeting agenda by queuing and listing items based on the meeting's date hashtags. 

The plugin will create posts preceding the agenda item with configured hashtag format and later will open a search with that hashtag to view the agenda list. 

The meeting settings for the channel can be configured in the Channel Header Dropwdown (supported in [this WebApp branch](https://github.com/mattermost/mattermost-webapp/tree/MM-19902))

![channel_header_menu](./assets/channelHeaderDropdown.png)

![settings_dialog](./assets/settingsDialog.png)

Meeting settings include:

- Schedule Day: Day of the week when the meeting is scheduled.
- Hashtag Format: The format of the hashtag for the meeting date.

Created as part of [Mattermost Hackathon 2019](https://github.com/mattermost/mattermost-hackathon-nov2019#how-do-i-submit-my-project)

## Usage

Slash Commands to manage the meeting agenda:

```
/agenda queue [next-week] message
```
Creates a post with the given `message` for the next meeting with the configured hashtag format preceding it.
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
Will update the given setting with the provided value for the meeting settings of that channel. 

`Field` can be one of:

- `schedule`: Day of the week of the meeting. It is an int based on [`time.Weekday`](https://golang.org/pkg/time/#Weekday)
- `hashtag`: Format of the hashtag for the meeting date. It is based on the format used in [`time.Format`](https://golang.org/pkg/time/#Time.Format)


## ToDo

- Complete HTTP hooks to Post on /settings in order to save the setting from the UI.
- Complete saving settings in the WebApp dialog.

## Improvements

- Handle multiple meeting days in a week.
- Add numbering to agenda list. 
- Handle time in meeting schedule. 
- Mark items as resolved or queue for next week. 
- Queue a post using a menu option in the post dot menu. 
