import React from 'react';
import { useNavigate } from 'react-router-dom';
import Title from 'components/Title';
import Icon, { ICON_NAMES } from 'components/Icon';
import { TooltipWrapper } from 'components/Tooltip';
import { SCAN_CONFIGS_PATH } from 'layout/Scans/Configurations';
import { ROUTES } from 'utils/systemConsts';

import './configuration-alert-link.scss';

const CONFIGURATION_ALERT_TEXT = (
    <span>
        Configuration has been modified since<br />
        the scan has performed and it might not<br />
        match the scan's configuration<br />
    </span>
)

const ConfigurationAlertLink = ({configData}) => {
    const navigate = useNavigate();

    return (
        <div className="configuration-alert-link">
            <Title medium removeMargin onClick={() => navigate(`/${ROUTES.SCANS}/${SCAN_CONFIGS_PATH}/${configData.id}`)}>Configuration</Title>
            <TooltipWrapper tooltipId="configuration-alert-tooltip" tooltipText={CONFIGURATION_ALERT_TEXT}>
                <Icon name={ICON_NAMES.WARNING} />
            </TooltipWrapper>
        </div>
    )
}

export default ConfigurationAlertLink;