import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import EmptyDisplay from 'components/EmptyDisplay';
import Table from 'components/Table';
import { SCAN_CONFIGS_PATH } from 'layout/Scans/Configurations';

const TABLE_TITLE = "scans";

const ScansTable = ({setScanConfigFormData}) => {
    const navigate = useNavigate();
    const {pathname} = useLocation();

    const columns = useMemo(() => [
        {
            Header: "Config Name",
            id: "name",
            accessor: "scanConfigId"
        },
        {
            Header: "Started",
            id: "startTime",
            accessor: "startTime"
        },
        {
            Header: "Ended",
            id: "endTime",
            accessor: "endTime"
        }
    ], []);

    return (
        <div className="scans-table-page-wrapper">
            <ContentContainer>
                <Table
                    columns={columns}
                    paginationItemsName={TABLE_TITLE.toLowerCase()}
                    url="scans"
                    noResultsTitle={TABLE_TITLE}
                    onLineClick={({id}) => navigate(`${pathname}/${id}`)}
                    customEmptyResultsDisplay={() => (
                        <EmptyDisplay
                            message={(
                                <>
                                    <div>No scans detected.</div>
                                    <div>Start your first scan to see your VM's issues.</div>
                                </>
                            )}
                            title="New scan configuration"
                            onClick={() => setScanConfigFormData({})}
                            subTitle="Start scan from config"
                            onSubClick={() => navigate(SCAN_CONFIGS_PATH)}
                        />
                    )}
                />
            </ContentContainer>
        </div>
    )
}

export default ScansTable;