import React from 'react';
import classnames from 'classnames';
import FieldError from 'components/Form/FieldError';
import FieldLabel from 'components/Form/FieldLabel';

import './StatelessRadioButtonGroup.scss';

const StatelessRadioButtonGroup = ({
    className,
    disabled = false,
    errorMessage,
    items,
    label,
    name,
    setValue,
    tooltipText,
    value: currentValue
}) => {

    return (
        <div className={classnames("form-field-wrapper", "radio-field-wrapper", className)}>
            {!!label && <FieldLabel tooltipId={`form-tooltip-${name}`} tooltipText={tooltipText}>{label}</FieldLabel>}
            {
                items.map(({ value, label }) => (
                    <label key={value} className="radio-field-item">
                        <span className="radio-text">{label}</span>
                        <input
                            checked={value === currentValue}
                            disabled={disabled}
                            name={name}
                            onChange={() => setValue(value)}
                            type="radio"
                            value={value}
                        />
                        <span className="checkmark"></span>
                    </label>
                ))
            }
            {errorMessage && <FieldError>{errorMessage}</FieldError>}
        </div>
    )
}

export default StatelessRadioButtonGroup;