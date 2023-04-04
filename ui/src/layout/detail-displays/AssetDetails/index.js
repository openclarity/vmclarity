import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useFetch } from 'hooks';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Title from 'components/Title';
import Button from 'components/Button';
import Loader from 'components/Loader';
import { ROUTES, APIS } from 'utils/systemConsts';
import { formatNumber } from 'utils/utils';
import { useFilterDispatch, setFilters, FILTER_TYPES } from 'context/FiltersProvider';

const AssetScansDisplay = ({assetName, targetId}) => {
    const {pathname} = useLocation();
    const navigate = useNavigate();
    const filtersDispatch = useFilterDispatch();

    const filter = `target/id eq '${targetId}'`;
    
    const onAssetScansClick = () => {
        setFilters(filtersDispatch, {
            type: FILTER_TYPES.ASSET_SCANS,
            filters: {filter, name: assetName, suffix: "asset", backPath: pathname},
            isSystem: true
        });

        navigate(ROUTES.ASSET_SCANS);
    }
    
    const [{loading, data, error}] = useFetch(APIS.ASSET_SCANS, {
        queryParams: {"$filter": filter, "$count": true, "$select": "id,target,summary,scan"}
    });
    
    if (error) {
        return null;
    }

    if (loading) {
        return <Loader absolute={false} small />
    }
    
    return (
        <>
            <Title medium>Asset scans</Title>
            <Button onClick={onAssetScansClick} >{`See all asset scans (${formatNumber(data?.count || 0)})`}</Button>
        </>
    )
}

const AssetDetails = ({assetData, withAssetScansLink=false}) => {
    const {id, targetInfo} = assetData;
    const {instanceID, objectType, location} = targetInfo || {};

    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <Title medium>Asset</Title>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Name">{instanceID}</TitleValueDisplay>
                        <TitleValueDisplay title="Type">{objectType}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                    <TitleValueDisplayRow>
                        <TitleValueDisplay title="Location">{location}</TitleValueDisplay>
                    </TitleValueDisplayRow>
                </>
            )}
            rightPlaneDisplay={!withAssetScansLink ? null : () => <AssetScansDisplay assetName={instanceID} targetId={id} />}
        />
    )
}

export default AssetDetails;