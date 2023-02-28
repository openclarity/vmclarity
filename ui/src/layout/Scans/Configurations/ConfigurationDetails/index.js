import React from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import DetailsPageWrapper from 'components/DetailsPageWrapper';
import TabbedPage from 'components/TabbedPage';
import ConfigurationActionsDisplay from '../ConfigurationActionsDisplay';
import TabConfiguration from './TabConfiguration';
import TabScans from './TabScans';

import './configuration-details.scss';

const getReplace = params => {
    const id = params["id"];
    const innerTab = params["*"];
    
    return !!innerTab ? `/${id}/${innerTab}` : `/${id}`;
}

const DetailsContent = ({data, setScanConfigFormData}) => {
    const navigate = useNavigate();
    const {pathname} = useLocation();
    const params = useParams();
    
    const {id} = data;
    
    return (
        <TabbedPage
            basePath={`${pathname.substring(0, pathname.indexOf(id))}${id}`}
            items={[
                {
                    id: "config",
                    title: "Configuration",
                    isIndex: true,
                    component: () => <TabConfiguration data={data} />
                },
                {
                    id: "scans",
                    title: "Scans",
                    path: "scans",
                    component: () => <TabScans data={data} />
                }
            ]}
            headerCustomDisplay={() => (
                <ConfigurationActionsDisplay
                    data={data}
                    setScanConfigFormData={setScanConfigFormData}
                    onDelete={() => navigate(pathname.replace(getReplace(params), ""))}
                />
            )}
            withInnerPadding={false}
        />
    )
}

const ConfigurationDetails = ({setScanConfigFormData}) => (
    <DetailsPageWrapper
        className="configuration-details-page-wrapper"
        backTitle="Scan configurations"
        url="scanConfigs"
        titleKey="name"
        detailsContent={props => <DetailsContent {...props} setScanConfigFormData={setScanConfigFormData} />}
        getReplace={params => getReplace(params)}
    />
)

export default ConfigurationDetails;