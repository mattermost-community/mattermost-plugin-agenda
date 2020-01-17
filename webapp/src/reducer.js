import {combineReducers} from 'redux';

import ActionTypes from './action_types';

const meetingSettingsModal = (state = false, action) => {
    switch (action.type) {
    case ActionTypes.OPEN_MEETING_SETTINGS_MODAL:
        return {
            ...state,
            visible: true,
            channelId: action.channelId,
        };
    case ActionTypes.CLOSE_MEETING_SETTINGS_MODAL:
        return {
            ...state,
            visible: false,
            channelId: '',
        };
    default:
        return state;
    }
};

function meetingSettings(state = {}, action) {
    switch (action.type) {
    case ActionTypes.RECEIVED_MEETING_SETTINGS: {
        return {
            ...state,
            meeting: action.data,
        };
    }
    default:
        return state;
    }
}

export default combineReducers({
    meetingSettingsModal,
    meetingSettings,
});

