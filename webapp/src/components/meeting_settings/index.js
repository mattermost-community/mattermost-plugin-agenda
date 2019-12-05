import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getMettingSettingsModalState, getMeetingSettings} from 'selectors';
import {closeMeetingSettingsModal, fetchMeetingSettings, saveMeetingSettings} from 'actions';

import MeetingSettingsModal from './meeting_settings';

function mapStateToProps(state) {
    return {
        visible: getMettingSettingsModalState(state).visible,
        channelId: getMettingSettingsModalState(state).channelId,
        meeting: getMeetingSettings(state).meeting,
        saveMeetingSettings,
    };
}

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeMeetingSettingsModal,
    fetchMeetingSettings,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(MeetingSettingsModal);
