import React, { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import TablePage from 'components/TablePage';
import { utils } from 'components/Table';
import ScanProgressBar, { SCAN_STATUS_ITEMS } from 'components/ScanProgressBar';
import EmptyDisplay from 'components/EmptyDisplay';
import { OPERATORS } from 'components/Filter';
import { ExpandableScopeDisplay } from 'layout/Scans/scopeDisplayUtils';
import { useModalDisplayDispatch, MODAL_DISPLAY_ACTIONS } from 'layout/Scans/ScanConfigWizardModal/ModalDisplayProvider';
import { APIS } from 'utils/systemConsts';
import { formatDate, getFindingsColumnsConfigList, getVulnerabilitiesColumnConfigItem, formatNumber, findingsColumnsFiltersConfig,
    vulnerabilitiesCountersColumnsFiltersConfig } from 'utils/utils';
import { FILTER_TYPES } from 'context/FiltersProvider';
import { SCANS_PATHS } from '../utils';
// import ScanActionsDisplay from '../ScanActionsDisplay';

const TABLE_TITLE = "scans";
const TIME_CELL_WIDTH = 110;

const START_TIME_SORT_IDS = ["startTime"];

const TimeDisplay = ({time}) => (
    !!time ? formatDate(time) : <utils.EmptyValue />
);

const ScansTable = () => {
    const modalDisplayDispatch = useModalDisplayDispatch();

    const navigate = useNavigate();

    const columns = useMemo(() => [
        {
            Header: "Config Name",
            id: "name",
            sortIds: ["scanConfigSnapshot.name"],
            accessor: "scanConfigSnapshot.name"
        },
        {
            Header: "Started",
            id: "startTime",
            sortIds: START_TIME_SORT_IDS,
            Cell: ({row}) => <TimeDisplay time={row.original.startTime} />,
            width: TIME_CELL_WIDTH
        },
        {
            Header: "Ended",
            id: "endTime",
            sortIds: ["endTime"],
            Cell: ({row}) => <TimeDisplay time={row.original.endTime} />,
            width: TIME_CELL_WIDTH
        },
        {
            Header: "Scope",
            id: "scope",
            sortIds: [
                "scanConfigSnapshot.scope.allRegions",
                "scanConfigSnapshot.scope.regions"
            ],
            Cell: ({row}) => {
                const {allRegions, regions} = row.original.scanConfigSnapshot?.scope;

                return <ExpandableScopeDisplay all={allRegions} regions={regions || []} />
            },
            width: 260
        },
        {
            Header: "Status",
            id: "status",
            sortIds: ["state"],
            Cell: ({row}) => {
                const {id, state, stateReason, stateMessage, summary} = row.original;
                const {jobsCompleted, jobsLeftToRun} = summary || {};

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
            width: 150
        },
        getVulnerabilitiesColumnConfigItem(TABLE_TITLE),
        ...getFindingsColumnsConfigList(TABLE_TITLE),
        {
            Header: "Scanned assets",
            id: "assets",
            sortIds: ["summary.jobsCompleted"],
            accessor: original => {
                const {jobsCompleted, jobsLeftToRun} = original.summary || {};
                
                return `${formatNumber(jobsCompleted)}/${formatNumber(jobsCompleted + jobsLeftToRun)}`;
            }
        }
    ], []);
    
    return (
        <TablePage
            columns={columns}
            url={APIS.SCANS}
            tableTitle={TABLE_TITLE}
            filterType={FILTER_TYPES.SCANS}
            filtersConfig={[
                {value: "scanConfigSnapshot.name", label: "Config name", operators: [
                    {...OPERATORS.eq, valueItems: [], creatable: true},
                    {...OPERATORS.ne, valueItems: [], creatable: true},
                    {...OPERATORS.contains, valueItems: [], creatable: true}
                ]},
                {value: "startTime", label: "Started", isDate: true, operators: [
                    {...OPERATORS.ge},
                    {...OPERATORS.le},
                ]},
                {value: "endTime", label: "Ended", isDate: true, operators: [
                    {...OPERATORS.ge},
                    {...OPERATORS.le},
                ]},
                {value: "scanConfigSnapshot.scope.regions", label: "Scope", operators: [
                    {...OPERATORS.contains, valueItems: [], creatable: true}
                ]},
                {value: "state", label: "Status", operators: [
                    {...OPERATORS.eq, valueItems: SCAN_STATUS_ITEMS},
                    {...OPERATORS.ne, valueItems: SCAN_STATUS_ITEMS}
                ]},
                ...vulnerabilitiesCountersColumnsFiltersConfig,
                ...findingsColumnsFiltersConfig,
                {value: "summary.jobsCompleted", label: "Scanned assets", isNumber: true, operators: [
                    {...OPERATORS.eq, valueItems: [], creatable: true},
                    {...OPERATORS.ne, valueItems: [], creatable: true},
                    {...OPERATORS.ge},
                    {...OPERATORS.le},
                ]}
            ]}
            defaultSortBy={{sortIds: START_TIME_SORT_IDS, desc: true}}
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
            absoluteSystemBanner
        />
    )
}

export default ScansTable;