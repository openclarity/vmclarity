import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import { APIS } from 'utils/systemConsts';

const TABLE_TITLE = "vulnerabilities";

const VulnerabilitiesTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();

    const columns = useMemo(() => [
        {
            Header: "Vulnerability name",
            id: "name",
            accessor: "id",
            disableSort: true
        },
        {
            Header: "Severity",
            id: "severity",
            accessor: "severity",
            disableSort: true
        },
        {
            Header: "Package name",
            id: "packageName",
            accessor: "packageName",
            disableSort: true
        },
        {
            Header: "Package version",
            id: "packageVersion",
            accessor: "packageVersion",
            disableSort: true
        },
        {
            Header: "Fix version",
            id: "fixVersion",
            accessor: "fixVersion",
            disableSort: true
        },
        {
            Header: "Exploits",
            id: "exploits",
            accessor: "exploits",
            disableSort: true
        },
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

export default VulnerabilitiesTable;