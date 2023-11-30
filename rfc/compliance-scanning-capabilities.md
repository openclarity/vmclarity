# [RFC] Add compliance scanning capabilities

*Note: this RFC template follows HashiCrop RFC format described [here](https://works.hashicorp.com/articles/rfc-template)*


|               |                                            |
| ------------- | ------------------------------------------ |
| **Created**   | 2023-11-29                                 |
| **Status**    | WIP\| **InReview** \| Approved \| Obsolete |
| **Owner**     | *ramizpolic*                               |
| **Approvers** | *TODO*                                     |

---

This RFC proposes the addition of a new family of scanners for compliance to enrich security findings on assets.

## Background

The migration of [CIS Docker Benchmark](https://github.com/goodwithtech/dockle) scanner from KubeClarity requires extensions of both the backend and UI code to accommodate the implemented logic.

The scanners in VMClarity are described based on their respective security scopes, for example, vulnerabilities and misconfigurations. This ensures that the types of findings on assets are well-defined and respected. In KubeClarity, the CIS Docker Benchmark scanner defines its own findings as described in the [API specs](https://github.com/openclarity/kubeclarity/blob/5ac3048b7a782c900a9bef846a91a7735ba77e24/api/swagger.yaml#L243C26-L243C26). This makes the migration of scanning capabilities to VMClarity problematic for two main reasons:

- Logic in the form of a new independent scanner family does not conform to any supported *security scopes*. CIS Docker Benchmark provides little benefit on its own due to scope constraints compared to the existing scanner families.

* Logic is *too specific* to be part of any existing scanner families. CIS Docker Benchmark findings cannot be uniformly converted to other findings without some loss of data validity and integrity.

## Proposal

Consider the CIS Docker Benchmark scanner as a part of a new, more generalized family of scanners for compliance checks on assets. The new **compliance scanner family** would be used to scan for more generic findings such as best practice violations and other common security issues not discoverable by other families. This approach benefits VMClarity in several ways:

* The new compliance findings define a better security scope which natively extends the overall security findings on assets

The compliance findings fit well with the current API specifications. The new security scope offers extra flexibility by supporting additional security findings in a generic way. The migration needs only to address minor changes required to convert the CIS Docker Benchmark results into generic compliance findings.

- The new family of scanners enables an idiomatic way to migrate the required scanning logic from KubeClarity

The required boilerplate logic for the new scanner family can be added to VMClarity in advance. This is accomplished by reusing the existing patterns to minimize changes. The migration of CIS Docker Benchmark logic can then be performed as an implementation of a specific scanner within the new compliance family.

*TL;DR Adding the compliance scanner family provides value in terms of better scanning capabilities. Security findings are extended in a generic and reusable way. The alternative approach adds the CIS Docker Benchmark scanner family directly to avoid overgeneralization of data. The main difference between the approaches relates to the API specifications which are reflected in the usage and extension possibilities.*

### Abandoned Ideas (Optional)

---

## Implementation

The implementation will be performed in three stages:

1. Add scanner family boilerplate to prepare for CIS Docker Benchmark migration from KubeClarity

We can start by defining a generalized API model to describe compliance findings. The `Compliance` model is loosely inspired by the CIS Docker Benchmark from KubeClarity to support direct conversion.

```yaml
Compliance:
  type: object
  properties:
    scannerName:
      type: string
    code: # CISDockerBenchmarkResultsEX.code
      type: string
      description: Compliance violation code, if applicable (e.g. https://github.com/goodwithtech/dockle/blob/master/CHECKPOINT.md checkpoint codes)
    location: # returned by the CIS Docker Benchmark scanner
      type: string
      description: Location within the asset where the violation was recorded (e.g. filesystem path)
    reason: # CISDockerBenchmarkResultsEX.title
      type: string
      description: Short info about why the compliance violation was flagged
    summary: # CISDockerBenchmarkResultsEX.desc
      type: string
      description: Additional context such as the potential impact
    remediation: # could be returned by new scanners in the future
      type: string
      description: Possible fix for this compliance violation
    severity: # CISDockerBenchmarkResultsEX.level
      type: string
      enum:
        - HIGH
        - MEDIUM
        - LOW
```

Additional API changes required to enable compliance scanners are summarized in the table below. Affected components not found in the table should be easy to identify during the actual implementation.


| Action | Models                                                                                                             |
| ------ | ------------------------------------------------------------------------------------------------------------------ |
| Create | `Compliance` `ComplianceScan` `ComplianceScanSummary` <br /> `ComplianceFindingInfo` `CompliancesConfig`           |
| Extend | `ScanType` `Finding` `FindingInfo` `ScanFamiliesConfig` <br />`AssetScan` `AssetScanStats` `AssetScanRelationship` |

The backend changes can be performed by following the existing patterns. For example, the vulnerability scanner family can serve as a reference. Extensions to the API on the UI backend side are also needed to add support for the new findings.

2. Extend UI to support the new API changes

UI components affected by the API changes can easily be identified from the same table. Supporting UI code can similarly be updated by reusing the patterns defined for the existing findings. Note that it might not be possible to reuse the UI-related code from KubeClarity (to be confirmed).

3. Migrate CIS Docker Benchmark scanner from KubeClarity as part of compliance scanners

CIS Docker Benchmark defined [here](https://github.com/openclarity/kubeclarity/tree/5ac3048b7a782c900a9bef846a91a7735ba77e24/cis_docker_benchmark_scanner) can be migrated under the newly created compliance family of scanners. However, the logic defined [here](https://github.com/openclarity/kubeclarity/blob/5ac3048b7a782c900a9bef846a91a7735ba77e24/cis_docker_benchmark_scanner/pkg/report/report.go) needs to be slightly changed to use the new compliance findings models.

## UX

Users will be able to see the new findings category on the dashboard. When creating scan configurations, users can enable or disable compliance scans. When checking the results of a specific scan, users can check the findings tab to get the compliance findings summary.

Users can get compliance findings and summaries for a specific asset or asset scan via a separate tab. The findings page will have an additional tab for compliance. Users can filter specific compliance findings to see basic details. When checking specific compliance findings, users can navigate between finding and asset tabs to get additional information.

## UI

Relevant UI changes can reuse existing logic.
