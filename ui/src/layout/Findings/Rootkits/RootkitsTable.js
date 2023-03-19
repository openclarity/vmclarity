import React, { useMemo } from 'react';
import TablePage from 'components/TablePage';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';
import { FILTER_TYPES } from 'context/FiltersProvider';

const TABLE_TITLE = "rootkits";

const RootkitsTable = () => {
    const columns = useMemo(() => [
        {
            Header: "Rootkit name",
            id: "rootkitName",
            accessor: "findingInfo.rootkitName",
            disableSort: true
        },
        {
            Header: "path",
            id: "path",
            accessor: "findingInfo.path",
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
            filters="findingInfo/objectType eq 'Rootkit'"
            expand="asset,scan"
            absoluteSystemBanner
        />
    )
}

export default RootkitsTable;