import React from 'react';
import { useFetch } from 'hooks';
import TitleValueDisplay, { TitleValueDisplayColumn, TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Loader from 'components/Loader';
import Title from 'components/Title';
import { ScopeDisplay, ScanTypesDisplay, InstancesDisplay } from 'layout/Scans/scopeDisplayUtils';
import { APIS } from 'utils/systemConsts';
import { formatDate } from 'utils/utils';
import ConfigurationAlertLink from './ConfigurationAlertLink';

const TabGeneral = ({data}) => {
    const {scanConfigId, scanFamiliesConfig, startTime, endTime} = data || {};

    const [{loading, data: configData, error}] = useFetch(`${APIS.SCAN_CONFIGS}/${scanConfigId}`);

    if (loading) {
        return <Loader />;
    }

    if (error) {
        return null;
    }
    
    const {all, regions, instanceTagSelector, instanceTagExclusion} = configData?.scope;
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <TitleValueDisplayColumn>
                    <ConfigurationAlertLink configData={configData} />
                    <TitleValueDisplay title="Scope"><ScopeDisplay all={all} regions={regions} /></TitleValueDisplay>
                    <TitleValueDisplay title="Included instances"><InstancesDisplay tags={instanceTagSelector}/></TitleValueDisplay>
                    <TitleValueDisplay title="Excluded instances"><InstancesDisplay tags={instanceTagExclusion}/></TitleValueDisplay>
                    <TitleValueDisplay title="Scan types"><ScanTypesDisplay scanFamiliesConfig={scanFamiliesConfig} /></TitleValueDisplay>
                </TitleValueDisplayColumn>
            )}
            rightPlaneDisplay={() => (
                <>
                    <Title medium>Status</Title>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Started">{formatDate(startTime)}</TitleValueDisplay>
                        <TitleValueDisplay title="Ended">{formatDate(endTime)}</TitleValueDisplay>
                        <TitleValueDisplay title="Duration"></TitleValueDisplay>
                    </TitleValueDisplayRow>
                </>
            )}
        />
    )
}

export default TabGeneral;