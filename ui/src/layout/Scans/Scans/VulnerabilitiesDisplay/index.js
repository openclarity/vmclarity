import React from 'react';
import { TooltipWrapper } from 'components/Tooltip';
import Icon, { ICON_NAMES } from 'components/Icon';
import { toCapitalized } from 'utils/utils';

import COLORS from 'utils/scss_variables.module.scss';

import './vulnerabilities-display.scss';

const SEVERITY_ITEMS = {
    critical: {value: "critical", color: COLORS["color-error-dark"], innerColor: "white"},
    high: {value: "high", color: COLORS["color-error"], innerColor: "white"},
    medium: {value: "medium", color: COLORS["color-warning"], innerColor: COLORS["color-grey-black"]},
    low: {value: "low", color: COLORS["color-warning-low"], innerColor: COLORS["color-grey-black"]},
    none: {value: "none", color: COLORS["color-status-blue"], innerColor: COLORS["color-grey-black"]}
}

const TooltipContentDisplay = ({total, counters}) => (
    <div className="vulnerabilities-minimized-tooltip-content">
        <div>{`Vulnerabilities: ${total}`}</div>
        <div className="vulnerabilities-tooltip-counters">
            {
                Object.values(SEVERITY_ITEMS).map(({value, color}) => (
                    <div key={value} className="vulnerabilities-tooltip-counters-item"><Icon name={ICON_NAMES.SHIELD} style={{color}} /><span>{counters[value]}</span></div>
                ))
            }
        </div>
    </div>
)

const MinimizedVulnerabilitiesDisplay = ({id, highestSeverity, total, counters}) => {
    const {color, innerColor} = SEVERITY_ITEMS[highestSeverity];

    return (
        <div className="vulnerabilities-minimized-display-wrapper">
            <TooltipWrapper tooltipId={`vulnerability-minimized-display-${id}`} tooltipText={<TooltipContentDisplay total={total} counters={counters} />}>
                <div className="vulnerabilities-minimized-display-summary-item" style={{color: innerColor, backgroundColor: color}}>{counters[highestSeverity]}</div>
            </TooltipWrapper>
        </div>
    )
}

const CounterItemDisplay = ({count, title, color}) => (
    <div className="vulnerabilities-display-counter-item">
        <div className="vulnerabilities-counter-item-count" style={{color}}>{count}</div>
        <div className="vulnerabilities-counter-item-title">{title}</div>
    </div>
)

const VulnerabilitiesDisplay = ({highestSeverity, total, counters}) => {
    const {color} = SEVERITY_ITEMS[highestSeverity];

    return (
        <div className="vulnerabilities-display-wrapper">
            <div className="vulnerabilities-display-summary">
                <Icon name={ICON_NAMES.SHIELD} style={{color}} />
                <CounterItemDisplay count={total} title="Vulnerabilities" color={COLORS["color-main"]} />
            </div>
            <div className="vulnerabilities-display-counters">
                {
                    Object.values(SEVERITY_ITEMS).map(({value, color}) => (
                        <CounterItemDisplay key={value} count={counters[value]} title={toCapitalized(value)} color={color} />
                    ))
                }
            </div>
        </div>
    )
}


const VulnerabilitiesDisplayWrapper = ({id, isMinimized}) => {
    const counters = {critical: 0, high: 44, medium: 8, low: 4, none: 0};

    const total = Object.values(SEVERITY_ITEMS).reduce((acc, curr) => {
        return acc + counters[curr.value];
    }, 0);

    const highestSeverity = (Object.values(SEVERITY_ITEMS).find(({value}) => counters[value] > 0) || SEVERITY_ITEMS.none).value;

    const DisplayComponent = isMinimized ? MinimizedVulnerabilitiesDisplay : VulnerabilitiesDisplay;
    
    return (
        <DisplayComponent id={id} highestSeverity={highestSeverity} total={total} counters={counters} />
    )
}

export default VulnerabilitiesDisplayWrapper;