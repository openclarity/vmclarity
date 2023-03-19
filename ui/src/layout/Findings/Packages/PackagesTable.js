import React, { useMemo } from 'react';
import TablePage from 'components/TablePage';
import ExpandableList from 'components/ExpandableList';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';
import { FILTER_TYPES } from 'context/FiltersProvider';

const TABLE_TITLE = "packages";

const PackagesTable = () => {
    const columns = useMemo(() => [
        {
            Header: "Package name",
            id: "name",
            accessor: "findingInfo.name",
            disableSort: true
        },
        {
            Header: "Version",
            id: "version",
            accessor: "findingInfo.version",
            disableSort: true
        },
        {
            Header: "Languege",
            id: "languege",
            accessor: "findingInfo.language",
            disableSort: true
        },
        {
            Header: "Licenses",
            id: "licenses",
            Cell: ({row}) => {
                const {licenses} = row.original.findingInfo || {};

                return (
                    <ExpandableList items={licenses || []} />
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
            filters="findingInfo/objectType eq 'Package'"
            expand="asset,scan"
            absoluteSystemBanner
        />
    )
}

export default PackagesTable;