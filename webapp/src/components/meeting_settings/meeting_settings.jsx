import React from 'react';
import PropTypes from 'prop-types';

import {Modal} from 'react-bootstrap';

import Select from 'react-select';

export default class MeetingSettingsModal extends React.PureComponent {
    static propTypes = {
        visible: PropTypes.bool.isRequired,
        channelId: PropTypes.string.isRequired,
        close: PropTypes.func.isRequired,
        meeting: PropTypes.object,
        fetchMeetingSettings: PropTypes.func.isRequired,
        saveMeetingSettings: PropTypes.func.isRequired,
    };

    options = [
        {value: 'Jan 2', label: 'month_day'},
        {value: '2 Jan', label: 'day_month'},
        {value: '1 2', label: 'month_day'},
        {value: '2 1', label: 'day_month'},
        {value: '2006 1 2', label: 'year_month_day'},
    ];

    customStyles = {
        menuList: (provided) => {
            return ({
                ...provided,
                height: 188,
            });
        },
        control: (provided, state) => ({
            ...provided,
            height: 34,
            minHeight: 34,
            border: '1px solid #ced4da',
            boxShadow: state.isFocused ? 0 : '1px solid #ced4da',
            '&:hover': {
                border: '1px solid #ced4da',
            },
        }),
        indicatorsContainer: (provided) => {
            return ({
                ...provided,
                height: 34,
            });
        },
        singleValue: (provided, state) => {
            const opacity = state.isDisabled ? 0.5 : 1;
            const transition = 'opacity 300ms';

            return {...provided, opacity, transition};
        },
    }

    constructor(props) {
        super(props);

        this.state = {
            hashtagPrefix: 'Prefix',
            weekdays: [1],
            dateFormat: '1-2', // dateFormat will be an object type => { value: string, label: string }
        };
    }

    componentDidUpdate(prevProps) {
        if (this.props.channelId && this.props.channelId !== prevProps.channelId) {
            this.props.fetchMeetingSettings(this.props.channelId);
        }

        if (this.props.meeting && this.props.meeting !== prevProps.meeting) {
            const splitResult = this.props.meeting.hashtagFormat.split('{{');// we know, date Format is preceded by {{
            const hashtagPrefix = splitResult[0];
            const dateFormatValue = splitResult[1].substring(0, splitResult[1].length - 2).trim(); // remove trailing }}
            const dateFormat = this.options.filter((i) => i.value === dateFormatValue)[0]; // extract value object
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

    handleDateFormat = (newValue) => {
        this.setState({dateFormat: newValue});
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
            hashtagFormat: `${this.state.hashtagPrefix}{{${this.state.dateFormat.value}}}`,
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
                <Modal.Body style={{overflow: 'visible'}}>
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
                                style={{padding: '5px', minWidth: '175px'}}
                            >
                                <label className='control-label'>{'Date Format'}</label>
                                <br/>
                                <Select
                                    name='format'
                                    className='form-select'
                                    styles={this.customStyles}
                                    isSearchable={false}
                                    value={this.state.dateFormat}
                                    options={this.options}
                                    onChange={this.handleDateFormat.bind(this)}
                                />
                            </div>
                        </div>

                        <p className='text-muted pt-1'>
                            <div
                                className='alert alert-warning'
                                role='alert'
                                style={{marginBottom: '3px'}}
                            >
                                {'Prefixes may use underscore'}<code>{'_'}</code>{'.'} {'Other special characters including'} <code>{'-'}</code> {'are not allowed.'}
                            </div>
                            {'Date would be appended to Hashtag Prefix, according to the chosen format.'}
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
