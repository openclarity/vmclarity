import TitleValueDisplay, { TitleValueDisplayRow } from 'components/TitleValueDisplay';
import { formatDate } from 'utils/utils';


export const getScanColumnsConfigList = () => ([]);
export const FindingsDetailsCommonFields = ({firstSeen, lastSeen}) => (
    <>
        <TitleValueDisplayRow>
            <TitleValueDisplay title="First seen">{formatDate(firstSeen)}</TitleValueDisplay>
        </TitleValueDisplayRow>
        <TitleValueDisplayRow>
            <TitleValueDisplay title="Last seen">{formatDate(lastSeen)}</TitleValueDisplay>
        </TitleValueDisplayRow>
    </>
)
