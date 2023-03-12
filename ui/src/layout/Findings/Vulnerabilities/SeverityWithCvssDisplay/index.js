import React from 'react';
import SeverityDisplay from 'components/SeverityDisplay';

import './severity-with-cvss-display.scss';

const SeverityWithCvssDisplay = ({severity, cvssScore}) => {
    return (
        <div className="severity-with-cvss-display">
            <SeverityDisplay severity={severity} />
            <div className="cvss-score-display">{`(cvss ${cvssScore})`}</div>
        </div>
    );
}

export default SeverityWithCvssDisplay;