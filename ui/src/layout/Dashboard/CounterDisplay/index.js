import React, {useMemo} from 'react';
import { useFetch } from 'hooks';
import { formatNumber } from 'utils/utils';

import COLORS from "../../../utils/scss_variables.module.scss";
import { APIS } from "../../../utils/systemConsts";

import './counter-display.scss';

export const ScanCounterDisplay = () => {
    const [{data, error, loading}] = useFetch(APIS.SCANS, {queryParams: {"$select": "state"}});

    const completedScans = useMemo(() => {
        let ret = 0;

        data?.items.forEach((element) => {
            if (element.state === "Aborted" || element.state === "Failed" || element.state === "Done") {
                ret++
            }
        })

        return ret
    }, [data])

    return (
        <div className="dashboard-counter" style={{background: COLORS["color-gradient-green"]}}>
            {loading || error ? "" :
                <div className="dashboard-counter-content">
                    <div className="dashboard-counter-count">{formatNumber(completedScans)}</div> Completed scans
                </div>
            }
        </div>
    )
}

export const CounterDisplay = ({url, title, background}) => {
    const [{data, error, loading}] = useFetch(url, {queryParams: {"$count": true, "$top": 1, "$select": "id"}});

    return (
        <div className="dashboard-counter" style={{background}}>
            {loading || error ? "" :
                <div className="dashboard-counter-content"><div className="dashboard-counter-count">{formatNumber(data.count)}</div>{` ${title}`}</div>
            }
        </div>
    )
}
