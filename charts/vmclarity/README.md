# vmclarity

![Version: v0.0.0-latest](https://img.shields.io/badge/Version-v0.0.0--latest-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

VMClarity is an open source tool for agentless detection and management of
Virtual Machine Software Bill Of Materials (SBOM) and security threats such
as vulnerabilities, exploits, malware, rootkits, misconfigurations and leaked
secrets.

**Homepage:** <openclarity.io>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| VMClarity Maintainers |  | <https://github.com/openclarity/vmclarity> |

## Source Code

* <https://github.com/openclarity/vmclarity>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.bitnami.com/bitnami | postgresql | 12.7.1 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| apiserver.containerSecurityContext.enabled | bool | `true` | API Server container security context enabled |
| apiserver.containerSecurityContext.runAsNonRoot | bool | `true` | Whether the API Server containers should run as a non-root user |
| apiserver.containerSecurityContext.runAsUser | int | `1001` | User ID which the API Server container should run as |
| apiserver.image.pullPolicy | string | `"IfNotPresent"` | API Server image pull policy |
| apiserver.image.registry | string | `"ghcr.io"` | API Server image registry |
| apiserver.image.repository | string | `"openclarity/vmclarity-apiserver"` | API Server image repositiory |
| apiserver.image.tag | string | `"latest"` | API Server image tag (immutable tags are recommended) |
| apiserver.logLevel | string | `"info"` | API Server log level |
| apiserver.podSecurityContext.enabled | bool | `true` | API Server pod's security context enabled |
| apiserver.podSecurityContext.fsGroup | int | `1001` | API Server pod's security context fsGroup |
| apiserver.replicas | int | `1` | Number of replicas for the API Server |
| apiserver.resources.limits | object | `{}` | The resources limits for the apiserver containers |
| apiserver.resources.requests | object | `{}` | The requested resources for the apiserver containers |
| exploitDBServer.containerSecurityContext.enabled | bool | `true` |  |
| exploitDBServer.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| exploitDBServer.containerSecurityContext.runAsUser | int | `1001` |  |
| exploitDBServer.image.pullPolicy | string | `"IfNotPresent"` |  |
| exploitDBServer.image.registry | string | `"ghcr.io"` |  |
| exploitDBServer.image.repository | string | `"openclarity/exploit-db-server"` |  |
| exploitDBServer.image.tag | string | `"v0.2.3"` |  |
| exploitDBServer.podSecurityContext.enabled | bool | `true` |  |
| exploitDBServer.podSecurityContext.fsGroup | int | `1001` |  |
| exploitDBServer.replicas | int | `1` |  |
| exploitDBServer.resources.limits | object | `{}` | The resources limits for the exploit-db-server containers |
| exploitDBServer.resources.requests | object | `{}` | The requested resources for the exploit-db-server containers |
| freshclamMirror.containerSecurityContext.enabled | bool | `false` |  |
| freshclamMirror.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| freshclamMirror.containerSecurityContext.runAsUser | int | `1001` |  |
| freshclamMirror.image.pullPolicy | string | `"IfNotPresent"` |  |
| freshclamMirror.image.registry | string | `"ghcr.io"` |  |
| freshclamMirror.image.repository | string | `"openclarity/freshclam-mirror"` |  |
| freshclamMirror.image.tag | string | `"v0.1.0"` |  |
| freshclamMirror.podSecurityContext.enabled | bool | `false` |  |
| freshclamMirror.podSecurityContext.fsGroup | int | `1001` |  |
| freshclamMirror.replicas | int | `1` |  |
| freshclamMirror.resources.limits | object | `{}` | The resources limits for the freshclam mirror containers |
| freshclamMirror.resources.requests | object | `{}` | The requested resources for the freshclam mirror containers |
| gateway.containerSecurityContext.enabled | bool | `false` |  |
| gateway.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| gateway.containerSecurityContext.runAsUser | int | `1001` |  |
| gateway.image.pullPolicy | string | `"IfNotPresent"` |  |
| gateway.image.registry | string | `"docker.io"` |  |
| gateway.image.repository | string | `"library/nginx"` |  |
| gateway.image.tag | string | `"1.25.1"` |  |
| gateway.podSecurityContext.enabled | bool | `false` |  |
| gateway.podSecurityContext.fsGroup | int | `1001` |  |
| gateway.replicas | int | `1` |  |
| gateway.resources.limits | object | `{}` | The resources limits for the gateway containers |
| gateway.resources.requests | object | `{}` | The requested resources for the gateway containers |
| global.imageRegistry | string | `""` |  |
| grypeServer.containerSecurityContext.enabled | bool | `true` |  |
| grypeServer.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| grypeServer.containerSecurityContext.runAsUser | int | `1001` |  |
| grypeServer.image.pullPolicy | string | `"IfNotPresent"` |  |
| grypeServer.image.registry | string | `"ghcr.io"` |  |
| grypeServer.image.repository | string | `"openclarity/grype-server"` |  |
| grypeServer.image.tag | string | `"v0.4.0"` |  |
| grypeServer.logLevel | string | `"info"` |  |
| grypeServer.podSecurityContext.enabled | bool | `true` |  |
| grypeServer.podSecurityContext.fsGroup | int | `1001` |  |
| grypeServer.replicas | int | `1` |  |
| grypeServer.resources.limits | object | `{}` | The resources limits for the grype server containers |
| grypeServer.resources.requests | object | `{}` | The requested resources for the grype server containers |
| orchestrator.aws.keypairName | string | `""` | KeyPair to use for the scanner instance |
| orchestrator.aws.region | string | `""` | Region where the control plane is running |
| orchestrator.aws.scannerAmiId | string | `""` | AMI to use for the scanner instance |
| orchestrator.aws.scannerInstanceType | string | `""` | InstanceType to use for the scanner instance |
| orchestrator.aws.scannerRegion | string | `""` | Region where the scanners will be created |
| orchestrator.aws.securityGroupId | string | `""` | Security Group to use for the scanner networking |
| orchestrator.aws.subnetId | string | `""` | Subnet where the scanners will be created |
| orchestrator.azure.scannerImageOffer | string | `""` | Scanner VM source image offer |
| orchestrator.azure.scannerImagePublisher | string | `""` | Scanner VM source image publisher |
| orchestrator.azure.scannerImageSku | string | `""` | Scanner VM source image sku |
| orchestrator.azure.scannerImageVersion | string | `""` | Scanner VM source image version |
| orchestrator.azure.scannerLocation | string | `""` | Location where the scanner instances will be run |
| orchestrator.azure.scannerPublicKey | string | `""` | SSH RSA Public Key to configure the scanner instances with |
| orchestrator.azure.scannerResourceGroup | string | `""` | ResourceGroup where the scanner instances will be run |
| orchestrator.azure.scannerSecurityGroup | string | `""` | Scanner VM security group |
| orchestrator.azure.scannerStorageAccountName | string | `""` | Storage account to use for transfering snapshots between regions |
| orchestrator.azure.scannerStorageContainerName | string | `""` | Storage container to use for transfering snapshots between regions |
| orchestrator.azure.scannerSubnetId | string | `""` | Subnet ID where the scanner instances will be run |
| orchestrator.azure.scannerVmSize | string | `""` | Scanner VM size |
| orchestrator.azure.subscriptionId | string | `""` | Subscription ID for discovery and scanning |
| orchestrator.containerSecurityContext.enabled | bool | `true` | Whether Orchestrator container secuirty context is enabled |
| orchestrator.containerSecurityContext.runAsNonRoot | bool | `true` | Whether the Orchestrator containers should as a non-root user |
| orchestrator.containerSecurityContext.runAsUser | int | `1001` | User ID which the Orchestrator container should run as |
| orchestrator.deleteJobPolicy | string | `"Always"` | Global policy used to determine when to clean up an AssetScan. Possible options are: Always - All AssetScans are cleaned up OnSuccess - Only Successful AssetScans are cleaned up, Failed ones are left for debugging Never - No AssetScans are cleaned up |
| orchestrator.freshclamMirrorAddress | string | `""` | Address that scanenrs can use to reach the freshclam mirror |
| orchestrator.gcp.projectId | string | `""` | Project ID for discovery and scanning |
| orchestrator.gcp.scannerMachineType | string | `""` | Scanner Machine type |
| orchestrator.gcp.scannerSourceImage | string | `""` | Scanner source image |
| orchestrator.gcp.scannerSubnet | string | `""` | Subnet where to run the scanner instances |
| orchestrator.gcp.scannerZone | string | `""` | Zone to where the scanner instances should run |
| orchestrator.grypeServerAddress | string | `""` | Address that scanners can use to reach the grype server |
| orchestrator.image.pullPolicy | string | `"IfNotPresent"` | Orchestrator image pull policy |
| orchestrator.image.registry | string | `"ghcr.io"` | Orchestrator image registry |
| orchestrator.image.repository | string | `"openclarity/vmclarity-orchestrator"` | Orchestrator image repository |
| orchestrator.image.tag | string | `"latest"` | Orchestrator image tag (immutable tags are recommended) |
| orchestrator.logLevel | string | `"info"` | Orchestrator service log level |
| orchestrator.podSecurityContext.enabled | bool | `true` | Whether Orchestrator pod security context is enabled |
| orchestrator.podSecurityContext.fsGroup | int | `1001` | Orchestrator pod security context fsGroup |
| orchestrator.provider | string | `"aws"` | Which provider to enable |
| orchestrator.replicas | int | `1` | Number of replicas for the Orchestrator service Currently 1 supported. |
| orchestrator.resources.limits | object | `{}` | The resources limits for the orchestrator containers |
| orchestrator.resources.requests | object | `{}` | The requested resources for the orchestrator containers |
| orchestrator.scannerApiserverAddress | string | `""` | Address that scanners can use to reach back to the API server |
| orchestrator.scannerImage.registry | string | `"ghcr.io"` | Scanner Container image registry |
| orchestrator.scannerImage.repository | string | `"openclarity/vmclarity-cli"` | Scanner Container image repository |
| orchestrator.scannerImage.tag | string | `"latest"` | Scanner Container image tag (immutable tags are recommended) |
| orchestrator.trivyServerAddress | string | `""` | Address that scanners can use to reach trivy server |
| postgresql.auth.database | string | `"vmclarity"` |  |
| postgresql.auth.password | string | `"password1"` |  |
| postgresql.auth.username | string | `"vmclarity"` |  |
| postgresql.containerSecurityContext.enabled | bool | `true` |  |
| postgresql.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| postgresql.containerSecurityContext.runAsUser | int | `1001` |  |
| postgresql.image.pullPolicy | string | `"IfNotPresent"` |  |
| postgresql.image.registry | string | `"docker.io"` |  |
| postgresql.image.repository | string | `"bitnami/postgresql"` |  |
| postgresql.image.tag | string | `"14.6.0-debian-11-r31"` |  |
| postgresql.podSecurityContext.enabled | bool | `true` |  |
| postgresql.podSecurityContext.fsGroup | int | `1001` |  |
| postgresql.resources.limits | object | `{}` | The resources limits for the postgresql containers |
| postgresql.resources.requests | object | `{}` | The requested resources for the postgresql containers |
| postgresql.service.ports.postgresql | int | `5432` |  |
| trivyServer.containerSecurityContext.enabled | bool | `true` |  |
| trivyServer.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| trivyServer.containerSecurityContext.runAsUser | int | `1001` |  |
| trivyServer.image.pullPolicy | string | `"IfNotPresent"` |  |
| trivyServer.image.registry | string | `"docker.io"` |  |
| trivyServer.image.repository | string | `"aquasec/trivy"` |  |
| trivyServer.image.tag | string | `"0.41.0"` |  |
| trivyServer.podSecurityContext.enabled | bool | `true` |  |
| trivyServer.podSecurityContext.fsGroup | int | `1001` |  |
| trivyServer.replicas | int | `1` |  |
| trivyServer.resources.limits | object | `{}` | The resources limits for the trivy server containers |
| trivyServer.resources.requests | object | `{}` | The requested resources for the trivy server containers |
| ui.containerSecurityContext.enabled | bool | `false` |  |
| ui.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| ui.containerSecurityContext.runAsUser | int | `1001` |  |
| ui.image.pullPolicy | string | `"IfNotPresent"` |  |
| ui.image.registry | string | `"ghcr.io"` |  |
| ui.image.repository | string | `"openclarity/vmclarity-ui"` |  |
| ui.image.tag | string | `"latest"` |  |
| ui.podSecurityContext.enabled | bool | `false` |  |
| ui.podSecurityContext.fsGroup | int | `1001` |  |
| ui.replicas | int | `1` |  |
| ui.resources.limits | object | `{}` | The resources limits for the UI containers |
| ui.resources.requests | object | `{}` | The requested resources for the UI containers |
| uibackend.containerSecurityContext.enabled | bool | `true` |  |
| uibackend.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| uibackend.containerSecurityContext.runAsUser | int | `1001` |  |
| uibackend.image.pullPolicy | string | `"IfNotPresent"` |  |
| uibackend.image.registry | string | `"ghcr.io"` |  |
| uibackend.image.repository | string | `"openclarity/vmclarity-uibackend"` |  |
| uibackend.image.tag | string | `"latest"` |  |
| uibackend.logLevel | string | `"info"` |  |
| uibackend.podSecurityContext.enabled | bool | `true` |  |
| uibackend.podSecurityContext.fsGroup | int | `1001` |  |
| uibackend.replicas | int | `1` |  |
| uibackend.resources.limits | object | `{}` | The resources limits for the UI backend containers |
| uibackend.resources.requests | object | `{}` | The requested resources for the UI backend containers |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
