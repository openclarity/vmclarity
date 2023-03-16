import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import ExpandableList from 'components/ExpandableList';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';

const TABLE_TITLE = "packages";

const PackagesTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();
    
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
        <ContentContainer>
            <Table
                columns={columns}
                paginationItemsName={TABLE_TITLE.toLowerCase()}
                url={APIS.FINDINGS}
                filters={{"$filter": `findingInfo/objectType eq 'Package'`, "$expand": "asset,scan"}}
                noResultsTitle={TABLE_TITLE}
                onLineClick={({id}) => navigate(`${pathname}/${id}`)}
            />
        </ContentContainer>
    )
}

export default PackagesTable;