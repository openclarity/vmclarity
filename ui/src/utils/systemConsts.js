import { ICON_NAMES } from 'components/Icon';

import COLORS from 'utils/scss_variables.module.scss';

export const ROUTES = {
    DEFAULT: "/",
    SCANS: "/scans",
    ASSETS: "/assets",
    ASSET_SCANS: "/assets-scans",
    FINDINGS: "/findings"
}

export const APIS = {
    SCANS: "scans",
    SCAN_CONFIGS: "scanConfigs",
    ASSETS: "targets",
    ASSET_SCANS: "scanResults",
    SCOPES_DISCOVERY: "discovery/scopes",
    FINDINGS: "findings"
}

export const FINDINGS_MAPPING = {
    EXPLOITS: {dataKey: "totalExploits", title: "Exploits", icon: ICON_NAMES.BOMB, color: COLORS["color-main"]},
    MISCONFIGURATIONS: {dataKey: "totalMisconfigurations", title: "Misconfigurations", icon: ICON_NAMES.COG, color: COLORS["color-findings-1"]},
    SECRETS: {dataKey: "totalSecrets", title: "Secrets", icon: ICON_NAMES.KEY, color: COLORS["color-findings-2"]},
    MALWARE: {dataKey: "totalMalware", title: "Malwares", icon: ICON_NAMES.BUG, color: COLORS["color-findings-3"]},
    ROOTKITS: {dataKey: "totalRootkits", title: "Rootkits", icon: ICON_NAMES.GHOST, color: COLORS["color-findings-4"]},
    PACKAGES: {dataKey: "totalPackages", title: "Packages", icon: ICON_NAMES.PACKAGE, color: COLORS["color-findings-5"]}
}