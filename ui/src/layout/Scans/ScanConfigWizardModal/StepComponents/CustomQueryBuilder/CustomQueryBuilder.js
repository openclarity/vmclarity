import React, { useCallback, useEffect } from 'react';
import classNames from "classnames";
import throttle from "lodash/throttle";
import { Query, Builder, Utils as QbUtils } from '@react-awesome-query-builder/mui';
import { useField } from 'formik';

import Button from 'components/Button';
import FieldError from 'components/Form/FieldError';
import Loader from 'components/Loader';
import { postFixQuery, updateSavedScopeConfig, updateSavedScopeTree } from './CustomQueryBuilder.functions';
import { EMPTY_JSON_SCOPE_TREE } from '../../ScanConfigWizardModal.constants';

import "@react-awesome-query-builder/ui/css/styles.scss";
import "./CustomQueryBuilder.scss";

const CustomQueryBuilder = ({
    errorMessage,
    initialConfig,
    initialTree,
    loading,
    queryState,
    setQueryState,
}) => {
    const [annotationsField, , annotationsHelpers] = useField("annotations");
    const { setValue: setAnnotationsValue } = annotationsHelpers;
    const { value: annotationsValue } = annotationsField;

    const [, , scopeHelpers] = useField("scanTemplate.scope");
    const { setValue: setScopeValue } = scopeHelpers;

    const resetValue = useCallback(() => {
        setQueryState(state => ({
            ...state,
            tree: QbUtils.checkTree(QbUtils.loadTree(initialTree), initialConfig),
        }));
    }, [initialConfig, initialTree, setQueryState]);

    const clearValue = useCallback(() => {
        setQueryState(state => ({
            ...state,
            tree: QbUtils.loadTree(EMPTY_JSON_SCOPE_TREE),
        }));
    }, [setQueryState]);

    const renderBuilder = useCallback((props) => (
        <div className="query-builder-container">
            <div className={classNames("query-builder", { "qb-lite": queryState.tree.size > 2 })}>
                <Builder {...props} />
            </div>
        </div>
    ), [queryState.tree.size]);

    const updateResult = useCallback((immutableTree, config) => {
        throttle(() => {
            setQueryState(prevState => ({ ...prevState, tree: immutableTree, config }));
            // eslint-disable-next-line
        }, 100);
    }, [setQueryState]);

    const onChange = useCallback((immutableTree, config) => {
        updateResult(immutableTree, config);
        const jsonTree = QbUtils.getTree(immutableTree);
        const annotationsWithTree = updateSavedScopeTree(annotationsValue, jsonTree);
        const annotationsWithTreeAndConfig = updateSavedScopeConfig(annotationsWithTree, config);
        setAnnotationsValue(annotationsWithTreeAndConfig);
    }, [annotationsValue, setAnnotationsValue, updateResult]);

    useEffect(() => {
        const query = QbUtils.queryString(queryState.tree, queryState.config);
        setScopeValue(postFixQuery(query));
        // eslint-disable-next-line
    }, [queryState.tree])

    return (
        <>
            <div className="query-builder-result">
                <div className="query-buttons">
                    <Button onClick={resetValue}>Reset</Button>
                    <Button className="query-buttons__clear-button" onClick={clearValue}>Clear</Button>
                </div>
            </div>
            {loading && <Loader absolute={false} />}
            {errorMessage && <FieldError>{errorMessage}</FieldError>}
            {Object.keys(queryState.config.fields).length > 0 &&
                <Query
                    {...queryState.config}
                    value={queryState.tree}
                    onChange={onChange}
                    renderBuilder={renderBuilder}
                />
            }
        </>
    )
}

export default CustomQueryBuilder;
