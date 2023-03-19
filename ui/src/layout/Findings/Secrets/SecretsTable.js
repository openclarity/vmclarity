import React, { useMemo } from 'react';
import TablePage from 'components/TablePage';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';
import { FILTER_TYPES } from 'context/FiltersProvider';

const TABLE_TITLE = "secrets";

const SecretsTable = () => {
    const columns = useMemo(() => [
        {
            Header: "Fingerprint",
            id: "fingerprint",
            accessor: "findingInfo.fingerprint",
            disableSort: true
        },
        {
            Header: "Description",
            id: "description",
            accessor: "findingInfo.description",
            disableSort: true
        },
        {
            Header: "FilePath",
            id: "findingInfo",
            accessor: "findingInfo.filePath",
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
            filters="findingInfo/objectType eq 'Secret'"
            expand="asset,scan"
            absoluteSystemBanner
        />
    )
}

export default SecretsTable;