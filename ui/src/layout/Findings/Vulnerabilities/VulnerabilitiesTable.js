import React, { useMemo } from 'react';
import TablePage from 'components/TablePage';
import ExpandableList from 'components/ExpandableList';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';
import { FILTER_TYPES } from 'context/FiltersProvider';
import SeverityWithCvssDisplay from './SeverityWithCvssDisplay';
import { getHigestVersionCvssData } from './utils';

const TABLE_TITLE = "vulnerabilities";

const VulnerabilitiesTable = () => {
    const columns = useMemo(() => [
        {
            Header: "Vulnerability name",
            id: "name",
            accessor: "findingInfo.vulnerabilityName",
            disableSort: true
        },
        {
            Header: "Severity",
            id: "severity",
            Cell: ({row}) => {
                const {id, findingInfo} = row.original;
                const {severity, cvss} = findingInfo || {};
                const cvssScoreData = getHigestVersionCvssData(cvss);
                
                return (
                    <SeverityWithCvssDisplay
                        severity={severity}
                        cvssScore={cvssScoreData.score}
                        cvssSeverity={cvssScoreData.severity.toLocaleUpperCase()}
                        compareTooltipId={`severity-compare-tooltip-${id}`}
                    />
                )
            },
            disableSort: true
        },
        {
            Header: "Package name",
            id: "packageName",
            accessor: "findingInfo.package.name",
            disableSort: true
        },
        {
            Header: "Package version",
            id: "packageVersion",
            accessor: "findingInfo.package.version",
            disableSort: true
        },
        {
            Header: "Fix versions",
            id: "fixVersions",
            Cell: ({row}) => {
                const {versions} = row.original.findingInfo?.fix || {};

                return (
                    <ExpandableList items={versions || []} />
                )
            },
            disableSort: true
        },
        ...getAssetAndScanColumnsConfigList()
    ], []);

    return (
        <TablePage
            columns={columns}
            url={APIS.FINDINGS}
            tableTitle={TABLE_TITLE}
            filterType={FILTER_TYPES.FINDINGS}
            filters="findingInfo/objectType eq 'Vulnerability'"
            expand="asset,scan"
            absoluteSystemBanner
        />
    )
}

export default VulnerabilitiesTable;