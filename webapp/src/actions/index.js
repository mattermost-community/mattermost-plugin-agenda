import {searchPostsWithParams} from 'mattermost-redux/actions/search';
import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {Client4} from 'mattermost-redux/client';

import Client from '../client';

import ActionTypes from '../action_types';

export function fetchMeetingSettings(channelId = '') {
    return async (dispatch) => {
        let data;
        try {
            data = await (new Client()).getMeetingSettings(channelId);
        } catch (error) {
            return {error};
        }

        dispatch({
            type: ActionTypes.RECEIVED_MEETING_SETTINGS,
            data,
        });

        return {data};
    };
}

export function saveMeetingSettings(meeting) {
    let data;
    try {
        data = (new Client()).saveMeetingSettings(meeting);
    } catch (error) {
        return {error};
    }

    return {data};
}

export const openMeetingSettingsModal = (channelId = '') => (dispatch) => {
    dispatch({
        type: ActionTypes.OPEN_MEETING_SETTINGS_MODAL,
        channelId,
    });
};

export const closeMeetingSettingsModal = () => (dispatch) => {
    dispatch({
        type: ActionTypes.CLOSE_MEETING_SETTINGS_MODAL,
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
        const teamId = getCurrentTeamId(getState());
        const config = getConfig(getState());
        const viewArchivedChannels = config.ExperimentalViewArchivedChannels === 'true';

        return dispatch(searchPostsWithParams(teamId, {terms, is_or_search: false, include_deleted_channels: viewArchivedChannels, page: 0, per_page: 20}, true));
    };
}

export function requeueItem(postId) {
    return async (dispatch, getState) => {
        const command = `/agenda requeue post ${postId}`;
        await clientExecuteCommand(dispatch, getState, command);
        return {data: true};
    };
}

export async function clientExecuteCommand(dispatch, getState, command) {
    const state = getState();
    const currentChannel = getCurrentChannel(state);
    const currentTeamId = getCurrentTeamId(state);
    const args = {
        channel_id: currentChannel?.id,
        team_id: currentTeamId,
    };

    try {
        return Client4.executeCommand(command, args);
    } catch (error) {
        return error;
    }
}
