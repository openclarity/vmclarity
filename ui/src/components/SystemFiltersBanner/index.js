import React from 'react';
import { useNavigate } from 'react-router-dom';
import classnames from 'classnames';
import CloseButton from 'components/CloseButton';
import Arrow, { ARROW_NAMES } from 'components/Arrow';

import './system-filter-banner.scss';

const SystemFilterBanner = ({onClose, displayText, backPath, absolute=false, customDisplay: CustomDisplay}) => {
    const navigate = useNavigate();

    return (
        <div className={classnames("system-filter-banner", {absolute})}>
            <div className="system-filter-content">
                <Arrow name={ARROW_NAMES.LEFT} small onClick={() => navigate(backPath)} />
                <div className="filter-content">{displayText}</div>
            </div>
            <div style={{display: "flex", alignItems: "center"}}>
                {!!CustomDisplay && <CustomDisplay />}
                <CloseButton small onClose={onClose} />
            </div>
        </div>
    )
}

export default SystemFilterBanner;