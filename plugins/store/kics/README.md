# KICS

> **KICS** is a scanner application, that uses [Checkmarx KICS](https://checkmarx.com/product/opensource/kics-open-source-infrastructure-as-code-project/) (Keeping Infrastructure as Code Secure) to scan your Infrastructure as Code (IaC) files for misconfigurations. It's designed to be used as a plugin for the [VMClarity](https://openclarity.io/docs/vmclarity/) platform.

## Usage

Make a POST request to the `/assetScans` endpoint, to initiate a KICS scan. The body of the POST request should include a JSON object with the configuration for the scan.

> NOTE: The follwing is a minimal example. Your actual configuration should have additional properties.

```json
{
    "name": "scan-name",
    "scanTemplate": {
        "scope": "contains(assetInfo.labels, '{\"key\":\"scanconfig\",\"value\":\"test\"}')",
        "assetScanTemplate": {
            "scanFamiliesConfig": {
                "plugins": {
                    "enabled": true,
                    "scanners_list": ["KICS"],
                    "scanners_config": {
                        "image_name": "vmclarity-kics-scanner:latest",
                        "output_file": "kics-scan-out",
                        "config": "kics-config.json"
                    }
                }
            }
        }
    }
}
```

### Important Notes

- The KICS scanner is designed to be started by the **VMClarity** runner, therefore running it as a standalone tool is not recommended.

- The `config` property in the POST request should point to a file on the host filesystem with the [parameters](https://github.com/Checkmarx/kics/blob/e387aa2505a3207e1087520972e0e52f7e0e6fdf/pkg/scan/client.go#L54) that the `KICS` client will use.

- The configuration file can be in any of the following formats: `JSON`, `TOML`, `YAML`, or `HCL`.

- Please note that not all `scan parameters` are currently supported by the scanner.

When the scan is done the output can be found at the `<specified output file>.json` formatted the following way:

> Each misconfiguration is represented as an object with the following properties, the output file will contain an array of these misconfiguration objects:

- `scannerName`: The name of the scanner that detected the misconfiguration.
- `id`: Check or test ID, if applicable (e.g. Lynis TestID, CIS Docker Benchmark checkpoint code, etc).
- `location`: Location within the asset where the misconfiguration was recorded (e.g. filesystem path).
- `category`: Specifies misconfiguration impact category.
- `message`: A short description of the misconfiguration.
- `description`: Additional context, such as the potential impact of the misconfiguration.
- `remediation`: A possible fix for the misconfiguration.
- `severity`: The severity of the misconfiguration, which can be one of the following:

  - `MisconfigurationHighSeverity`
  - `MisconfigurationMediumSeverity`
  - `MisconfigurationLowSeverity`
  - `MisconfigurationInfoSeverity`
