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

    doGet = async (url, headers = {}) => {
        return this.doFetch(url, { headers });
    }

    doPost = async (url, body, headers = {}) => {
        return this.doFetch(url, {
            method : 'POST',
            body: JSON.stringify(body),
            headers: {...headers, ...{
                'Content-Type' : 'application/json'
            }}
        });
    }

    doFetch = async (url, { method = 'GET', body = null, headers = {} }) => {
        const response = await fetch(url, {
            method,
            body,
            headers: {...this.getCommonHeaders(), ...headers}
        });

        return response.json();
    }

    getCommonHeaders = () => {
        const { headers } = Client4.getOptions([]);

        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();
        headers['Accept'] = 'application/json';

        return headers;
    }
}
