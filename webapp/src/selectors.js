import {id as pluginId} from './manifest';

import {getChannelSettings} from './client'

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const channelSettingsModalState = (state) => getPluginState(state).channelSettingsModal;

export function getChannelSetting (channelId) {
    return getChannelSetting(channelId);
}
