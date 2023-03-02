import React, { useMemo } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import EmptyDisplay from 'components/EmptyDisplay';
import Table, { utils } from 'components/Table';
import Icon, { ICON_NAMES } from 'components/Icon';
import ProgressBar from 'components/ProgressBar';
import { ExpandableScopeDisplay } from 'layout/Scans/scopeDisplayUtils';
import { SCAN_CONFIGS_PATH } from 'layout/Scans/Configurations';
import { useModalDisplayDispatch, MODAL_DISPLAY_ACTIONS } from 'layout/Scans/ScanConfigWizardModal/ModalDisplayProvider';
import { APIS } from 'utils/systemConsts';
import { formatDate } from 'utils/utils';
import VulnerabilitiesDisplay from '../VulnerabilitiesDisplay';
import ScanActionsDisplay from '../ScanActionsDisplay'

const TABLE_TITLE = "scans";

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
            accessor: "scanConfigId"
        },
        {
            Header: "Started",
            id: "startTime",
            Cell: ({row}) => <TimeDisplay time={row.original.startTime} />
        },
        {
            Header: "Ended",
            id: "endTime",
            Cell: ({row}) => <TimeDisplay time={row.original.endTime} />
        },
        {
            Header: "Scope",
            id: "scope",
            Cell: ({row}) => {
                const {all, regions} = row.original;

                return <ExpandableScopeDisplay all={all} regions={regions || []} />
            }
        },
        {
            Header: "Status",
            id: "status",
            Cell: ({row}) => <ProgressBar itemsCompleted={10} itemsLeft={8} width="80px" />
        },
        {
            Header: <Icon name={ICON_NAMES.SHIELD} />,
            id: "vulnerabilities",
            Cell: ({row}) => {
                const {id} = row.original;

                return (
                    <VulnerabilitiesDisplay id={id} isMinimized />
                )
            }
        },
        {
            Header: <Icon name={ICON_NAMES.BOMB} />,
            id: "exploits",
            accessor: "exploits"
        },
        {
            Header: "Scanned assets",
            id: "assets",
            accessor: original => "10/10"
        },
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
                    actionsComponent={({original}) => (
                        <ScanActionsDisplay data={original} />
                    )}
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
                            onSubClick={() => navigate(SCAN_CONFIGS_PATH)}
                        />
                    )}
                />
            </ContentContainer>
        </div>
    )
}

export default ScansTable;