import {searchPostsWithParams} from 'mattermost-redux/actions/search';

import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {OPEN_CHANNEL_SETTINGS_MODAL, CLOSE_CHANNEL_SETTINGS_MODAL} from './actions_types.js';

export const openChannelSettingsModal = (channelId = '') => (dispatch) => {
    dispatch({
        type: OPEN_CHANNEL_SETTINGS_MODAL,
        channelId,
    });
};

export const closeChannelSettingsModal = () => (dispatch) => {
    dispatch({
        type: CLOSE_CHANNEL_SETTINGS_MODAL,
    });
};

// Hackathon hack: "Copying" these actions below directly from webapp

export function updateSearchTerms(terms) {
    return {
        type: 'UPDATE_RHS_SEARCH_TERMS',
        terms,
    };
}

export function updateSearchResultsTerms(terms) {
    return {
        type: 'UPDATE_RHS_SEARCH_RESULTS_TERMS',
        terms,
    };
}

export function updateRhsState(rhsState) {
    return {
        type: 'UPDATE_RHS_STATE',
        state: rhsState,
    };
}

export function performSearch(terms) {
    return (dispatch, getState) => {
        console.log('CALLING SEARCH'); // eslint-disable-line no-console

        const teamId = getCurrentTeamId(getState());
        const config = getConfig(getState());
        const viewArchivedChannels = config.ExperimentalViewArchivedChannels === 'true';

        return dispatch(searchPostsWithParams(teamId, {terms, is_or_search: false, include_deleted_channels: viewArchivedChannels, page: 0, per_page: 20}, true));
    };
}