import request from 'superagent';

import {id as pluginId} from './manifest';

export default class Client {
    constructor() {
        this.url = `/plugins/${pluginId}/api/v1`;
    }

    getMeetingSettings = async (channelId) => {
        return this.doGet(`${this.url}/settings?channelId=${channelId}`);
    }

    saveMeetingSettings = async (meeting) => {
        return this.doPost(`${this.url}/settings`, meeting);
    }

    doGet = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        const response = await request.
            get(url).
            set(headers).
            accept('application/json');

        return response.body;
    }

    doPost = async (url, body, headers = {}) => {
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
}
