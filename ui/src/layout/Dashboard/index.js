import React from 'react';
import { APIS } from 'utils/systemConsts';
import CounterDisplay from './CounterDisplay';
import WidgetWrapper from './WidgetWrapper';

import COLORS from 'utils/scss_variables.module.scss';

import './dashboard.scss';

const COUNTERS_CONFIG = [
    {url: APIS.SCANS, title: "Completed scans", background: COLORS["color-gradient-green"]},
    {url: APIS.ASSETS, title: "Scanned assets", background: COLORS["color-gradient-blue"]},
    {url: APIS.FINDINGS, title: "Total risky findings", background: COLORS["color-gradient-yellow"]}
];

const Dashboard = () => {
    return (
        <div className="dashboard-page-wrapper">
            {
                COUNTERS_CONFIG.map(({url, title, background}, index) => (
                    <CounterDisplay key={index} url={url} title={title} background={background} />
                ))
            }
            <WidgetWrapper className="riskiest-regions" title="Riskiest regions">
                TBD
            </WidgetWrapper>
            <WidgetWrapper className="findings-trend" title="Findings trend">
                TBD
            </WidgetWrapper>
            <WidgetWrapper className="riskiest-assets" title="Riskiest assets">
                TBD
            </WidgetWrapper>
            <WidgetWrapper className="findings-impact" title="Findings-impact">
                TBD
            </WidgetWrapper>
        </div>
    )
}

export default Dashboard;