import React from 'react';
import { useLocation } from 'react-router-dom';
import DetailsPageWrapper from 'components/DetailsPageWrapper';
import TabbedPage from 'components/TabbedPage';
import { APIS } from 'utils/systemConsts';
import { AssetDetails, ScanDetails } from 'layout/detail-displays';
import TabSecretDetails from './TabSecretDetails';

const SECRET_DETAILS_PATHS = {
    SECRET_DETAILS: "",
    ASSET_DETAILS: "asset",
    SCAN_DETAILS: "scan"
}

const DetailsContent = ({data}) => {
    const {pathname} = useLocation();
    
    const {id} = data;
    
    return (
        <TabbedPage
            basePath={`${pathname.substring(0, pathname.indexOf(id))}${id}`}
            items={[
                {
                    id: "general",
                    title: "Secret details",
                    isIndex: true,
                    component: () => <TabSecretDetails data={data} />
                },
                {
                    id: "asset",
                    title: "Asset details",
                    path: SECRET_DETAILS_PATHS.ASSET_DETAILS,
                    component: () => <AssetDetails assetData={data} />
                },
                {
                    id: "scan",
                    title: "Scan details",
                    path: SECRET_DETAILS_PATHS.SCAN_DETAILS,
                    component: () => <ScanDetails scanData={data} />
                }
            ]}
            withInnerPadding={false}
        />
    )
}

const SecretDetails = () => (
    <DetailsPageWrapper
        // className="asset-details-page-wrapper"
        backTitle="Secrets"
        getUrl={({id}) => `${APIS.SCANS}/${id}?$expand=ScanConfig`}
        getTitleData={({scanConfigSnapshot, startTime}) => ({title: "Fingerprint"})}
        detailsContent={props => <DetailsContent {...props} />}
    />
)

export default SecretDetails;