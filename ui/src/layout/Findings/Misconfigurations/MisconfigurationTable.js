import React, { useMemo } from 'react';
import TablePage from 'components/TablePage';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';
import { FILTER_TYPES } from 'context/FiltersProvider';

const TABLE_TITLE = "misconfiguration";

const MisconfigurationsTable = () => {
    const columns = useMemo(() => [
        {
            Header: "Path",
            id: "path",
            accessor: "findingInfo.path",
            disableSort: true
        },
        {
            Header: "Description",
            id: "description",
            accessor: "findingInfo.description",
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
            filters="findingInfo/objectType eq 'Misconfiguration'"
            expand="asset,scan"
            absoluteSystemBanner
        />
    )
}

export default MisconfigurationsTable;