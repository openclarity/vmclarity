import React from 'react';
import { useLocation } from 'react-router-dom';
import DetailsPageWrapper from 'components/DetailsPageWrapper';
import TabbedPage from 'components/TabbedPage';
import { APIS } from 'utils/systemConsts';
import { Findings } from 'layout/detail-displays';
import TabAssetScanDetails from './TabAssetScanDetails';

const ASSET_SCAN_DETAILS_PATHS = {
    ASSET_SCAN_DETAILS: "",
    FINDINGHS: "findings"
}

const DetailsContent = ({data}) => {
    const {pathname} = useLocation();
    
    const {id, name, summary} = data;
    
    return (
        <TabbedPage
            basePath={`${pathname.substring(0, pathname.indexOf(id))}${id}`}
            items={[
                {
                    id: "general",
                    title: "Asset scan details",
                    isIndex: true,
                    component: () => <TabAssetScanDetails data={data} />
                },
                {
                    id: "findings",
                    title: "Findings",
                    path: ASSET_SCAN_DETAILS_PATHS.FINDINGHS,
                    component: () => (
                        <Findings
                            findingsSummary={summary}
                            findingsFilter={`assetScan/id eq '${id}'`}
                            findingsFilterTitle={`${name}`}
                        />
                    )
                }
            ]}
            withInnerPadding={false}
        />
    )
}

const AssetScanDetails = () => (
    <DetailsPageWrapper
        backTitle="Asset scans"
        getUrl={({id}) => `${APIS.ASSET_SCANS}/${id}?$select=id,name,summary,status&$expand=scan($select=id,name,startTime,endTime),target($select=id,targetInfo/objectType,targetInfo/location,targetInfo/instanceID)`}
        getTitleData={({name}) => {
            return ({
                title: name,
            })
        }}
        detailsContent={props => <DetailsContent {...props} />}
        withPadding
    />
)

export default AssetScanDetails;
