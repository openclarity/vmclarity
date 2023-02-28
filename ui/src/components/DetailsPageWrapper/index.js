import React from 'react';
import classnames from 'classnames';
import { useLocation, useParams } from 'react-router-dom';
import BackRouteButton from 'components/BackRouteButton';
import ContentContainer from 'components/ContentContainer';
import Loader from 'components/Loader';
import { useFetch } from 'hooks';

import './details-page-wrapper.scss';

const DetailsContentWrapper = ({data, titleKey, detailsContent: DetailsContent}) => (
    <div className="details-page-content-wrapper">
        <div className="details-page-title">{data[titleKey]}</div>
        <ContentContainer><DetailsContent data={data} /></ContentContainer>
    </div>
)

const DetailsPageWrapper = ({className, backTitle, url, getUrl, getReplace, titleKey, detailsContent}) => {
    const {pathname} = useLocation();
    const params = useParams();
    const {id} = params;
    
    const [{loading, data, error}, fetchData] = useFetch(!!url ? `${url}/${id}` : getUrl(params));

    return (
        <div className={classnames("details-page-wrapper", className)}>
            <BackRouteButton title={backTitle} pathname={pathname.replace(!!getReplace ? getReplace(params) : `/${id}`, "")} />
            {loading ? <Loader /> : (!!error ? null : <DetailsContentWrapper detailsContent={detailsContent} titleKey={titleKey} data={data} />)}
        </div>
    )
}

export default DetailsPageWrapper;