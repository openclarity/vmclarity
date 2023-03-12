import React from 'react';
import moment from 'moment';
import TitleValueDisplay, { TitleValueDisplayColumn, TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Title from 'components/Title';
import ScanProgressBar from 'components/ScanProgressBar';
import { ScopeDisplay, ScanTypesDisplay, InstancesDisplay } from 'layout/Scans/scopeDisplayUtils';
import { formatDate } from 'utils/utils';
import ConfigurationAlertLink from './ConfigurationAlertLink';
import Button from 'components/Button';
import { useNavigate } from 'react-router-dom';

export const calculateDuration = (startTime, endTime) => {
    const startMoment = moment(startTime);
    const endMoment = moment(endTime);

    const range = ["days", "hours", "minutes", "seconds"].map(item => ({diff: endMoment.diff(startMoment, item), label: item}))
        .find(({diff}) => diff > 1);

    return !!range ? `${range.diff} ${range.label}` : null;
}

const ScanDetails = ({scanData, withAssetScansLink=false}) => {
    const navigate = useNavigate();

    const {scanConfig, scanConfigSnapshot, startTime, endTime, summary, state, stateMessage, stateReason} = scanData || {};
    const {scope, scanFamiliesConfig} = scanConfigSnapshot;
    const {all, regions, instanceTagSelector, instanceTagExclusion, shouldScanStoppedInstances} = scope;
    const {jobsCompleted, jobsLeftToRun} = summary;
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <TitleValueDisplayColumn>
                    <ConfigurationAlertLink updatedConfigData={scanConfig} scanConfigData={scanConfigSnapshot} />
                    <TitleValueDisplay title="Scope"><ScopeDisplay all={all} regions={regions} /></TitleValueDisplay>
                    <TitleValueDisplay title="Instances">
                        <div style={{margin: "10px 0 20px 0"}}>
                            {shouldScanStoppedInstances ? "Running and non-running instances" : "Running instances only"}
                        </div>
                        <TitleValueDisplay title="Included instances" isSubItem><InstancesDisplay tags={instanceTagSelector}/></TitleValueDisplay>
                        <TitleValueDisplay title="Excluded instances" isSubItem><InstancesDisplay tags={instanceTagExclusion}/></TitleValueDisplay>
                    </TitleValueDisplay>
                    <TitleValueDisplay title="Scan types"><ScanTypesDisplay scanFamiliesConfig={scanFamiliesConfig} /></TitleValueDisplay>
                </TitleValueDisplayColumn>
            )}
            rightPlaneDisplay={() => (
                <>
                    <Title medium>Status</Title>
                    <ScanProgressBar
                        state={state}
                        stateReason={stateReason}
                        stateMessage={stateMessage}
                        itemsCompleted={jobsCompleted}
                        itemsLeft={jobsLeftToRun}
                    />
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Started">{formatDate(startTime)}</TitleValueDisplay>
                        <TitleValueDisplay title="Ended">{formatDate(endTime)}</TitleValueDisplay>
                        <TitleValueDisplay title="Duration">{calculateDuration(startTime, endTime)}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    {withAssetScansLink &&
                        <div style={{marginTop: "50px"}}>
                            <Title medium>Asset scans</Title>
                            <Button onClick={() => navigate(0)}>See asset scans (221/340)</Button>
                        </div>
                    }
                </>
            )}
        />
    )
}

export default ScanDetails;