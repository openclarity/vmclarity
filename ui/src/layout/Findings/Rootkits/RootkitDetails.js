import React from 'react';
import { useLocation } from 'react-router-dom';
import DetailsPageWrapper from 'components/DetailsPageWrapper';
import TabbedPage from 'components/TabbedPage';
import { APIS } from 'utils/systemConsts';
import { AssetDetails, ScanDetails } from 'layout/detail-displays';
import TabRootkitDetails from './TabRootkitDetails';

const ROOTKIT_DETAILS_PATHS = {
    PACKAGE_DETAILS: "",
    ASSET_DETAILS: "asset",
    SCAN_DETAILS: "scan"
}

const DetailsContent = ({data}) => {
    const {pathname} = useLocation();
    
    const {id, assetScan} = data;
    const {target: asset, scan} = assetScan;
    
    return (
        <TabbedPage
            basePath={`${pathname.substring(0, pathname.indexOf(id))}${id}`}
            items={[
                {
                    id: "general",
                    title: "Rootkit details",
                    isIndex: true,
                    component: () => <TabRootkitDetails data={data} />
                },
                {
                    id: "asset",
                    title: "Asset details",
                    path: ROOTKIT_DETAILS_PATHS.ASSET_DETAILS,
                    component: () => <AssetDetails assetData={asset} withAssetLink />
                },
                {
                    id: "scan",
                    title: "Scan details",
                    path: ROOTKIT_DETAILS_PATHS.SCAN_DETAILS,
                    component: () => <ScanDetails scanData={scan} withScanLink />
                }
            ]}
            withInnerPadding={false}
        />
    )
}

const RootkitDetails = () => (
    <DetailsPageWrapper
        backTitle="Rootkits"
        getUrl={({id}) => `${APIS.FINDINGS}/${id}?$expand=assetScan($expand=target,scan)`}
        getTitleData={({findingInfo}) => ({title: findingInfo.rootkitName})}
        detailsContent={props => <DetailsContent {...props} />}
    />
)

export default RootkitDetails;
