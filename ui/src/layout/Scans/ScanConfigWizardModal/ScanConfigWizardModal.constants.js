import { Utils as QbUtils } from '@react-awesome-query-builder/mui';

const EMPTY_JSON_SCOPE_TREE = { "id": QbUtils.uuid(), "type": "group" };
const SCOPE_TREE_KEY = "openclarity.io/vmclarity-ui/query-builder/data";
const SCOPE_CONFIG_KEY = "openclarity.io/vmclarity-ui/query-builder/config";

export { EMPTY_JSON_SCOPE_TREE, SCOPE_TREE_KEY, SCOPE_CONFIG_KEY };
