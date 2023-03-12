import React from 'react';
import { useLocation } from 'react-router-dom';
import DetailsPageWrapper from 'components/DetailsPageWrapper';
import TabbedPage from 'components/TabbedPage';
import { APIS } from 'utils/systemConsts';
import { AssetDetails, ScanDetails } from 'layout/detail-displays';
import TabPackageDetails from './TabPackageDetails';

const PACKAGE_DETAILS_PATHS = {
    PACKAGE_DETAILS: "",
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
                    title: "Package details",
                    isIndex: true,
                    component: () => <TabPackageDetails data={data} />
                },
                {
                    id: "asset",
                    title: "Asset details",
                    path: PACKAGE_DETAILS_PATHS.ASSET_DETAILS,
                    component: () => <AssetDetails assetData={data} />
                },
                {
                    id: "scan",
                    title: "Scan details",
                    path: PACKAGE_DETAILS_PATHS.SCAN_DETAILS,
                    component: () => <ScanDetails scanData={data} />
                }
            ]}
            withInnerPadding={false}
        />
    )
}

const PackageDetails = () => (
    <DetailsPageWrapper
        // className="asset-details-page-wrapper"
        backTitle="Packages"
        getUrl={({id}) => `${APIS.SCANS}/${id}?$expand=ScanConfig`}
        getTitleData={({scanConfigSnapshot, startTime}) => ({title: "zlib1g", subTitle: "1.2.11+deb1"})}
        detailsContent={props => <DetailsContent {...props} />}
    />
)

export default PackageDetails;