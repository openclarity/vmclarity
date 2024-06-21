import React from 'react';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import { FindingsDetailsCommonFields } from '../utils';
import { formatNumber } from 'utils/utils';

const TabSecretDetails = ({data}) => {
    const {findingInfo, firstSeen, lastSeen, assetCount} = data;
    const {fingerprint, description, startLine, endLine, filePath} = findingInfo;

    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Fingerprint">{fingerprint}</TitleValueDisplay>
                        <TitleValueDisplay title="Description">{description}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Start Line">{startLine}</TitleValueDisplay>
                        <TitleValueDisplay title="End line">{endLine}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="File path">{filePath}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <FindingsDetailsCommonFields firstSeen={firstSeen} lastSeen={lastSeen} />
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Asset count">{formatNumber(assetCount)}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                </>
            )}
        />
    )
}

export default TabSecretDetails;
