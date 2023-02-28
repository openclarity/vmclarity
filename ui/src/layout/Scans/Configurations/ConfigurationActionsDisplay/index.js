import React, { useEffect, useState } from 'react';
import { isNull } from 'lodash';
import { useDelete, usePrevious } from 'hooks';
import Icon, { ICON_NAMES } from 'components/Icon';
import { TooltipWrapper } from 'components/Tooltip';
import Modal from 'components/Modal';
import { BoldText } from 'utils/utils';
import { SCAN_CONFIGS_URL } from '../utils';

import './configuration-actions-display.scss';

const ConfigurationActionsDisplay = ({data, setScanConfigFormData, onDelete}) => {
    const {id} = data;

    const [deleteConfigmationData, setDeleteConfigmationData] = useState(null);
    const closeDeleteConfigmation = () => setDeleteConfigmationData(null);

    const [{deleting}, deleteConfiguration] = useDelete(SCAN_CONFIGS_URL);
    const prevDeleting = usePrevious(deleting);

    useEffect(() => {
        if (prevDeleting && !deleting) {
            onDelete();
        }
    }, [prevDeleting, deleting, onDelete])

    return (
        <>
            <div className="configuration-actions-display">
                <TooltipWrapper tooltipId={`${id}-duplicate`} tooltipText="Duplicate scan configuration" >
                    <Icon
                        name={ICON_NAMES.DUPLICATE}
                        onClick={event => {
                            event.stopPropagation();
                            event.preventDefault();
                            
                            setScanConfigFormData({...data, id: null, name: ""});
                        }}
                    />
                </TooltipWrapper>
                <TooltipWrapper tooltipId={`${id}-edit`} tooltipText="Edit scan configuration" >
                    <Icon
                        name={ICON_NAMES.EDIT}
                        onClick={event => {
                            event.stopPropagation();
                            event.preventDefault();
                            
                            setScanConfigFormData(data);
                        }}
                    />
                </TooltipWrapper>
                <TooltipWrapper tooltipId={`${id}-delete`} tooltipText="Delete scan configuration" >
                    <Icon
                        name={ICON_NAMES.DELETE}
                        onClick={event => {
                            event.stopPropagation();
                            event.preventDefault();

                            setDeleteConfigmationData(data);
                        }}
                    />
                </TooltipWrapper>
            </div>
            {!isNull(deleteConfigmationData) &&
                <Modal
                    title="Delete configmation"
                    className="scan-config-delete-confirmation"
                    onClose={closeDeleteConfigmation}
                    height={250}
                    doneTitle="Delete"
                    onDone={() => {
                        deleteConfiguration(deleteConfigmationData.id);
                        closeDeleteConfigmation();
                    }}
                >
                    <span>{`Once `}</span><BoldText>{deleteConfigmationData.name}</BoldText><span>{` will be deleted, the action cannot be reverted`}</span><br />
                    <span>{`Are you sure you want to delete ${deleteConfigmationData.name}?`}</span>
                </Modal>
            }
        </>
    );
}

export default ConfigurationActionsDisplay;