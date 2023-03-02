import React from 'react';
import ProgressBar from 'components/ProgressBar';

import './scan-status-display.scss';

const ScanStatusDisplay = ({itemsCompleted, itemsLeft, errorMessage}) => (
    <div className="scan-status-display-wrapper">
        <ProgressBar itemsCompleted={itemsCompleted} itemsLeft={itemsLeft} />
        {errorMessage &&
            <div className="error-display-wrapper">
                <div className="error-display-title">Some of the elements were failed to be scanned.</div>
                <div className="error-display-message">{errorMessage}</div>
            </div>
        }
    </div>
)

export default ScanStatusDisplay;