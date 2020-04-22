import request from 'superagent';
import {Client4} from 'mattermost-redux/client';

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
        const response = await request.
            get(url).
            set({...this.getCommonHeaders(), ...headers}).
            accept('application/json');

        return response.body;
    }

    doPost = async (url, body, headers = {}) => {
        const response = await request.
            post(url).
            send(body).
            set({...this.getCommonHeaders(), ...headers}).
            type('application/json').
            accept('application/json');

        return response.body;
    }

    getCommonHeaders = () => {
        const { headers } = Client4.getOptions([]);

        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        return headers;
    }
}
