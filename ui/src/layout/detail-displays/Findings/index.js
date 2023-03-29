import React from 'react';
import { useLocation } from 'react-router-dom';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Title from 'components/Title';
import LinksList from 'components/LinksList';
import VulnerabilitiesDisplay, { getTotlalVulnerabilitiesFromCounters } from 'components/VulnerabilitiesDisplay';
import { FINDINGS_MAPPING, VULNERABIITY_FINDINGS_ITEM } from 'utils/systemConsts';
import { useFilterDispatch, setFilters, FILTER_TYPES } from 'context/FiltersProvider';
import FindingsCounterDisplay from './FindingsCounterDisplay';
import FindingsSystemFilterLinks from './FindingsSystemFilterLinks';

const Findings = ({findingsSummary={}, findingsFilter, findingsFilterTitle, findingsFilterSuffix=""}) => {
    const {totalVulnerabilities} = findingsSummary;

    const {pathname} = useLocation();
    const filtersDispatch = useFilterDispatch();

    const onFindingsClick = () => {
        setFilters(filtersDispatch, {
            type: FILTER_TYPES.FINDINGS,
            filters: {
                filter: findingsFilter,
                name: findingsFilterTitle,
                suffix: findingsFilterSuffix,
                backPath: pathname,
                customDisplay: () => (
                    <FindingsSystemFilterLinks
                        findingsSummary={findingsSummary}
                        totalVulnerabilitiesCount={getTotlalVulnerabilitiesFromCounters(totalVulnerabilities)}
                    />
                )
            },
            isSystem: true
        });
    }
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <Title medium>Findings</Title>
                    <LinksList
                        items={[
                            {
                                path: VULNERABIITY_FINDINGS_ITEM.appRoute,
                                component: () => <VulnerabilitiesDisplay counters={totalVulnerabilities} />,
                                callback: onFindingsClick
                            },
                            ...Object.keys(FINDINGS_MAPPING).map(findingType => {
                                const {totalKey, title, icon, color, appRoute} = FINDINGS_MAPPING[findingType];
                                
                                return {
                                    path: appRoute,
                                    component: () => (
                                        <FindingsCounterDisplay key={totalKey} icon={icon} color={color} count={findingsSummary[totalKey] || 0} title={title} />
                                    ),
                                    callback: onFindingsClick
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