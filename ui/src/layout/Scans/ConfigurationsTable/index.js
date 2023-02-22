import React, { useEffect, useMemo, useState } from 'react';
import { useDelete, usePrevious } from 'hooks';
import ButtonWithIcon from 'components/ButtonWithIcon';
import Icon, { ICON_NAMES } from 'components/Icon';
import ContentContainer from 'components/ContentContainer';
import EmptyDisplay from 'components/EmptyDisplay';
import Table from 'components/Table';
import { TooltipWrapper } from 'components/Tooltip';

import './configurations-table.scss';

const TABLE_TITLE = "scan configurations";

const SCAN_CONFIGS_URL = "scanConfigs";

const ConfigurationsTable = ({setScanConfigFormData}) => {
    const columns = useMemo(() => [
        {
            Header: "Name",
            id: "name",
            accessor: "name"
        }
    ], []);

    const [refreshTimestamp, setRefreshTimestamp] = useState(Date());
    const doRefreshTimestamp = () => setRefreshTimestamp(Date());

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
    
                                            deleteScan(id)
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
        </div>
    )
}

export default ConfigurationsTable;