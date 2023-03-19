import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import ContentContainer from 'components/ContentContainer';
import Table from 'components/Table';
import SystemFilterBanner from 'components/SystemFiltersBanner';
import { toCapitalized, BoldText } from 'utils/utils';
import { useFilterState, useFilterDispatch, resetSystemFilters } from 'context/FiltersProvider';

const TablePage = ({tableTitle, filterType, filters, expand, withMargin, absoluteSystemBanner=false, ...tableProps}) => {
    const navigate = useNavigate();
    const {pathname} = useLocation();

    const filtersState = useFilterState();
    const {systemFilters} = filtersState[filterType];
    const filtersDispatch = useFilterDispatch();
    
    const {name: systemFilterName, suffix: systemSuffix, backPath: systemFilterBackPath, filter: systemFilter, customDisplay} = systemFilters;

    const onSystemFilterClose = () => resetSystemFilters(filtersDispatch, filterType);
    
    const fitlersList = [
        ...(!!filters ? [filters] : []),
        ...(!!systemFilter ? [systemFilter]  : [])
    ]
    
    return (
        <>
            {!!systemFilterName &&
                <SystemFilterBanner
                    displayText={<span>{`${toCapitalized(tableTitle)} for `}<BoldText>{systemFilterName}</BoldText>{` ${systemSuffix}`}</span>}
                    onClose={onSystemFilterClose}
                    backPath={systemFilterBackPath}
                    absolute={absoluteSystemBanner}
                    customDisplay={customDisplay}
                />
            }
            <ContentContainer withMargin={withMargin}>
                <Table
                    paginationItemsName={tableTitle.toLowerCase()}
                    filters={{
                        ...(!!expand ? {"$expand": expand} : {}),
                        ...(fitlersList.length > 0 ? {"$filter": fitlersList.join(" and ")} : {})
                    }}
                    noResultsTitle={tableTitle}
                    onLineClick={({id}) => navigate(`${pathname}/${id}`)}
                    {...tableProps}
                />
            </ContentContainer>
        </>
    )
}

export default TablePage;