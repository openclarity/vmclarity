import React from 'react';
import TitleValueDisplay, { TitleValueDisplayRow, ValuesListDisplay } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import { FindingsDetailsCommonFields } from '../utils';
import { formatNumber } from 'utils/utils';

const TabPackageDetails = ({data}) => {
    const {pathname} = useLocation();
    const filtersDispatch = useFilterDispatch();

    const {id, findingInfo, firstSeen, lastSeen, summary} = data;
    const {totalVulnerabilities} = summary || {};
    const {name, version, language, licenses} = findingInfo;

    const onVulnerabilitiesClick = () => {
        setFilters(filtersDispatch, {
            type: FILTER_TYPES.FINDINGS_VULNERABILITIES,
            filters: {
                filter: `findingInfo/package/name eq '${name}' and findingInfo/package/version eq '${version}'`,
                name: `Package ${id}`,
                suffix: "finding",
                backPath: pathname
            },
            isSystem: true
        });
    }

    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Package name">{name}</TitleValueDisplay>
                        <TitleValueDisplay title="Version">{version}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Language">{language}</TitleValueDisplay>
                        <TitleValueDisplay title="Licenses"><ValuesListDisplay values={licenses} /></TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <FindingsDetailsCommonFields firstSeen={firstSeen} lastSeen={lastSeen} />
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Asset count">{formatNumber(assetCount)}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                </>
            )}
            rightPlaneDisplay={() => (
                <>
                    <Title medium>Package Vulnerabilities</Title>
                    <LinksList
                        items={[
                            {
                                path: ROUTES.FINDINGS,
                                component: () => <VulnerabilitiesDisplay counters={totalVulnerabilities} />,
                                callback: onVulnerabilitiesClick
                            }
                        ]}
                    />
                </>
            )}
        />
    )
}

export default TabPackageDetails;
