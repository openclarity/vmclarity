import React from 'react';
import { useLocation } from 'react-router-dom';
import { useFetch } from 'hooks';
import TabbedPage from 'components/TabbedPage';
import FindingsDetailsPage from '../FindingsDetailsPage';
import TabRootkitDetails from './TabRootkitDetails';
import Loader from 'components/Loader';
import { APIS } from 'utils/systemConsts';

const ROOTKIT_DETAILS_PATHS = {
    PACKAGE_DETAILS: "",
}

const DetailsContent = ({data}) => {
    const {pathname} = useLocation();
    
    const {id} = data;
    
    const filter = `finding/id eq '${id}'`;

    const [{loading, data: assetFindingData, error}] = useFetch(APIS.ASSET_FINDINGS, {
        queryParams: {"$filter": filter, "$count": true, "$expand": "asset"}
    });

    if (error) {
        return null;
    }

    if (loading) {
        return <Loader large />
    }

    data.assetCount = assetFindingData.count;

    return (
        <TabbedPage
            basePath={`${pathname.substring(0, pathname.indexOf(id))}${id}`}
            items={[
                {
                    id: "general",
                    title: "Rootkit details",
                    isIndex: true,
                    component: () => <TabRootkitDetails data={data} />
                }
            ]}
            withInnerPadding={false}
        />
    )
}

const RootkitDetails = () => (
    <FindingsDetailsPage
        backTitle="Rootkits"
        getTitleData={({findingInfo}) => ({title: findingInfo.rootkitName})}
        detailsContent={DetailsContent}
    />
)

export default RootkitDetails;
