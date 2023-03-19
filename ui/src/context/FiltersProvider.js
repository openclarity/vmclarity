import { create } from './utils';

export const FILTER_TYPES = {
    ASSETS: "ASSETS",
    ASSET_SCANS: "ASSET_SCANS",
    PACKAGES: "PACKAGES",
    SCANS: "SCANS",
    SCAN_CONFIGURATIONS: "SCAN_CONFIGURATIONS",
    FINDINGS: "FINDINGS"
}

const initialState = Object.keys(FILTER_TYPES).reduce((acc, curr) => ({
    ...acc,
    [curr]: {
        tableFilters: [],
        systemFilters: {}
    }
}), {});

const FITLER_ACTIONS = {
    SET_TABLE_FILTERS_BY_KEY: "SET_TABLE_FILTERS_BY_KEY",
    SET_SYSTEM_FILTERS_BY_KEY: "SET_SYSTEM_FILTERS_BY_KEY",
    RESET_ALL_FILTERS: "RESET_ALL_FILTERS",
    RESET_FILTERS_BY_KEY: "RESET_FILTERS_BY_KEY"
}

const reducer = (state, action) => {
    switch (action.type) {
        case FITLER_ACTIONS.SET_TABLE_FILTERS_BY_KEY: {
            const {filterType, filterData} = action.payload;

            return {
                ...state,
                [filterType]: {
                    ...state[filterType],
                    tableFilters: filterData
                }
            };
        }
        case FITLER_ACTIONS.SET_SYSTEM_FILTERS_BY_KEY: {
            const {filterType, filterData} = action.payload;
            
            return {
                ...state,
                [filterType]: {
                    ...state[filterType],
                    tableFilters: [...initialState[filterType].tableFilters],
                    systemFilters: filterData
                }
            };
        }
        case FITLER_ACTIONS.RESET_ALL_FILTERS: {
            return {
                ...state,
                ...initialState
            };
        }
        case FITLER_ACTIONS.RESET_FILTERS_BY_KEY: {
            const {filterTypes} = action.payload;
            
            return {
                ...state,
                ...filterTypes.reduce((acc, curr) => ({
                    ...acc,
                    [curr]: {...initialState[curr]}
                }), {})
            };
        }
        default:
            return state;
    }
}

const [FiltersProvider, useFilterState, useFilterDispatch] = create(reducer, initialState);

const setFilters = (dispatch, {type, filters, isSystem=false}) => dispatch({
    type: isSystem ? FITLER_ACTIONS.SET_SYSTEM_FILTERS_BY_KEY : FITLER_ACTIONS.SET_TABLE_FILTERS_BY_KEY,
    payload: {filterType: type, filterData: filters}
});
const resetAllFilters = (dispatch) => dispatch({type: FITLER_ACTIONS.RESET_ALL_FILTERS});
const resetFilters = (dispatch, filterTypes) => dispatch({type: FITLER_ACTIONS.RESET_FILTERS_BY_KEY, payload: {filterTypes}});
const resetSystemFilters = (dispatch, type) => setFilters(dispatch, {type, filters: {}, isSystem: true})

export {
    FiltersProvider,
    useFilterState,
    useFilterDispatch,
    setFilters,
    resetAllFilters,
    resetFilters,
    resetSystemFilters
};