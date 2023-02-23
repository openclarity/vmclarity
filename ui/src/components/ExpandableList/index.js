import React, { useState } from 'react';
import classnames from 'classnames';
import Arrow, { ARROW_NAMES } from 'components/Arrow';

import './expandable-list.scss';

const MINIMAL_LEN = 1;

const ExpandableList = ({items, withTagWrap=false}) => {
    const allItems = items || [];
    const minimalItems = allItems.slice(0, MINIMAL_LEN);

    const [itemsToDisplay, setItemsToDisplay] = useState(allItems.length > MINIMAL_LEN ? minimalItems : allItems);
    const isOpen = itemsToDisplay.length === allItems.length;

    return (
        <div>
            <div className="expandable-list-display-wrapper">
                <div className="expandable-list-items">
                    {
                        itemsToDisplay.map((item, index) => (
                            <div key={index} className="expandable-list-item-wrapper">
                                <div className={classnames("expandable-list-item", {"with-tag-wrap": withTagWrap})}>{item}</div>
                            </div>
                        ))
                    }
                </div>
                {minimalItems.length !== allItems.length &&
                    <Arrow
                        name={isOpen ? ARROW_NAMES.TOP : ARROW_NAMES.BOTTOM}
                        onClick={event => {
                            event.stopPropagation();
                            event.preventDefault();
                            
                            setItemsToDisplay(isOpen ? minimalItems : allItems);
                        }}
                        small
                    />
                }
            </div>
        </div>
    )
}

export default ExpandableList;