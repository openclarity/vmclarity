import React, { useState } from 'react';
import { isNull } from 'lodash';
import { useMountMultiFetch } from 'hooks';
import TabbedPage from 'components/TabbedPage';
import Loader from 'components/Loader';
import EmptyDisplay from 'components/EmptyDisplay';
import ScanConfigWizardModal from './ScanConfigWizardModal';
import ScansTable from './ScansTable';
import ConfigurationsTable, { SCAN_CONFIGS_PATH } from './ConfigurationsTable';

const Scans = React.memo(({setScanConfigFormData, data}) => {
    const {scans, scanConfigs} = data;
    
    return (
        <React.Fragment>
            {(!scans?.total && !scanConfigs?.total) ?
                <EmptyDisplay
                    message={(
                        <>
                            <div>No scans detected.</div>
                            <div>Create your first scan configuration to see your VM's issues.</div>
                        </>
                    )}
                    title="New scan configuration"
                    onClick={() => setScanConfigFormData({})}
                /> :
                <TabbedPage
                    items={[
                        {
                            id: "scans",
                            title: "Scans",
                            isIndex: true,
                            component: () => <ScansTable setScanConfigFormData={setScanConfigFormData} />
                        },
                        {
                            id: "configs",
                            title: "Configurations",
                            path: SCAN_CONFIGS_PATH,
                            component: () => <ConfigurationsTable setScanConfigFormData={setScanConfigFormData} />
                        }
                    ]}
                />
            }
        </React.Fragment>
    
    )
}, () => true);

const ScansWrapper = () => {
    const [scanConfigFormData, setScanConfigFormData] = useState(null);
    const closeConfigForm = () => setScanConfigFormData(null);

    const [{data, error, loading}, fetchData] = useMountMultiFetch([
        {key: "scans", url: "scans"},
        // {key: "scanConfigs", url: "scanConfigs"}
    ]);

    if (loading) {
        return <Loader />;
    }

    if (error) {
        return null;
    }

    return (
        <>
            <Scans setScanConfigFormData={setScanConfigFormData} data={data} />
            {!isNull(scanConfigFormData) && 
                <ScanConfigWizardModal
                    initialData={scanConfigFormData}
                    onClose={closeConfigForm}
                    onSubmitSuccess={() => {
                        closeConfigForm();
                        fetchData();
                    }}
                />
            }
        </>
    )
}

export default ScansWrapper;


