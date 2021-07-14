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
            hashtagPrefix: 'Prefix',
            weekdays: [1],
            dateFormat: '1-2',
        };
    }

    componentDidUpdate(prevProps) {
        if (this.props.channelId && this.props.channelId !== prevProps.channelId) {
            this.props.fetchMeetingSettings(this.props.channelId);
        }

        if (this.props.meeting && this.props.meeting !== prevProps.meeting) {
            const splitResult = this.props.meeting.hashtagFormat.split('{{');// we know, date Format is preceded by {{
            const hashtagPrefix = splitResult[0];
            const dateFormat = splitResult[1].substring(0, splitResult[1].length - 2); // remove trailing }}
            // eslint-disable-next-line react/no-did-update-set-state
            this.setState({
                hashtagPrefix,
                dateFormat,
                weekdays: this.props.meeting.schedule || [],
            });
        }
    }

    handleHashtagChange = (e) => {
        this.setState({
            hashtagPrefix: e.target.value,
        });
    }

    handleDateFormat = (event) => {
        this.setState({
            dateFormat: event.target.value,
        });
    };

    handleCheckboxChanged = (e) => {
        const changeday = Number(e.target.value);
        let changedWeekdays = Object.assign([], this.state.weekdays);

        if (e.target.checked && !this.state.weekdays.includes(changeday)) {
            // Add the checked day
            changedWeekdays = [...changedWeekdays, changeday];
        } else if (!e.target.checked && this.state.weekdays.includes(changeday)) {
            // Remove the unchecked day
            changedWeekdays.splice(changedWeekdays.indexOf(changeday), 1);
        }

        this.setState({
            weekdays: changedWeekdays,
        });
    }

    onSave = () => {
        this.props.saveMeetingSettings({
            channelId: this.props.channelId,
            hashtagFormat: `${this.state.hashtagPrefix}{{${this.state.dateFormat}}}`,
            schedule: this.state.weekdays.sort(),
        });

        this.props.close();
    }

    getDaysCheckboxes() {
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
                        value={i}
                        checked={this.state.weekdays.includes(i)}
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
                        <div style={{display: 'flex'}}>
                            <div
                                className='fifty'
                                style={{padding: '5px'}}
                            >
                                <label className='control-label'>{'Hashtag Prefix'}</label>
                                <input
                                    onChange={this.handleHashtagChange}
                                    className='form-control'
                                    value={this.state.hashtagPrefix ? this.state.hashtagPrefix : ''}
                                />
                            </div>
                            <div
                                className='fifty'
                                style={{padding: '5px'}}
                            >
                                <label className='control-label'>{'Date Format'}</label>
                                <br/>
                                <select
                                    name='format'
                                    value={this.state.dateFormat}
                                    onChange={this.handleDateFormat}
                                    style={{height: '35px', border: '1px solid #ced4da'}}
                                    className='form-select'
                                >
                                    <option value='Jan 2'>{'Month_day'}</option>
                                    <option value='2 Jan'>{'day_Month'}</option>
                                    <option value='1 2'>{'month_day'}</option>
                                    <option value='2 1'>{'day_month'}</option>
                                    <option value='2006 1 2'>{'year_month_day'}</option>

                                </select>
                            </div>
                        </div>

                        <p className='text-muted pt-1'>
                            <div
                                className='alert alert-warning'
                                role='alert'
                                style={{marginBottom: '3px'}}
                            >
                                {'You may use underscore'}<code>{'_'}</code>{'.'} {'Other special characters including'} <code>{'-'}</code>{','} {'not allowed.'}
                            </div>
                            {'Date would be appended to Hashtag Prefix, according to format chosen.'}
                        </p>
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
