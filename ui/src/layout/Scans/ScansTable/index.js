import React, { useMemo } from 'react';
import ContentContainer from 'components/ContentContainer';
import EmptyDisplay from 'components/EmptyDisplay';
import Table from 'components/Table';

const TABLE_TITLE = "scans";

const ScansTable = ({setScanConfigFormData}) => {
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
                            onSubClick={() => {debugger}}
                        />
                    )}
                />
            </ContentContainer>
        </div>
    )
}

export default ScansTable;