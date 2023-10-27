import React from 'react';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import { TagsList } from 'components/Tag';
import { formatTagsToStringsList } from 'utils/utils';
import prettyBytes from 'pretty-bytes';


export const ContainerInfoDetails = ({assetData}) => {
    const {containerName, containerID, image} = assetData.assetInfo || {};
    const {repoTags, labels, architecture, os, size} = image || {};


    return (
        <>
            <TitleValueDisplayRow>
                <TitleValueDisplay title="Container Name">{containerName}</TitleValueDisplay>
                <TitleValueDisplay title="Container ID">{containerID}</TitleValueDisplay>
            </TitleValueDisplayRow>

            <TitleValueDisplayRow>
                <TitleValueDisplay title="Architecture">{architecture}</TitleValueDisplay>
                <TitleValueDisplay title="OS">{os}</TitleValueDisplay>
                <TitleValueDisplay title="Size">{prettyBytes(size)}</TitleValueDisplay>
            </TitleValueDisplayRow>

            <TitleValueDisplayRow>
                <TitleValueDisplay title="Labels"><TagsList items={formatTagsToStringsList(labels)} /></TitleValueDisplay>
            </TitleValueDisplayRow>

            <TitleValueDisplayRow>
                <TitleValueDisplay title="Repo tags">{repoTags?.[0]}</TitleValueDisplay>
            </TitleValueDisplayRow>
        </>
    )
}
