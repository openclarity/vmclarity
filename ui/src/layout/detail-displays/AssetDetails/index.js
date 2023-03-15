import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useFetch } from 'hooks';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Title from 'components/Title';
import Button from 'components/Button';
import Loader from 'components/Loader';
import { ROUTES, APIS } from 'utils/systemConsts';

const AssetScansDisplay = ({targetId}) => {
    const navigate = useNavigate();

    const [{loading, data, error}] = useFetch(APIS.ASSET_SCANS, {queryParams: {"$filter": `target/id eq '${targetId}'`, "$count": true}});
    
    if (error) {
        return null;
    }

    if (loading) {
        return <Loader />
    }
    
    return (
        <>
            <Title medium>Asset scans</Title>
            <Button onClick={() => navigate(ROUTES.ASSET_SCANS)} >{`See all asset scans (${data?.count || 0})`}</Button>
        </>
    )
}

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
            rightPlaneDisplay={!withAssetScansLink ? null : () => <AssetScansDisplay targetId={id} />}
        />
    )
}

export default AssetDetails;