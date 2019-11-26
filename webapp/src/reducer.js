import {combineReducers} from 'redux';

import {OPEN_CHANNEL_SETTINGS_MODAL, CLOSE_CHANNEL_SETTINGS_MODAL} from './actions/actions_types';

const channelSettingsModal = (state = false, action) => {
    switch (action.type) {
    case OPEN_CHANNEL_SETTINGS_MODAL:
        return {
            ...state,
            visible: true,
            channelId: action.channelId,
        };
    case CLOSE_CHANNEL_SETTINGS_MODAL:
        return {
            ...state,
            visible: false,
            channelId: '',
        };
    default:
        return state;
    }
};

export default combineReducers({
    channelSettingsModal,
});

