import React from 'react';
import DoublePaneDisplay from 'components/DoublePaneDisplay';
import Title from 'components/Title';
import { ICON_NAMES } from 'components/Icon';
import LinksList from 'components/LinksList';
import VulnerabilitiesDisplay from '../VulnerabilitiesDisplay';
import FindingsCounterDisplay from './FindingsCounterDisplay';

const FINDINGS_ITEMS = [
    {title: "Exploits", key: "test", icon: ICON_NAMES.LOCK, path: "test"},
    {title: "Misconfigurations", key: "test", icon: ICON_NAMES.COG, path: "test"},
    {title: "Secrets", key: "test", icon: ICON_NAMES.KEY, path: "test"},
    {title: "Malwares", key: "test", icon: ICON_NAMES.BUG, path: "test"},
    {title: "Rootkits", key: "test", icon: ICON_NAMES.GHOST, path: "test"}
]

const TabFindings = ({data}) => {
    const {id} = data || {};
    
    return (
        <DoublePaneDisplay
            leftPaneDisplay={() => (
                <>
                    <Title medium>Findings</Title>
                    <LinksList
                        items={[
                            {path: "test", component: () => <VulnerabilitiesDisplay />},
                            ...FINDINGS_ITEMS.map(({title, key, icon, path}) => ({path: path, component: () => (
                                <FindingsCounterDisplay key={title} icon={icon} count={10} title={title} />
                            )}))
                        ]}
                    />
                </>
            )}
        />
    )
}

export default TabFindings;