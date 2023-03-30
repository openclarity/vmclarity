import React from 'react';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import { MISCONFIGURATION_SEVERITY_MAP } from './utils';

const TabMisconfigurationDetails = ({data}) => {
    const {testID, severity, testDescription, scannerName, scannedPath, remediation, testCategory, message} = data.findingInfo;

    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Test ID">{testID}</TitleValueDisplay>
                        <TitleValueDisplay title="Severity">{MISCONFIGURATION_SEVERITY_MAP[severity]}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Scanner Name">{scannerName}</TitleValueDisplay>
                        <TitleValueDisplay title="File path">{scannedPath}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Category">{testCategory}</TitleValueDisplay>
                        <TitleValueDisplay title="Remediation">{remediation}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Message" withOpen defaultOpen>{message}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Description" withOpen defaultOpen>{testDescription}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                </>  
            )}
        />
    )
}

export default TabMisconfigurationDetails;