// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import {getUser} from 'mattermost-redux/selectors/entities/users';

import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';

import {fetchQueuedItems} from '../../actions';
import {id as pluginId} from '../../manifest';

import {getQueuedItems} from 'src/selectors';

import SidebarRight from './sidebar_right.jsx';

function mapStateToProps(state) {
    return {
        username: state[`plugins-${pluginId}`].username,
        channelId: getCurrentChannel(state).id,
        teamId: getCurrentTeam(state).id,
        queuedItems: getQueuedItems(state).items,
        getUser,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            fetchQueuedItems,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarRight);
