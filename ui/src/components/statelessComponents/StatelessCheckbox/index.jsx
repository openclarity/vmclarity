import React from 'react';
import classnames from 'classnames';
import InfoIcon from 'components/InfoIcon';
import FieldError from 'components/Form/FieldError';
import FieldLabel from 'components/Form/FieldLabel';

import './StatelessCheckbox.scss';

const StatelessCheckbox = ({
    className,
    disabled,
    errorMessage,
    label,
    name,
    setValue,
    title,
    tooltipText,
    value
}) => {

    const tooltipId = `form-tooltip-${name}`;

    return (
        <div className={classnames("form-field-wrapper", "checkbox-field-wrapper", className)}>
            {!!label && <FieldLabel tooltipId={tooltipId} tooltipText={tooltipText}>{label}</FieldLabel>}
            <label className={classnames("checkbox-wrapper", { disabled })}>
                <div className="inner-checkbox-wrapper">
                    <input type="checkbox" value={value} name={name} onChange={() => disabled ? null : setValue(e => !e)} />
                    <span className="checkmark"></span>
                </div>
                <span className="checkbox-title">{title}</span>
                {!label && !!tooltipText && <div style={{ marginLeft: "5px" }}><InfoIcon tooltipId={tooltipId} tooltipText={tooltipText} /></div>}
            </label>
            {errorMessage && <FieldError>{errorMessage}</FieldError>}
        </div>
    )
}

export default StatelessCheckbox;