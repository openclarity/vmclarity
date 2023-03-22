import React from 'react';
import TitleValueDisplay, { ValuesListDisplay } from 'components/TitleValueDisplay';
import { TagsList } from 'components/Tag';
import ExpandableList from 'components/ExpandableList';
import { formatRegionsToStrings, formatTagsToStringInstances, getEnabledScanTypesList } from 'layout/Scans/utils';

const SCOPE_ALL = "All";

export const ExpandableScopeDisplay = ({all, regions}) => (
    all ? SCOPE_ALL : <ExpandableList items={formatRegionsToStrings(regions)} />
)

const ScopeDisplay = ({all, regions}) => {
    if (all) {
        return SCOPE_ALL;
    }

    return ( 
        <ValuesListDisplay values={formatRegionsToStrings(regions)} />
    )
}

const InstancesDisplay = ({tags}) => (
    <TagsList items={formatTagsToStringInstances(tags)} />
)

export const ConfigurationReadOnlyDisplay = ({scope, scanFamiliesConfig}) => {
    const {allRegions, regions, instanceTagSelector, instanceTagExclusion, shouldScanStoppedInstances} = scope;

    return (
        <>
            <TitleValueDisplay title="Scope"><ScopeDisplay all={allRegions} regions={regions} /></TitleValueDisplay>
            <TitleValueDisplay title="Instances">
                <div style={{margin: "10px 0 20px 0"}}>
                    {shouldScanStoppedInstances ? "Running and non-running instances" : "Running instances only"}
                </div>
                <TitleValueDisplay title="Included instances" isSubItem><InstancesDisplay tags={instanceTagSelector}/></TitleValueDisplay>
                <TitleValueDisplay title="Excluded instances" isSubItem><InstancesDisplay tags={instanceTagExclusion}/></TitleValueDisplay>
            </TitleValueDisplay>
            <TitleValueDisplay title="Scan types"><ValuesListDisplay values={getEnabledScanTypesList(scanFamiliesConfig)} /></TitleValueDisplay>
        </>
    )
}