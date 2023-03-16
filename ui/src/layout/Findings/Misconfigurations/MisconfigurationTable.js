import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import { APIS } from 'utils/systemConsts';
import { getAssetAndScanColumnsConfigList } from 'layout/Findings/utils';

const TABLE_TITLE = "misconfiguration";

const MisconfigurationsTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();
    
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
        <ContentContainer>
            <Table
                columns={columns}
                paginationItemsName={TABLE_TITLE.toLowerCase()}
                url={APIS.FINDINGS}
                filters={{"$filter": `findingInfo/objectType eq 'Misconfiguration'`, "$expand": "asset,scan"}}
                noResultsTitle={TABLE_TITLE}
                onLineClick={({id}) => navigate(`${pathname}/${id}`)}
            />
        </ContentContainer>
    )
}

export default MisconfigurationsTable;