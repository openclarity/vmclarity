import React, { useMemo, useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { isUndefined } from 'lodash';
import ExpandableList from 'components/ExpandableList';
import ToggleButton from 'components/ToggleButton';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import Loader from 'components/Loader';
import { getFindingsColumnsConfigList, getVulnerabilitiesColumnConfigItem, formatTagsToStringsList, formatDate, getAssetName} from 'utils/utils';
import { APIS } from 'utils/systemConsts';
import { useFilterDispatch, useFilterState, setFilters, FILTER_TYPES } from 'context/FiltersProvider';

const TABLE_TITLE = "assets";

const NAME_SORT_IDS = ["asset.assetInfo.instanceID", "asset.assetInfo.podName", "asset.assetInfo.dirName", "asset.assetInfo.imageID", "asset.assetInfo.containerName"];
const LABEL_SORT_IDS = ["asset.assetInfo.tags", "asset.assetInfo.labels"];
const LOCATION_SORT_IDS = ["asset.assetInfo.location"];

const ASSETS_FILTER_TYPE = FILTER_TYPES.ASSETS;

const AssetsForFindingTable = (id) => {
    const {findingId} = id;

    const navigate = useNavigate();
    const filtersDispatch = useFilterDispatch();
    const filtersState = useFilterState();

    const {customFilters} = filtersState[ASSETS_FILTER_TYPE];
    const {hideTerminated} = customFilters;

    const setHideTerminated = useCallback(hideTerminated => setFilters(filtersDispatch, {
        type: ASSETS_FILTER_TYPE,
        filters: {hideTerminated},
        isCustom: true
    }), [filtersDispatch]);

    useEffect(() => {
        if (isUndefined(hideTerminated)) {
            setHideTerminated(true);
        }
    }, [hideTerminated, setHideTerminated]);

    const [showFindingCounts, setShowFindingCounts] = useState(false);

    const columns = useMemo(() => [
        {
            Header: "Name",
            id: "instanceID",
            sortIds: NAME_SORT_IDS,
            accessor: (original) => getAssetName(original.asset.assetInfo),
        },
        {
            Header: "Labels",
            id: "tags",
            sortIds: LABEL_SORT_IDS,
            Cell: ({row}) => {
                const {tags, labels} = row.original.asset.assetInfo;

                return (
                    <ExpandableList items={formatTagsToStringsList(tags ?? labels)} withTagWrap />
                )
            },
            alignToTop: true
        },
        {
            Header: "Type",
            id: "objectType",
            sortIds: ["asset.assetInfo.objectType"],
            accessor: "asset.assetInfo.objectType"
        },
        {
            Header: "Location",
            id: "location",
            sortIds: LOCATION_SORT_IDS,
            accessor: (original) => original.asset.assetInfo.location || original.asset.assetInfo.repoDigests?.[0],
        },
        {
            Header: "Last Seen",
            id: "lastSeen",
            sortIds: ["lastSeen"],
            accessor: original => formatDate(original.asset.lastSeen)
        },
        ...(hideTerminated ? [] : [{
            Header: "Terminated On",
            id: "terminatedOn",
            sortIds: ["terminatedOn"],
            accessor: original => formatDate(original.asset?.terminatedOn)
        }]),
        ...(!showFindingCounts ? [] : [
            getVulnerabilitiesColumnConfigItem({tableTitle: TABLE_TITLE, withAssetPrefix: true}),
            ...getFindingsColumnsConfigList({tableTitle: TABLE_TITLE, withAssetPrefix: true})
        ]),
    ], [hideTerminated, showFindingCounts]);

    if (isUndefined(hideTerminated) || isUndefined(showFindingCounts)) {
        return <Loader />;
    }

    let filtersList = [`(finding/id eq '${findingId}')`]
    if (hideTerminated) {
        filtersList.push("(asset/terminatedOn eq null)");
    }
    let select = "asset/id,asset/assetInfo,asset/lastSeen"
    if (!hideTerminated) {
        select += ",asset/terminatedOn"
    }
    if (showFindingCounts) {
        select += ",asset/summary"
    }
    const expand = "asset"

    return (
        <div style={{marginTop: "20px"}}>
            <div style={{float: "right", marginRight: "36px"}}>
                <ToggleButton title="Show finding counts" checked={showFindingCounts} onChange={setShowFindingCounts}/>
            </div>
            <div style={{float: "right", marginRight: "36px"}}>
                <ToggleButton title="Hide terminated" checked={hideTerminated} onChange={setHideTerminated}/>
            </div>
            <div style={{position: "relative"}}>
                <ContentContainer withMargin>
                    <Table
                        paginationItemsName={TABLE_TITLE.toLowerCase()}
                        filters={{
                            ...(!!expand ? {"$expand": expand} : {}),
                            ...(!!select ? {"$select": select} : {}),
                            ...(filtersList.length > 0 ? {"$filter": filtersList.join(" and ")} : {})
                        }}
                        noResultsTitle={TABLE_TITLE}
                        onLineClick={({asset}) => navigate(`/${APIS.ASSETS}/${asset.id}`)}
                        columns={columns}
                        url={APIS.ASSET_FINDINGS}
                        defaultSortBy={{sortIds: ["asset/lastSeen", "asset/terminatedOn"], desc: true}}
                        defaultPageSize={10}
                    />
                </ContentContainer>
            </div>
        </div>
    )
}

export default AssetsForFindingTable;
