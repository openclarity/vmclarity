import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import { APIS } from 'utils/systemConsts';
import { getFindingsColumnsConfigList, getVulnerabilitiesColumnConfigItem } from 'utils/utils';

const TABLE_TITLE = "assets";

const AssetsTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();

    const columns = useMemo(() => [
        {
            Header: "Name",
            id: "name",
            accessor: "targetInfo.instanceID",
            disableSort: true
        },
        {
            Header: "Type",
            id: "type",
            accessor: "targetInfo.objectType",
            disableSort: true
        },
        {
            Header: "Location",
            id: "location",
            accessor: "targetInfo.location",
            disableSort: true
        },
        getVulnerabilitiesColumnConfigItem({tableTitle: TABLE_TITLE}),
        ...getFindingsColumnsConfigList({tableTitle: TABLE_TITLE})
    ], []);

    return (
        <ContentContainer withMargin>
            <Table
                columns={columns}
                paginationItemsName={TABLE_TITLE.toLowerCase()}
                url={APIS.ASSETS}
                noResultsTitle={TABLE_TITLE}
                onLineClick={({id}) => navigate(`${pathname}/${id}`)}
            />
        </ContentContainer>
    )
}

export default AssetsTable;