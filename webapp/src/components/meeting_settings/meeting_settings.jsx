import React from 'react';
import PropTypes from 'prop-types';

import {Modal} from 'react-bootstrap';

export default class MeetingSettingsModal extends React.PureComponent {
    static propTypes = {
        visible: PropTypes.bool.isRequired,
        channelId: PropTypes.string.isRequired,
        close: PropTypes.func.isRequired,
        meeting: PropTypes.object,
        fetchMeetingSettings: PropTypes.func.isRequired,
        saveMeetingSettings: PropTypes.func.isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            hashtag: null,
            weekday: 0,
        };
    }

    componentDidUpdate(prevProps) {
        if (this.props.channelId && this.props.channelId !== prevProps.channelId) {
            this.props.fetchMeetingSettings(this.props.channelId);
        }

        if (this.props.meeting && this.props.meeting !== prevProps.meeting) {
            // eslint-disable-next-line react/no-did-update-set-state
            this.setState({
                hashtag: this.props.meeting.hashtagFormat,
                weekday: this.props.meeting.schedule,
            });
        }
    }

    handleHashtagChange = (e) => {
        this.setState({
            hashtag: e.target.value,
        });
    }

    handleCheckboxChanged = (e) => {
        this.setState({
            weekday: Number(e.target.value),
        });
    }

    onSave = () => {
        this.props.saveMeetingSettings({
            channelId: this.props.channelId,
            hashtagFormat: this.state.hashtag,
            schedule: this.state.weekday,
        });

        this.props.close();
    }

    getDaysCheckboxes() {
        // weekday is based on the server time.Weekday
        // where the week array starts on Sunday = 0
        const meetingDay = this.state.weekday ? this.state.weekday : 0;

        const weekDays = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

        const checkboxes = weekDays.map((weekday, i) => {
            return (
                <label
                    className='checkbox-inline pl-3 pr-2'
                    key={weekday}
                >
                    <input
                        key={weekday}
                        type='checkbox'
                        value={i + 1}
                        checked={i + 1 === meetingDay}
                        onChange={this.handleCheckboxChanged}
                    /> {weekday}
                </label>);
        });

        return checkboxes;
    }

    render() {
        return (
            <Modal
                dialogClassName='a11y__modal modal-xl'
                onHide={this.props.close}
                show={this.props.visible}
                role='dialog'
                aria-labelledby='agendaPluginMeetingSettingsModalLabel'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title
                        componentClass='h1'
                        id='agendaPluginMeetingSettingsModalLabel'
                    >
                        {'Channel Agenda Settings'}
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <div className='form-group'>
                        <label className='control-label'>
                            {'Meeting Day'}
                        </label>
                        <div className='p-2'>
                            {this.getDaysCheckboxes()}
                        </div>
                    </div>
                    <div className='form-group'>
                        <label className='control-label'>{'Hashtag Format'}</label>
                        <input
                            onInput={this.handleHashtagChange}
                            className='form-control'
                            value={this.state.hashtag ? this.state.hashtag : ''}
                        />
                    </div>
                </Modal.Body>
                <Modal.Footer>
                    <button
                        type='button'
                        className='btn btn-link'
                        onClick={this.props.close}
                    >
                        {'Cancel'}
                    </button>
                    <button
                        onClick={this.onSave}
                        id='save-button'
                        className='btn btn-primary save-button'
                    >
                        {'Save'}
                    </button>
                </Modal.Footer>
            </Modal>
        );
    }
}
