import React from 'react';
import DoublePaneDisplay from 'components/DoublePaneDisplay';

const TabScans = ({data}) => {
    const {id} = data || {};
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                "test"
            )}
            rightPlaneDisplay={() => null}
        />
    )
}

export default TabScans;
