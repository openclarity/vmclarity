import React, { useEffect, useCallback } from 'react';
import { isUndefined } from 'lodash';
import TablePage from 'components/TablePage';
import { OPERATORS } from 'components/Filter';
import ToggleButton from 'components/ToggleButton';
import InfoIcon from 'components/InfoIcon';
import Loader from 'components/Loader';
import { getAssetColumnsFiltersConfig, scanColumnsFiltersConfig } from 'utils/utils';
import { APIS } from 'utils/systemConsts';
import { useFilterDispatch, useFilterState, setFilters, FILTER_TYPES } from 'context/FiltersProvider';

const FindingsTablePage = ({tableTitle, findingsObjectType, columns, filterType, filtersConfig}) => {
    const filtersDispatch = useFilterDispatch();
    const filtersState = useFilterState();

    const {customFilters} = filtersState[filterType];
    const {hideHistory} = customFilters;

    const setHideHistory = useCallback(hideHistory => setFilters(filtersDispatch, {
        type: filterType,
        filters: {hideHistory},
        isCustom: true
    }), [filterType, filtersDispatch]);
    
    useEffect(() => {
        if (isUndefined(hideHistory)) {
            setHideHistory(true);
        }
    }, [hideHistory, setHideHistory]);

    if (isUndefined(hideHistory)) {
        return <Loader />;
    }
    
    return (
        <div style={{position: "relative"}}>
            <div style={{position: "absolute", top: 0, right: 0, zIndex: 1, display: "flex", alignItems: "center"}}>
                <ToggleButton title="Hide history" checked={hideHistory} onChange={setHideHistory} />
                <div style={{marginLeft: "5px"}}>
                    <InfoIcon tooltipId="hide-hostory-info-icon" tooltipText="Hide findings that were replaced by a newer asset scan results of that type" />
                </div>
            </div>
            <TablePage
                columns={columns}
                url={APIS.FINDINGS}
                tableTitle={tableTitle}
                filterType={filterType}
                filtersConfig={[
                    ...filtersConfig,
                    ...getAssetColumnsFiltersConfig({prefix: "asset.targetInfo", withType: false, withLabels: false}),
                    ...scanColumnsFiltersConfig,
                    {value: "foundOn", label: "Found on", isDate: true, operators: [
                        {...OPERATORS.ge},
                        {...OPERATORS.le},
                    ]}
                ]}
                systemFilterType={FILTER_TYPES.FINDINGS_GENERAL}
                filters={[
                    `(findingInfo.objectType eq '${findingsObjectType}')`,
                    ...(hideHistory ? ["(invalidatedOn eq null)"] : [])
                ].join(` and `)}
                expand="asset,scan"
                defaultSortBy={{sortIds: ["scan.startTime"], desc: true}}
            />
        </div>
    )
}

export default FindingsTablePage;