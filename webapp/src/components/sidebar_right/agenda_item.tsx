import React from 'react';

type Props = {
    title: string;
    icon: string;
    username: string;
    cardId: string;
    boardId: string;
    parentId: string;
    teamId: string;
}

const AgendaItem: React.FC<Props> = ({title, username, icon, teamId, boardId, cardId}: Props) => {
    const cardLink = 'http://localhost:8065/boards/team/' + teamId + '/' + boardId + '/vx3t46h8cxfrefkjmxqtkisf58o/' + cardId;

    return (
        <div style={style.item}>
            <div>
                <div style={style.titleContainer}>
                    { icon ? <div style={style.icon}>{icon}</div> : undefined }
                    <a
                        style={style.itemTitle}
                        href={cardLink}
                        target='_blank'
                        rel='noopener noreferrer'

                    > {title}</a>
                </div>
                <div style={style.itemDescrition}>{`Created by: ${username}`}</div>
            </div>
        </div>
    );
};

const style = {
    sectionHeader: {
        padding: '15px',
    },
    titleContainer: {
        display: 'flex',
        flex: '1 1 auto',
    },
    item: {
        display: 'flex',
        border: '1px solid #000000a3',
        boxShadow: 'rgba(0,0,0, 0.1) 0 0 0 1px,rgba(0,0,0, 0.1) 0 2px 4px',
        borderRadius: '4px',
        alignItems: 'center',
        padding: '9px',
        margin: '10px',
    },
    itemTitle: {
        fontWeight: 600,
    },
    itemDescrition: {
        margin: '0px',
    },
    icon: {
        fontSize: '16px',
        marginRight: '4px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: '20px',
        height: '20px',
    },
};
export default AgendaItem;