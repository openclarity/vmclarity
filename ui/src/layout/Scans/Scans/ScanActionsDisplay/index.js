import React from 'react';
import Icon, { ICON_NAMES } from 'components/Icon';
import { TooltipWrapper } from 'components/Tooltip';

import './scan-actions-display.scss';

const ScanActionsDisplay = ({data}) => {
    const {id} = data;

    return (
        <div className="scan-actions-display">
            <TooltipWrapper tooltipId={`${id}-stop`} tooltipText="Stop scan" >
                <Icon
                    name={ICON_NAMES.STOP}
                    onClick={event => {
                        event.stopPropagation();
                        event.preventDefault();
                        
                        console.log("stop scan");
                    }}
                />
            </TooltipWrapper>
        </div>
    );
}

export default ScanActionsDisplay;