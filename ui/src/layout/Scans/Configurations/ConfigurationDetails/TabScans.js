import React from 'react';
import DoublePaneDisplay from 'components/DoublePaneDisplay';

const TabScans = ({data}) => {
    const {id} = data || {};
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                "TBD"
            )}
        />
    )
}

export default TabScans;
