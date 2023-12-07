import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { FETCH_METHODS, useFetch } from 'hooks';
import { Utils as QbUtils } from '@react-awesome-query-builder/core';
import OpenAPIParser from '@readme/openapi-parser';

import {
    CRON_QUICK_OPTIONS,
    SCHEDULE_TYPES_ITEMS,
    StepAdvancedSettings,
    StepGeneralProperties,
    StepScanTypes,
    StepTimeConfiguration,
} from './StepComponents';
import WizardModal from 'components/WizardModal';
import { APIS } from 'utils/systemConsts';
import {
    BASIC_CONFIG,
    collectProperties,
} from './StepComponents/CustomQueryBuilder';
import { EMPTY_JSON_SCOPE_TREE } from "./ScanConfigWizardModal.constants";

import './scan-config-wizard-modal.scss';

const padDateTime = time => String(time).padStart(2, "0");

const ScanConfigWizardModal = ({ initialData, onClose, onSubmitSuccess }) => {
    const { id, name, scanTemplate, scheduled } = initialData || {};
    const { scope, maxParallelScanners, assetScanTemplate } = scanTemplate || {};
    const { operationTime, cronLine } = scheduled || {};

    const { scanFamiliesConfig, scannerInstanceCreationConfig } = assetScanTemplate || {}
    const { useSpotInstances } = scannerInstanceCreationConfig || {};

    const [{ loading, data, error }] = useFetch(`${window.location.origin}/api/openapi.json`, { isAbsoluteUrl: true });

    const [isQueryBuilderVisible, setIsQueryBuilderVisible] = useState(true);

    const isEditForm = useMemo(() => !!id, [id]);

    const [configWithFields, setConfigWithFields] = useState(BASIC_CONFIG);

    const [queryState, setQueryState] = useState({
        config: configWithFields,
        tree: QbUtils.checkTree(QbUtils.loadTree(EMPTY_JSON_SCOPE_TREE), configWithFields),
    });

    const INITIAL_SCAN_CONFIG_FORM_VALUES = useMemo(() => ({
        annotations: [],
        id: id || null,
        name: name || "",
        scanFamiliesConfig: {
            sbom: { enabled: true },
            vulnerabilities: { enabled: true },
            malware: { enabled: false },
            rootkits: { enabled: false },
            secrets: { enabled: false },
            misconfigurations: { enabled: false },
            infoFinder: { enabled: false },
            exploits: { enabled: false }
        },
        scanTemplate: {
            scope: scope || "",
            maxParallelScanners: maxParallelScanners || 2,
            assetScanTemplate: {
                scanFamiliesConfig: {
                    sbom: { enabled: true },
                    vulnerabilities: { enabled: true },
                    malware: { enabled: false },
                    rootkits: { enabled: false },
                    secrets: { enabled: false },
                    misconfigurations: { enabled: false },
                    infoFinder: { enabled: false },
                    exploits: { enabled: false }
                },
                scannerInstanceCreationConfig: {
                    useSpotInstances: useSpotInstances || false
                }
            }
        },
        scheduled: {
            scheduledSelect: !!cronLine ? SCHEDULE_TYPES_ITEMS.REPETITIVE.value : SCHEDULE_TYPES_ITEMS.NOW.value,
            laterDate: "",
            laterTime: "",
            cronLine: cronLine || CRON_QUICK_OPTIONS[0].value
        }
    }), [cronLine, id, maxParallelScanners, name, scope, useSpotInstances]);

    if (!!operationTime && !cronLine) {
        const dateTime = new Date(operationTime);
        INITIAL_SCAN_CONFIG_FORM_VALUES.scheduled.scheduledSelect = SCHEDULE_TYPES_ITEMS.LATER.value;
        INITIAL_SCAN_CONFIG_FORM_VALUES.scheduled.laterTime = `${padDateTime(dateTime.getHours())}:${padDateTime(dateTime.getMinutes())}`;
        INITIAL_SCAN_CONFIG_FORM_VALUES.scheduled.laterDate = `${dateTime.getFullYear()}-${padDateTime(dateTime.getMonth() + 1)}-${padDateTime(dateTime.getDate())}`;
    }

    Object.keys(scanFamiliesConfig || {}).forEach(type => {
        const { enabled } = scanFamiliesConfig[type];
        INITIAL_SCAN_CONFIG_FORM_VALUES.scanTemplate.assetScanTemplate.scanFamiliesConfig[type].enabled = enabled;
    })

    const steps = [
        {
            id: "general",
            title: "General properties",
            component: StepGeneralProperties,
            componentProps: {
                configWithFields,
                error,
                isEditForm,
                isQueryBuilderVisible,
                loading,
                queryState,
                setIsQueryBuilderVisible,
                setQueryState
            }
        },
        {
            id: "scanTypes",
            title: "Scan types",
            component: StepScanTypes
        },
        {
            id: "time",
            title: "Time configuration",
            component: StepTimeConfiguration
        },
        {
            id: "advance",
            title: "Advanced settings",
            component: StepAdvancedSettings
        }
    ];

    const readYamlFile = useCallback(
        async (rawApiData) => {
            if (rawApiData) {
                try {
                    const apiData = await OpenAPIParser.dereference(rawApiData);
                    const properties = collectProperties(apiData.components.schemas.Asset);
                    setConfigWithFields(previousConfig => ({ ...previousConfig, fields: properties }))
                } catch (err) {
                    console.error(err);
                }
            }
        },
        [],
    );

    useEffect(() => {
        const currentTree = queryState.tree;
        setQueryState({
            config: configWithFields,
            tree: QbUtils.checkTree(currentTree, configWithFields)
        });
        // eslint-disable-next-line
    }, [configWithFields])

    useEffect(() => {
        readYamlFile(data);
        // eslint-disable-next-line
    }, [data])

    return (
        <WizardModal
            extended={isQueryBuilderVisible}
            title={`${isEditForm ? "Edit" : "New"} scan config`}
            onClose={onClose}
            steps={steps}
            initialValues={INITIAL_SCAN_CONFIG_FORM_VALUES}
            submitUrl={APIS.SCAN_CONFIGS}
            getSubmitParams={formValues => {
                const { id, scheduled, ...submitData } = formValues;

                const { scheduledSelect, laterDate, laterTime, cronLine } = scheduled;
                const isNow = scheduledSelect === SCHEDULE_TYPES_ITEMS.NOW.value;

                let formattedDate = new Date();

                if (!isNow) {
                    const [hours, minutes] = laterTime.split(":");
                    formattedDate = new Date(laterDate);
                    formattedDate.setHours(hours, minutes);
                }

                submitData.scheduled = {};

                if (scheduledSelect === SCHEDULE_TYPES_ITEMS.REPETITIVE.value) {
                    submitData.scheduled.cronLine = cronLine;
                } else {
                    submitData.scheduled.operationTime = formattedDate.toISOString();
                }

                return !isEditForm ? { submitData } : {
                    method: FETCH_METHODS.PUT,
                    formatUrl: url => `${url}/${id}`,
                    submitData
                }
            }}
            onSubmitSuccess={onSubmitSuccess}
            removeTitleMargin={true}
        />
    )
}

export default ScanConfigWizardModal;
