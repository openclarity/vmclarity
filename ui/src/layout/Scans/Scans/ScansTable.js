import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import EmptyDisplay from 'components/EmptyDisplay';
import Table, { utils } from 'components/Table';
import ScanProgressBar from 'components/ScanProgressBar';
import { ExpandableScopeDisplay } from 'layout/Scans/scopeDisplayUtils';
import { useModalDisplayDispatch, MODAL_DISPLAY_ACTIONS } from 'layout/Scans/ScanConfigWizardModal/ModalDisplayProvider';
import { APIS } from 'utils/systemConsts';
import { SCANS_PATHS } from '../utils';
import { formatDate, getFindingsColumnsConfigList, getVulnerabilitiesColumnConfigItem } from 'utils/utils';
// import ScanActionsDisplay from '../ScanActionsDisplay';

const TABLE_TITLE = "scans";
const TIME_CELL_WIDTH = 110;

const TimeDisplay = ({time}) => (
    !!time ? formatDate(time) : <utils.EmptyValue />
);

const ScansTable = () => {
    const modalDisplayDispatch = useModalDisplayDispatch();

    const navigate = useNavigate();
    const {pathname} = useLocation();

    const columns = useMemo(() => [
        {
            Header: "Config Name",
            id: "name",
            accessor: "scanConfigSnapshot.name",
            disableSort: true
        },
        {
            Header: "Started",
            id: "startTime",
            Cell: ({row}) => <TimeDisplay time={row.original.startTime} />,
            width: TIME_CELL_WIDTH,
            disableSort: true
        },
        {
            Header: "Ended",
            id: "endTime",
            Cell: ({row}) => <TimeDisplay time={row.original.endTime} />,
            width: TIME_CELL_WIDTH,
            disableSort: true
        },
        {
            Header: "Scope",
            id: "scope",
            Cell: ({row}) => {
                const {all, regions} = row.original.scanConfigSnapshot?.scope;

                return <ExpandableScopeDisplay all={all} regions={regions || []} />
            },
            width: 260,
            disableSort: true
        },
        {
            Header: "Status",
            id: "status",
            Cell: ({row}) => {
                const {id, state, stateReason, stateMessage, summary} = row.original;
                const {jobsCompleted, jobsLeftToRun} = summary;

                return (
                    <ScanProgressBar
                        state={state}
                        stateReason={stateReason}
                        stateMessage={stateMessage}
                        itemsCompleted={jobsCompleted}
                        itemsLeft={jobsLeftToRun}
                        barWidth="80px"
                        isMinimized
                        minimizedTooltipId={id}
                    />
                )
            },
            width: 150,
            disableSort: true
        },
        getVulnerabilitiesColumnConfigItem({tableTitle: TABLE_TITLE}),
        ...getFindingsColumnsConfigList({tableTitle: TABLE_TITLE}),
        {
            Header: "Scanned assets",
            id: "assets",
            accessor: original => {
                const {jobsCompleted, jobsLeftToRun} = original.summary;
                
                return `${jobsCompleted}/${jobsCompleted + jobsLeftToRun}`;
            },
            disableSort: true
        }
    ], []);

    return (
        <div className="scans-table-page-wrapper">
            <ContentContainer>
                <Table
                    columns={columns}
                    paginationItemsName={TABLE_TITLE.toLowerCase()}
                    url={APIS.SCANS}
                    noResultsTitle={TABLE_TITLE}
                    onLineClick={({id}) => navigate(`${pathname}/${id}`)}
                    // actionsComponent={({original}) => (
                    //     <ScanActionsDisplay data={original} />
                    // )}
                    customEmptyResultsDisplay={() => (
                        <EmptyDisplay
                            message={(
                                <>
                                    <div>No scans detected.</div>
                                    <div>Start your first scan to see your VM's issues.</div>
                                </>
                            )}
                            title="New scan configuration"
                            onClick={() => modalDisplayDispatch({type: MODAL_DISPLAY_ACTIONS.SET_MODAL_DISPLAY_DATA, payload: {}})}
                            subTitle="Start scan from config"
                            onSubClick={() => navigate(SCANS_PATHS.CONFIGURATIONS)}
                        />
                    )}
                />
            </ContentContainer>
        </div>
    )
}

export default ScansTable;