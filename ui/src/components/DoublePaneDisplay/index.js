import React from 'react';

import './double-pane-display.scss';

const DoublePaneDisplay = ({rightPlaneDisplay: RightPlaneDisplay, leftPaneDisplay: LeftPaneDisplay}) => (
    <div className="double-pane-display-wrapper">
        <div className="left-pane-display"><LeftPaneDisplay /></div>
        <div className="right-pane-display"><RightPlaneDisplay /></div>
    </div>
)

export default DoublePaneDisplay;