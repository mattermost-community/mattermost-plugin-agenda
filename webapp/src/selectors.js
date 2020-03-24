import {id as pluginId} from './manifest';
import reducer from './reducer';

const defaultState = reducer({}, {});
const getPluginState = (state) => state['plugins-' + pluginId] || defaultState;

export const getMeetingSettingsModalState = (state) => getPluginState(state).meetingSettingsModal;
export const getMeetingSettings = (state) => getPluginState(state).meetingSettings;
