import React, { useEffect, useMemo, useState } from 'react';
import { useField } from 'formik';
import { Utils as QbUtils } from '@react-awesome-query-builder/mui';

import { FieldLabel, TextAreaField, TextField, validators } from 'components/Form'; // useFormikContext,
import { StatelessCheckbox, StatelessRadioButtonGroup } from 'components/statelessComponents';
import { EMPTY_JSON_SCOPE_TREE, SCOPE_CONFIG_KEY, SCOPE_TREE_KEY } from "../ScanConfigWizardModal.constants";
import { CustomQueryBuilder } from './CustomQueryBuilder';

const SCOPES = [
    { label: "All", value: true },
    { label: "Define scope", value: false }
];

const StepGeneralProperties = ({
    configWithFields,
    error,
    isEditForm,
    isQueryBuilderVisible,
    loading,
    queryState,
    setIsQueryBuilderVisible,
    setQueryState
}) => {
    const [fullScope, setFullScope] = useState(SCOPES[0].value);
    const [isManualScope, setIsManualScope] = useState(false);
    const [showHumanFriendlyScope, setShowHumanFriendlyScope] = useState(false);
    const [scopeField, , scopeHelpers] = useField("scanTemplate.scope");
    const { setValue: setScopeValue } = scopeHelpers;
    const { value: scopeValue } = scopeField;

    const [annotationsField, ,] = useField("annotations");
    const { value: annotations } = annotationsField;

    const savedScopeTree = useMemo(
        () => {
            const treeObject = (annotations?.length >= 0 ? annotations : []).find(f => Object.keys(f).includes(SCOPE_TREE_KEY));
            const valueString = treeObject[SCOPE_TREE_KEY];
            if (valueString) return JSON.parse(valueString);
            return undefined
        },
        [annotations]
    );

    const savedScopeConfig = useMemo(
        () => {
            const configObject = (annotations?.length >= 0 ? annotations : []).find(f => Object.keys(f).includes(SCOPE_CONFIG_KEY));
            const valueString = configObject[SCOPE_CONFIG_KEY];
            if (valueString) return JSON.parse(valueString);
            return undefined
        },
        [annotations]
    );

    const initialConfig = useMemo(
        () => (isEditForm && savedScopeConfig) ? QbUtils.decompressConfig(savedScopeConfig, configWithFields) : configWithFields,
        [configWithFields, isEditForm, savedScopeConfig]
    );

    const initialJsonTree = useMemo(
        () => (savedScopeTree && Object.keys(savedScopeTree).length > 0) ? savedScopeTree : EMPTY_JSON_SCOPE_TREE,
        [savedScopeTree]);

    const initialTree = useMemo(() => QbUtils.checkTree(QbUtils.loadTree(initialJsonTree), initialConfig), [initialConfig, initialJsonTree]);

    useEffect(() => {
        setQueryState({
            tree: initialTree,
            config: initialConfig
        })
        // eslint-disable-next-line
    }, [initialConfig, initialTree])

    useEffect(() => {
        if (!fullScope) {
            setIsQueryBuilderVisible(true);
            setIsManualScope(false)
        } else {
            setIsQueryBuilderVisible(false);
            setScopeValue("");
        }
        // eslint-disable-next-line
    }, [fullScope])

    useEffect(() => {
        if (!fullScope) {
            setIsQueryBuilderVisible(!isManualScope);
        }
        // eslint-disable-next-line
    }, [isManualScope])

    useEffect(() => {
        if (savedScopeConfig && savedScopeTree) {
            setFullScope(false);
            setIsManualScope(false);
            setIsQueryBuilderVisible(true);
        }
        // eslint-disable-next-line
    }, [savedScopeConfig, savedScopeTree])

    return (
        <div className="scan-config-general-step">
            <div className="manual-query-container">
                <TextField
                    label="Scan config name*"
                    name="name"
                    placeholder="Type a scan config name..."
                    validate={validators.validateRequired}
                />
                <StatelessRadioButtonGroup
                    initialValue={SCOPES[0].value}
                    items={SCOPES}
                    label="Scope*"
                    name="fullScope"
                    setValue={setFullScope}
                    value={fullScope}
                    tooltipText="You can narrow the scope of scanning here"
                />
                {!fullScope &&
                    <>
                        <StatelessCheckbox
                            className="checkbox"
                            title={`${scopeValue ? "Edit" : "Create"} scope only manually`}
                            value={isManualScope}
                            setValue={setIsManualScope}
                        />
                        <div className='query-builder-result__section'>
                            {showHumanFriendlyScope
                                ?
                                <>
                                    <div className='query-builder-result__odata'>
                                        {QbUtils.queryString(queryState.tree, queryState.config, true) ?? "-"}
                                    </div>
                                </>
                                :
                                <>
                                    <FieldLabel>Manual scope editor (odata query)*</FieldLabel>
                                    <div className='query-builder-result__odata'>
                                        <span className='query-builder-result__odata--details'>(This query is going to be used by the scanner)</span>
                                        <TextAreaField
                                            name="scanTemplate.scope"
                                            placeholder="You can type a scope manually..."
                                        />
                                    </div>
                                </>}
                        </div>
                        <div className='query-builder-result__section'>
                            {/* <span className='query-builder-result__title'>Human friendly scope:{" "}</span> */}
                            <StatelessCheckbox
                                className="checkbox"
                                title={`${showHumanFriendlyScope ? "Hide" : "Show"} human friendly scope`}
                                value={showHumanFriendlyScope}
                                setValue={setShowHumanFriendlyScope}
                            />
                        </div>
                    </>
                }
            </div>
            {
                isQueryBuilderVisible &&
                <div className='query-builder-wrapper'>
                    <FieldLabel>Scope builder</FieldLabel>
                    <CustomQueryBuilder
                        errorMessage={error?.errorMessage}
                        initialConfig={initialConfig}
                        initialTree={initialTree}
                        loading={loading}
                        queryState={queryState}
                        setQueryState={setQueryState}
                    />
                </div>
            }
        </div>
    )
}

export default StepGeneralProperties;
