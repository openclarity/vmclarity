import React from 'react';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import { FindingsDetailsCommonFields } from '../utils';
import { formatNumber } from 'utils/utils';

const TabPackageDetails = ({data}) => {
    const {findingInfo, firstSeen, lastSeen, assetCount} = data;
    const {rootkitName, message} = findingInfo;

    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Rootkit name">{rootkitName}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Message">{message}</TitleValueDisplay>
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

export default TabPackageDetails;
