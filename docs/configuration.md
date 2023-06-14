# Configuration

## Orchestrator

| Environment Variable                      | Required  | Default | Description                                  |
|-------------------------------------------|-----------|---------|----------------------------------------------|
| `DELETE_JOB_POLICY`                       |           |         |                                              |
| `SCANNER_CONTAINER_IMAGE`                 |           |         |                                              |
| `GITLEAKS_BINARY_PATH`                    |           |         |                                              |
| `CLAM_BINARY_PATHCLAM_BINARY_PATH`        |           |         |                                              |
| `FRESHCLAM_BINARY_PATH`                   |           |         |                                              |
| `ALTERNATIVE_FRESHCLAM_MIRROR_URL`        |           |         |                                              |
| `LYNIS_INSTALL_PATH`                      |           |         |                                              |
| `SCANNER_VMCLARITY_BACKEND_ADDRESS`       |           |         |                                              |
| `EXPLOIT_DB_ADDRESS`                      |           |         |                                              |
| `TRIVY_SERVER_ADDRESS`                    |           |         |                                              |
| `TRIVY_SERVER_TIMEOUT`                    |           |         |                                              |
| `GRYPE_SERVER_ADDRESS`                    |           |         |                                              |
| `GRYPE_SERVER_TIMEOUT`                    |           |         |                                              |
| `CHKROOTKIT_BINARY_PATH`                  |           |         |                                              |
| `SCAN_CONFIG_POLLING_INTERVAL`            |           |         |                                              |
| `SCAN_CONFIG_RECONCILE_TIMEOUT`           |           |         |                                              |
| `SCAN_POLLING_INTERVAL`                   |           |         |                                              |
| `SCAN_RECONCILE_TIMEOUT`                  |           |         |                                              |
| `SCAN_TIMEOUT`                            |           |         |                                              |
| `SCAN_RESULT_POLLING_INTERVAL`            |           |         |                                              |
| `SCAN_RESULT_RECONCILE_TIMEOUT`           |           |         |                                              |
| `SCAN_RESULT_PROCESSOR_POLLING_INTERVAL`  |           |         |                                              |
| `SCAN_RESULT_PROCESSOR_RECONCILE_TIMEOUT` |           |         |                                              |
| `DISCOVERY_INTERVAL`                      |           |         |                                              |
| `CONTROLLER_STARTUP_DELAY`                |           |         |                                              |
| `PROVIDER`                                | **yes**   | `aws`   | Provider used for Target discovery and scans |

## Provider

### AWS

| Environment Variable                         | Required  | Default                                      | Description                                                                                                                                                                                                                             |
|----------------------------------------------|-----------|----------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `VMCLARITY_AWS_REGION`                       | **yes**   |                                              | Region where the Scanner instance needs to be created                                                                                                                                                                                   |
| `VMCLARITY_AWS_SUBNET_ID`                    | **yes**   |                                              | SubnetID where the Scanner instance needs to be created                                                                                                                                                                                 |
| `VMCLARITY_AWS_SECURITY_GROUP_ID`            | **yes**   |                                              | SecurityGroupId which needs to be attached to the Scanner instance                                                                                                                                                                      |
| `VMCLARITY_AWS_KEYPAIR_NAME`                 |           |                                              | Name of the SSH KeyPair to use for Scanner instance launch                                                                                                                                                                              |
| `VMCLARITY_AWS_IMAGE_NAME_FILTER`            |           | `ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-*` | The Name filter used for finding AMI to instantiate Scanner instance. See: [DescribeImages](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeImages.html)                                                             |
| `VMCLARITY_AWS_IMAGE_OWNERS`                 |           | `099720109477`                               | Comma separated list of OwnerID(s)/OwnerAliases used as Owners filter for finding AMI to instantiate Scanner instance. See: [DescribeImages](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeImages.html)            |
| `VMCLARITY_AWS_INSTANCE_MAPPING`             |           | `x86_64:t3.large,arm64:t4g.large`            | Comma separated list of architecture:instance_type pairs used for VMClarity Scanner instance                                                                                                                                            |
| `VMCLARITY_AWS_SCANNER_INSTANCE_ARCH_TO_USE` |           |                                              | Architecture to be used for Scanner instance which prevent the Provider to dynamically determine it based on the Target architecture. The Provider will use this value to lookup for instance type in `VMCLARITY_AWS_INSTANCE_MAPPING`. |
| `VMCLARITY_AWS_BLOCK_DEVICE_NAME`            |           | `xvdh`                                       | Block device name used for attaching Scanner volume to the Scanner instance                                                                                                                                                             |
