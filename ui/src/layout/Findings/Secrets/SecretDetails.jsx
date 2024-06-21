import React from 'react';
import { useLocation } from 'react-router-dom';
import { useFetch } from 'hooks';
import TabbedPage from 'components/TabbedPage';
import FindingsDetailsPage from '../FindingsDetailsPage';
import TabSecretDetails from './TabSecretDetails';
import Loader from 'components/Loader';
import { APIS } from 'utils/systemConsts';

const SECRET_DETAILS_PATHS = {
    SECRET_DETAILS: "",
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

    return (
        <TabbedPage
            basePath={`${pathname.substring(0, pathname.indexOf(id))}${id}`}
            items={[
                {
                    id: "general",
                    title: "Secret details",
                    isIndex: true,
                    component: () => <TabSecretDetails data={data} />
                }
            ]}
            withInnerPadding={false}
        />
    )
}

const SecretDetails = () => (
    <FindingsDetailsPage
        backTitle="Secrets"
        getTitleData={({findingInfo}) => ({title: findingInfo.fingerprint})}
        detailsContent={DetailsContent}
    />
)

export default SecretDetails;
