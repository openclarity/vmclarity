import React from 'react';
import ListAndDetailsRouter from 'components/ListAndDetailsRouter';
import ConfigurationsTable from './ConfigurationsTable';
import ConfigurationDetails from './ConfigurationDetails';

export const SCAN_CONFIGS_PATH = "configs";

const Configurations = ({setScanConfigFormData}) => (
    <ListAndDetailsRouter
        listComponent={() => <ConfigurationsTable setScanConfigFormData={setScanConfigFormData} />}
        detailsComponent={() => <ConfigurationDetails setScanConfigFormData={setScanConfigFormData} />}
    />
)


export default Configurations;