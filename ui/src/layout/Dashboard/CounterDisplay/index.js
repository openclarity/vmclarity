import React, {useMemo} from 'react';
import { useFetch } from 'hooks';
import { formatNumber } from 'utils/utils';

import './counter-display.scss';

const CounterDisplay = ({url, title, background}) => {
    const [{data, error, loading}] = useFetch(url, {queryParams: {"$select": "state"}});

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
        <div className="dashboard-counter" style={{background}}>
            {loading || error ? "" :
                <div className="dashboard-counter-content"><div className="dashboard-counter-count">{formatNumber(completedScans)}</div>{` ${title}`}</div>
            }
        </div>
    )
}

export default CounterDisplay;
