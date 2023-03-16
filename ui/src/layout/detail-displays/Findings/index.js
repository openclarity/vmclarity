import React from 'react';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Title from 'components/Title';
import LinksList from 'components/LinksList';
import VulnerabilitiesDisplay from 'components/VulnerabilitiesDisplay';
import { ROUTES, FINDINGS_MAPPING } from 'utils/systemConsts';
import { FINDINGS_PATHS } from 'layout/Findings'
import FindingsCounterDisplay from './FindingsCounterDisplay';

const Findings = ({findingsSummary}) => {
    const {totalVulnerabilities} = findingsSummary || {};
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <Title medium>Findings</Title>
                    <LinksList
                        items={[
                            {path: `${ROUTES.FINDINGS}/${FINDINGS_PATHS.VULNERABILITIES}`, component: () => <VulnerabilitiesDisplay counters={totalVulnerabilities} />},
                            ...Object.keys(FINDINGS_MAPPING).map(findingType => {
                                const {dataKey, title, icon, color, appRoute} = FINDINGS_MAPPING[findingType];
                                
                                return {
                                    path: appRoute,
                                    component: () => (
                                        <FindingsCounterDisplay key={dataKey} icon={icon} color={color} count={findingsSummary[dataKey] || 0} title={title} />
                                    )
                                }
                            })
                        ]}
                    />
                </>
            )}
        />
    )
}

export default Findings;