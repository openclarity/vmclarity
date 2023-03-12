import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import { APIS } from 'utils/systemConsts';
import { getVulnerabilitiesColumnConfigItem } from 'utils/utils';

const TABLE_TITLE = "packages";

const PackagesTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();
    
    const columns = useMemo(() => [
        {
            Header: "Package name",
            id: "name",
            accessor: "id",
            disableSort: true
        },
        {
            Header: "Version",
            id: "version",
            accessor: "version",
            disableSort: true
        },
        {
            Header: "Languege",
            id: "languege",
            accessor: "languege",
            disableSort: true
        },
        {
            Header: "License",
            id: "license",
            accessor: "license",
            disableSort: true
        },
        getVulnerabilitiesColumnConfigItem({tableTitle: TABLE_TITLE}),
        {
            Header: "Asset name",
            id: "assetName",
            accessor: "assetName",
            disableSort: true
        },
        {
            Header: "Asset location",
            id: "assetLocation",
            accessor: "assetLocation",
            disableSort: true
        },
        {
            Header: "Scan",
            id: "scan",
            accessor: "scan",
            disableSort: true
        }
    ], []);

    return (
        <ContentContainer>
            <Table
                columns={columns}
                paginationItemsName={TABLE_TITLE.toLowerCase()}
                url={APIS.SCANS}
                noResultsTitle={TABLE_TITLE}
                onLineClick={({id}) => navigate(`${pathname}/${id}`)}
            />
        </ContentContainer>
    )
}

export default PackagesTable;