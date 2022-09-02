// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import Scrollbars from 'react-custom-scrollbars';

import AgendaItem from './agenda_item';

export function renderView(props) {
    return (
        <div
            {...props}
            className='scrollbar--view'
        />);
}

export function renderThumbHorizontal(props) {
    return (
        <div
            {...props}
            className='scrollbar--horizontal'
        />);
}

export function renderThumbVertical(props) {
    return (
        <div
            {...props}
            className='scrollbar--vertical'
        />);
}

export default class SidebarRight extends React.PureComponent {
    static propTypes = {
        username: PropTypes.string,
        queuedItems: PropTypes.array,
        channelId: PropTypes.string,
        teamId: PropTypes.string,
        actions: PropTypes.shape({
            fetchQueuedItems: PropTypes.func.isRequired,
        }).isRequired,
    };
å
componentDidMount() {
    this.props.actions.fetchQueuedItems(this.props.channelId);
}

getUserName() {
    return 'maria.nunez';
}

render() {
    return (
        <React.Fragment>
            <Scrollbars
                autoHide={true}
                autoHideTimeout={500}
                autoHideDuration={500}
                renderThumbHorizontal={renderThumbHorizontal}
                renderThumbVertical={renderThumbVertical}
                renderView={renderView}
            >
                <div style={style.sectionHeader}>
                    <strong>
                        <a

                            //href={listUrl}
                            target='_blank'
                            rel='noopener noreferrer'
                        >{'Up Next Agenda Items'}</a>
                    </strong>
                </div>
                <div>
                    {this.props.queuedItems && this.props.queuedItems.length > 0 ? this.props.queuedItems.map((item) => {
                        return (<AgendaItem
                            title={item.title}
                            username={this.getUserName()}
                            icon={item.fields.icon}
                            teamId={this.props.teamId}
                            boardId={item.boardId}
                            cardId={item.id}
                        />);
                    }) :
                        <p> {'No items queued'}</p>
                    }
                </div>
            </Scrollbars>
        </React.Fragment>
    );
}
}

const style = {
    sectionHeader: {
        padding: '15px',
    },
};
