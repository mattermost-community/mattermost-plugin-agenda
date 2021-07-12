import {AnyAction, Dispatch} from 'redux';
import {searchPostsWithParams} from 'mattermost-redux/actions/search';
import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {GetStateFunc} from 'mattermost-redux/types/actions';
import {Client4} from 'mattermost-redux/client';
import {IntegrationTypes} from 'mattermost-redux/action_types';


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

export function setTriggerId(triggerId) {
    return {
        type: IntegrationTypes.RECEIVED_DIALOG_TRIGGER_ID,
        data: triggerId,
    };
}

export function requeueItem(itemId) {
    return (dispatch, getState) => {
        const command = `/agenda requeue ${itemId}`;
        clientExecuteCommand(dispatch, getState, command).then(r => {console.log({r})});

        return {data: true};
    };
}


export async function clientExecuteCommand(dispatch, getState, command) {
    let currentChannel = getCurrentChannel(getState());
    const currentTeamId = getCurrentTeamId(getState());

    // Default to town square if there is no current channel (i.e., if Mattermost has not yet loaded)
    if (!currentChannel) {
        currentChannel = await Client4.getChannelByName(currentTeamId, 'town-square');
    }

    const args = {
        channel_id: currentChannel?.id,
        team_id: currentTeamId,
    };

    try {
        //@ts-ignore Typing in mattermost-redux is wrong
        const data = await Client4.executeCommand(command, args);
        dispatch(setTriggerId(data?.trigger_id));
    } catch (error) {
        console.error(error); //eslint-disable-line no-console
    }
}