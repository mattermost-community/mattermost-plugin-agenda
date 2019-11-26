import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {channelSettingsModalState} from 'selectors';

import {closeChannelSettingsModal} from 'actions';

import ChannelSettingsModal from './channel_settings';

function mapStateToProps(state) {
    return {
        visible: channelSettingsModalState(state).visible,
        channelId: channelSettingsModalState(state).channelId,
    };
}

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeChannelSettingsModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(ChannelSettingsModal);
