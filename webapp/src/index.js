
import {updateSearchTerms, updateSearchResultsTerms, updateRhsState, performSearch, openChannelSettingsModal} from './actions';

import reducer from './reducer';

import ChannelSettingsModal from './components/channel_settings';

import {id as pluginId} from './manifest';
export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerWebSocketEventHandler(
            'custom_' + pluginId + '_list',
            handleSearchHashtag(store)
        );

        registry.registerRootComponent(ChannelSettingsModal);
        registry.registerChannelHeaderMenuAction('Agenda Plugin Settings',
            (channelId) => {
                store.dispatch(openChannelSettingsModal(channelId));
            });

        registry.registerReducer(reducer);
    }
}

function handleSearchHashtag(store) {
    return (msg) => {
        if (!msg.data) {
            return;
        }
        store.dispatch(updateSearchTerms(msg.data.hashtag));
        store.dispatch(updateSearchResultsTerms(msg.data.hashtag));

        store.dispatch(updateRhsState('search'));
        store.dispatch(performSearch(msg.data.hashtag));
    };
}

window.registerPlugin(pluginId, new Plugin());
