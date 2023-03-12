import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import { APIS } from 'utils/systemConsts';

const TABLE_TITLE = "secrets";

const SecretsTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();
    
    const columns = useMemo(() => [
        {
            Header: "Fingerprint",
            id: "fingerprint",
            accessor: "id",
            disableSort: true
        },
        {
            Header: "Description",
            id: "description",
            accessor: "description",
            disableSort: true
        },
        {
            Header: "FilePath",
            id: "filePath",
            accessor: "filePath",
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

export default SecretsTable;