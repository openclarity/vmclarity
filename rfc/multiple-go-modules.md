# [RFC] Split VMClarity into multiple go modules

*Note: this RFC template follows HashiCrop RFC format described [here](https://works.hashicorp.com/articles/rfc-template)*

|               |                                             |
| ------------- | ------------------------------------------- |
| **Created**   | 2024-01-17                                  |
| **Status**    | **WIP** \| InReview \| Approved \| Obsolete |
| **Owner**     | paralta                                     |
| **Approvers** | vmclarity-maintainers                       |

---

This RFC proposes a new project structure which splits the main module into multiple modules to allow us to version modules separately and reduce dependencies.

## Background

VMClarity consists of a single module containing all the packages required to run VMClarity. However, there is a clear boundary between the following packages: api, cli, orchestrator, ui and uibackend.

The build and deployment of these packages is already performed separately and independently.

## Proposal

The proposal here is to split the VMClarity repository into multiple modules:

- **api**. Interface between all the services in VMClarity including the DB. Composed by API model, backend client and server.
- **scanner** or **cli**. Responsible for running a scan in an asset and report the results back to api. Contains the logic to configure, run and manage different analysers and scanners. 
- **orchestrator**. Responsible for managing scan configurations, scans, assets and estimations.
- **provider**. Responsible for discovery and scan infrastructure setup for each provider. Contains logic to find assets and run scans on AWS, GCP, Azure, Docker and Kubernetes.
- **uibackend**. Responsible for offloading the ui from data processing and filtering. Slightly coupled with ui. Composed by API model, backend client and server.
- **utils**. Contains packages shared between modules.

Each module will have its own go.mod file and each module will be versioned independently.

## Implementation

The scope of this RFC is not to change code logic but to change code structure. Therefore, the following table describes the path changes for each package impacted.

| Module       | Current path                  | New path                      |
| ------------ | ----------------------------- | ----------------------------- |
| api          | pkg/apiserver                 | api/server                    |
| api          | pkg/shared/backendclient      | api/client                    |
| scanner      | pkg/cli                       | scanner/cli                   |
| scanner      | pkg/shared/analyzer           | scanner/analyzer              |
| scanner      | pkg/shared/config             | scanner/config                |
| scanner      | pkg/shared/converter          | scanner/converter             |
| scanner      | pkg/shared/families           | scanner/families              |
| scanner      | pkg/shared/findingkey         | scanner/findingkey            |
| scanner      | pkg/shared/job_manager        | scanner/jobmanager            |
| scanner      | pkg/shared/scanner            | scanner/scanner               |
| scanner      | pkg/shared/utils              | scanner/utils                 |
| orchestrator | pkg/orchestrator              | orchestrator                  |
| orchestrator | pkg/containerruntimediscovery | orchestrator/runtimediscovery |
| uibackend    | pkg/uibackend                 | uibackend                     |
| uibackend    | pkg/uibackend/rest            | uibackend/server              |
| uibackend    | pkg/shared/uibackendclient    | uibackend/client              |
| utils        | pkg/version                   | utils/version                 |
| utils        | pkg/shared/command            | utils/command                 |
| utils        | pkg/shared/fsutils            | utils/fsutils                 |
| utils        | pkg/shared/log                | utils/log                     |
| utils        | pkg/shared/manifest           | utils/manifest                |


Furthermore, the provider could be removed from the orchestrator.

| Module       | Current path                  | New path                      |
| ------------ | ----------------------------- | ----------------------------- |
| provider     | pkg/orchestrator/provider     | provider                      |

The Dockerfiles for each package will be moved to the corresponding directory. Makefile, GitHub workflows and other files will need to be updated.

## UX/UI

This RFC has no user-impacting changes.
