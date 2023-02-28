import React from 'react';
import TitleValueDisplay, { TitleValueDisplayColumn, ValuesListDisplay } from 'components/TitleValueDisplay';
import Tag from 'components/Tag';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import { formatRegionsToStrings, formatTagsToStringInstances, getEnabledScanTypesList } from 'layout/Scans/utils';

const ScopeDisplay = ({all, regions}) => {
    if (all) {
        return "All";
    }

    return (
        <ValuesListDisplay values={formatRegionsToStrings(regions)} />
    )
}

const ScanTypesDisplay = ({scanFamiliesConfig}) => (
    <ValuesListDisplay values={getEnabledScanTypesList(scanFamiliesConfig)} />
)

const InstancesDisplay = ({tags}) => (
    <div className="configuration-instances-tags-display">
        {
            formatTagsToStringInstances(tags).map(tag => <Tag key={tag}>{tag}</Tag>)
        }
    </div>
)

const TabConfiguration = ({data}) => {
    const {scope, scanFamiliesConfig} = data || {};
    const {all, regions, instanceTagSelector, instanceTagExclusion} = scope;
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <TitleValueDisplayColumn>
                    <TitleValueDisplay title="Scope"><ScopeDisplay all={all} regions={regions} /></TitleValueDisplay>
                    <TitleValueDisplay title="Included instances"><InstancesDisplay tags={instanceTagSelector}/></TitleValueDisplay>
                    <TitleValueDisplay title="Excluded instances"><InstancesDisplay tags={instanceTagExclusion}/></TitleValueDisplay>
                </TitleValueDisplayColumn>
            )}
            rightPlaneDisplay={() => (
                <TitleValueDisplayColumn>
                    <TitleValueDisplay title="Scan types"><ScanTypesDisplay scanFamiliesConfig={scanFamiliesConfig} /></TitleValueDisplay>
                </TitleValueDisplayColumn>
            )}
        />
    )
}

export default TabConfiguration;