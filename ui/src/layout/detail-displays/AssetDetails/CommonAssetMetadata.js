import React from 'react';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import { formatDate } from 'utils/utils';


export const CommonAssetMetadata = ({assetData}) => {
    const {firstSeen, lastSeen, terminatedOn} = assetData;

    return (
        <>
            <TitleValueDisplayRow>
                <TitleValueDisplay title="Type">{assetData.assetInfo.objectType}</TitleValueDisplay>
                <TitleValueDisplay title="First Seen">{formatDate(firstSeen)}</TitleValueDisplay>
            </TitleValueDisplayRow>
            <TitleValueDisplayRow>
                
                <TitleValueDisplay title="Last Seen">{formatDate(lastSeen)}</TitleValueDisplay>
                <TitleValueDisplay title="Terminated On">{formatDate(terminatedOn)}</TitleValueDisplay>
            </TitleValueDisplayRow>
        </>
    )
}
