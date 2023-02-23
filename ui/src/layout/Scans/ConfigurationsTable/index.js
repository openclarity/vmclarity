import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { isNull } from 'lodash';
import { useDelete, usePrevious } from 'hooks';
import ButtonWithIcon from 'components/ButtonWithIcon';
import Icon, { ICON_NAMES } from 'components/Icon';
import ContentContainer from 'components/ContentContainer';
import EmptyDisplay from 'components/EmptyDisplay';
import Table from 'components/Table';
import { TooltipWrapper } from 'components/Tooltip';
import Modal from 'components/Modal';
import ExpandableList from 'components/ExpandableList';
import { BoldText, toCapitalized, formatDate } from 'utils/utils';
import { formatTagsToStringInstances, formatRegionsToStrings } from '../utils';

import './configurations-table.scss';

const TABLE_TITLE = "scan configurations";

const SCAN_CONFIGS_URL = "scanConfigs";

export const SCAN_CONFIGS_PATH = "configs";

const ConfigurationsTable = ({setScanConfigFormData}) => {
    const columns = useMemo(() => [
        {
            Header: "Name",
            id: "name",
            accessor: "name"
        },
        {
            Header: "Scope",
            id: "scope",
            Cell: ({row}) => {
                const {all, regions} = row.original.scope;

                return (
                    all ? "All" : <ExpandableList items={formatRegionsToStrings(regions)} />
                )
            }
        },
        {
            Header: "Excluded instances",
            id: "instanceTagExclusion",
            Cell: ({row}) => {
                const {instanceTagExclusion} = row.original.scope;
                
                return (
                    <ExpandableList items={formatTagsToStringInstances(instanceTagExclusion)} withTagWrap />
                )
            }
        },
        {
            Header: "Included instances",
            id: "instanceTagSelector",
            Cell: ({row}) => {
                const {instanceTagSelector} = row.original.scope;
                
                return (
                    <ExpandableList items={formatTagsToStringInstances(instanceTagSelector)} withTagWrap />
                )
            }
        },
        {
            Header: "Time config",
            id: "timeConfig",
            Cell: ({row}) => {
                const {operationTime} = row.original.scheduled;
                const isScheduled = (Date.now() - (new Date(operationTime)).valueOf() <= 0);
                
                return (
                    <div>
                        {!!isScheduled && <BoldText>Scheduled</BoldText>}
                        <div>{formatDate(operationTime)}</div>
                    </div>
                )
            }
        },
        {
            Header: "Scan types",
            id: "scanTypes",
            Cell: ({row}) => {
                const {scanFamiliesConfig} = row.original;

                return (
                    <div>
                        {
                            Object.keys(scanFamiliesConfig).map(type => {
                                const {enabled} = scanFamiliesConfig[type];

                                return enabled ? toCapitalized(type) : null;
                            }).filter(type => !isNull(type)).join(" - ")
                        }
                    </div>
                )
            }
        }
    ], []);

    const [refreshTimestamp, setRefreshTimestamp] = useState(Date());
    const doRefreshTimestamp = useCallback(() => setRefreshTimestamp(Date()), []);

    const [deleteConfigmationData, setDeleteConfigmationData] = useState(null);
    const closeDeleteConfigmation = () => setDeleteConfigmationData(null);

    const [{deleting}, deleteScan] = useDelete(SCAN_CONFIGS_URL);
    const prevDeleting = usePrevious(deleting);

    useEffect(() => {
        if (prevDeleting && !deleting) {
            doRefreshTimestamp();
        }
    }, [prevDeleting, deleting, doRefreshTimestamp])

    return (
        <div className="scan-configs-table-page-wrapper">
            <ButtonWithIcon iconName={ICON_NAMES.PLUS} onClick={() => setScanConfigFormData({})}>
                New scan configuration
            </ButtonWithIcon>
            <ContentContainer>
                <Table
                    columns={columns}
                    paginationItemsName={TABLE_TITLE.toLowerCase()}
                    url={SCAN_CONFIGS_URL}
                    refreshTimestamp={refreshTimestamp}
                    noResultsTitle={TABLE_TITLE}
                    actionsComponent={({original}) => {
                        const {id} = original;
                        const deleteTooltipId = `${id}-delete`;
                        const editTooltipId = `${id}-edit`;
    
                        return (
                            <div className="config-row-actions">
                                <TooltipWrapper tooltipId={editTooltipId} tooltipText="Edit scan configuration" >
                                    <Icon
                                        name={ICON_NAMES.EDIT}
                                        onClick={event => {
                                            event.stopPropagation();
                                            event.preventDefault();
                                            
                                            setScanConfigFormData(original);
                                        }}
                                    />
                                </TooltipWrapper>
                                <TooltipWrapper tooltipId={deleteTooltipId} tooltipText="Delete scan configuration" >
                                    <Icon
                                        name={ICON_NAMES.DELETE}
                                        onClick={event => {
                                            event.stopPropagation();
                                            event.preventDefault();
    
                                            setDeleteConfigmationData(original);
                                        }}
                                    />
                                </TooltipWrapper>
                            </div>
                        );
                    }}
                    customEmptyResultsDisplay={() => (
                        <EmptyDisplay
                            message={(
                                <>
                                    <div>No scan configurations detected.</div>
                                    <div>Create your first scan configuration to see your VM's issues.</div>
                                </>
                            )}
                            title="New scan configuration"
                            onClick={() => setScanConfigFormData({})}
                        />
                    )}
                />
            </ContentContainer>
            {!isNull(deleteConfigmationData) &&
                <Modal
                    title="Delete configmation"
                    className="scan-config-delete-confirmation"
                    onClose={closeDeleteConfigmation}
                    height={250}
                    doneTitle="Delete"
                    onDone={() => {
                        deleteScan(deleteConfigmationData.id);
                        closeDeleteConfigmation();
                    }}
                >
                    <span>{`Once `}</span><BoldText>{deleteConfigmationData.name}</BoldText><span>{` will be deleted, the action cannot be reverted`}</span><br />
                    <span>{`Are you sure you want to delete ${deleteConfigmationData.name}?`}</span>
                </Modal>
            }
        </div>
    )
}

export default ConfigurationsTable;