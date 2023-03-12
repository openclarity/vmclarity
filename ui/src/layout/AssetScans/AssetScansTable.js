import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import { APIS } from 'utils/systemConsts';
import { getScanName, getFindingsColumnsConfigList, getVulnerabilitiesColumnConfigItem } from 'utils/utils';

const TABLE_TITLE = "asset scans";

const AssetScansTable = () => {
    const navigate = useNavigate();
    const {pathname} = useLocation();

    const columns = useMemo(() => [
        {
            Header: "Asset name",
            id: "name",
            accessor: "target.targetInfo.instanceID",
            disableSort: true
        },
        {
            Header: "Asset type",
            id: "type",
            accessor: "target.targetInfo.objectType",
            disableSort: true
        },
        {
            Header: "Asset location",
            id: "location",
            accessor: "target.targetInfo.location",
            disableSort: true
        },
        {
            Header: "Scan",
            id: "scan",
            accessor: original => {
                const {startTime, scanConfigSnapshot} = original.scan;
                
                return getScanName({name: scanConfigSnapshot?.name, startTime});
            },
            disableSort: true
        },
        getVulnerabilitiesColumnConfigItem({tableTitle: TABLE_TITLE, idKey: "scan.id", summaryKey: "scan.summary"}),
        ...getFindingsColumnsConfigList({tableTitle: TABLE_TITLE, summaryKey: "scan.summary"})
    ], []);

    return (
        <ContentContainer withMargin>
            <Table
                columns={columns}
                paginationItemsName={TABLE_TITLE.toLowerCase()}
                url={APIS.ASSET_SCANS}
                noResultsTitle={TABLE_TITLE}
                onLineClick={({id}) => navigate(`${pathname}/${id}`)}
            />
        </ContentContainer>
    )
}

export default AssetScansTable;