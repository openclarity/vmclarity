import React from 'react';
import { useFetch } from 'hooks';
import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import Loader from 'components/Loader';
import { APIS } from 'utils/systemConsts';
import { formatNumber } from 'utils/utils';

const AssetCountDisplay = (findingId) => {
    const filter = `finding/id eq '${findingId}'`;
    const [{loading, data, error}] = useFetch(APIS.ASSET_FINDINGS, {
        queryParams: {"$filter": filter, "$count": true, "$select": "count"}
    });

    if (error) {
        return null;
    }

    if (loading) {
        return <Loader absolute={false} small />
    }

    return (
        <TitleValueDisplayRow>
            <TitleValueDisplay title="Asset count">{formatNumber(data?.count || 0)}</TitleValueDisplay>
        </TitleValueDisplayRow>
    )
}

export default AssetCountDisplay;
