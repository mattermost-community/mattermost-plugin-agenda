import request from 'superagent';

import {id as pluginId} from './manifest';

export async function getChannelSettings (channelId) {
    const url = `/plugins/${pluginId}/api/v1`;
    
    return doGet(`${url}/settings?channelId=${channelId}`);
}


async function doGet (url, body, headers = {}) {
    headers['X-Requested-With'] = 'XMLHttpRequest';
    headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

    const response = await request.
        get(url).
        set(headers).
        accept('application/json');

    return response.body;
}

async function doPost (url, body, headers = {}) {
    headers['X-Requested-With'] = 'XMLHttpRequest';
    headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

    const response = await request.
        post(url).
        send(body).
        set(headers).
        type('application/json').
    accept('application/json');

    return response.body;
}