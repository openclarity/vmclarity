import React from 'react';
import { useNavigate } from 'react-router-dom';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Title from 'components/Title';
import Button from 'components/Button';
import { ROUTES } from 'utils/systemConsts';

const AssetScansDisplay = ({onClick}) => (
    <>
        <Title medium>Asset scans</Title>
        <Button onClick={onClick}>See all asset scans (100)</Button>
    </>
)

const AssetDetails = ({assetData, withAssetLink=false, withAssetScansLink=false}) => {
    const navigate = useNavigate();

    const {id, targetInfo} = assetData;
    const {instanceID, objectType, location} = targetInfo || {};

    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <Title medium onClick={withAssetLink ? () => navigate(`${ROUTES.ASSETS}/${id}`) : undefined}>Asset</Title>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Name">{instanceID}</TitleValueDisplay>
                        <TitleValueDisplay title="Type">{objectType}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Location">{location}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                </>
            )}
            rightPlaneDisplay={!withAssetScansLink ? null : () => <AssetScansDisplay onClick={() => navigate(ROUTES.ASSET_SCANS)} />}
        />
    )
}

export default AssetDetails;