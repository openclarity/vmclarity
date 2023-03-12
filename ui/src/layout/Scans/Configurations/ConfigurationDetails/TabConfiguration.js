import React from 'react';
import { useNavigate } from 'react-router-dom';
import TitleValueDisplay, { TitleValueDisplayColumn } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Button from 'components/Button';
import Title from 'components/Title';
import { ScopeDisplay, ScanTypesDisplay, InstancesDisplay } from 'layout/Scans/scopeDisplayUtils';
import { ROUTES } from 'utils/systemConsts';

const TabConfiguration = ({data}) => {
    const navigate = useNavigate();

    const {scope, scanFamiliesConfig} = data || {};
    const {all, regions, instanceTagSelector, instanceTagExclusion} = scope;
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <Title medium>Configuration</Title>
                    <TitleValueDisplayColumn>
                        <TitleValueDisplay title="Scope"><ScopeDisplay all={all} regions={regions} /></TitleValueDisplay>
                        <TitleValueDisplay title="Included instances"><InstancesDisplay tags={instanceTagSelector}/></TitleValueDisplay>
                        <TitleValueDisplay title="Excluded instances"><InstancesDisplay tags={instanceTagExclusion}/></TitleValueDisplay>
                        <TitleValueDisplay title="Scan types"><ScanTypesDisplay scanFamiliesConfig={scanFamiliesConfig} /></TitleValueDisplay>
                    </TitleValueDisplayColumn>
                </>
            )}
            rightPlaneDisplay={() => (
                <>
                    <Title medium>Configuration's scans</Title>
                    <Button onClick={() => navigate(ROUTES.SCANS)}>See all scans (100)</Button>
                </>
            )}
        />
    )
}

export default TabConfiguration;