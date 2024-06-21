import React from 'react';
import { useLocation } from 'react-router-dom';
import { useFetch } from 'hooks';
import TabbedPage from 'components/TabbedPage';
import FindingsDetailsPage from '../FindingsDetailsPage';
import TabMisconfigurationDetails from './TabMisconfigurationDetails';
import Loader from 'components/Loader';
import { APIS } from 'utils/systemConsts';

const MISCONFIGURATION_DETAILS_PATHS = {
    MISCONFIGURATION_DETAILS: "",
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
                    title: "Misconfiguration details",
                    isIndex: true,
                    component: () => <TabMisconfigurationDetails data={data} />
                }
            ]}
            withInnerPadding={false}
        />
    )
}

const MisconfigurationDetails = () => (
    <FindingsDetailsPage
        backTitle="Misconfigurations"
        getTitleData={({findingInfo}) => ({title: findingInfo.id})}
        detailsContent={DetailsContent}
    />
)

export default MisconfigurationDetails;
