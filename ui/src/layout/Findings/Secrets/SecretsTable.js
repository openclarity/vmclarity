import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';

const TABLE_TITLE = "secrets";

const SecretsTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();
    
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
        <ContentContainer>
            <Table
                columns={columns}
                paginationItemsName={TABLE_TITLE.toLowerCase()}
                url={APIS.FINDINGS}
                filters={{"$filter": `findingInfo/objectType eq 'Secret'`, "$expand": "asset,scan"}}
                noResultsTitle={TABLE_TITLE}
                onLineClick={({id}) => navigate(`${pathname}/${id}`)}
            />
        </ContentContainer>
    )
}

export default SecretsTable;