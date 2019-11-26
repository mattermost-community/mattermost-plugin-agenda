import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

import {getChannelSetting} from '../../client'

import {Modal} from 'react-bootstrap';

export default class ChannelSettingsModal extends React.PureComponent {
    static propTypes = {
        visible: PropTypes.bool.isRequired,
        channelId: PropTypes.string.isRequired,
        close: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
        subMenu: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
    };

    render() {
        return (
            <Modal
                dialogClassName='a11y__modal'
                onHide={this.props.close}
                show={this.props.visible}
                role='dialog'
                aria-labelledby='agendaPluginChannelSettingsModalLabel'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title
                        componentClass='h1'
                        id='renameChannelModalLabel'
                    >
                        <FormattedMessage
                            id='rename_channel.title'
                            defaultMessage='Agenda Plugin Channel Settings'
                        />
                    </Modal.Title>
                </Modal.Header>
                <form role='form'>
                    <Modal.Body>
                        <div className='form-group'>
                            <label className='control-label'>
                                {
                                    'Meeting Schedule Days'
                                }
                            </label>
                            <div className='form-control'>
                                <label className='checkbox-inline'>
                                    <input
                                        type='checkbox'
                                        value=''
                                    /> {'Monday'}
                                </label>
                                <label className='checkbox-inline'>
                                    <input
                                        type='checkbox'
                                        value=''
                                    /> {'Tuesday'}
                                </label>
                                <label className='checkbox-inline'>
                                    <input
                                        type='checkbox'
                                        value=''
                                        checked
                                    />{'Wednesday'}
                                </label>
                                <label className='checkbox-inline'>
                                    <input
                                        type='checkbox'
                                        value=''
                                    />{'Thrusday'}
                                </label>
                                <label className='checkbox-inline'>
                                    <input
                                        type='checkbox'
                                        value=''
                                    />{'Friday'}
                                </label>
                            </div>
                        </div>
                        <div className='form-group'>
                            <label className='control-label'>{'Hashtag Format'}</label>
                            <input className='form-control' value="DEV-Jan02"/>
                        </div>
                    </Modal.Body>
                    <Modal.Footer>
                        <button
                            type='button'
                            className='btn btn-link'
                            onClick={this.props.close}
                        >
                            <FormattedMessage
                                id='rename_channel.cancel'
                                defaultMessage='Cancel'
                            />
                        </button>
                        <button
                            onClick={this.props.close}
                            type='submit'
                            id='save-button'
                            className='btn btn-primary'
                        >
                            <FormattedMessage
                                id='rename_channel.save'
                                defaultMessage='Save'
                            />
                        </button>
                    </Modal.Footer>
                </form>
            </Modal>
        );
    }
}