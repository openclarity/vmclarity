# [RFC] Extend asset-finding relationship logic

*Note: this RFC template follows HashiCrop RFC format described [here](https://works.hashicorp.com/articles/rfc-template)*


|               |                                            |
|---------------|--------------------------------------------|
| **Created**   | 2023-01-15                                 |
| **Status**    | WIP\| **InReview** \| Approved \| Obsolete |
| **Owner**     | *ramizpolic*                               |
| **Approvers** | *github handles*                           |

---

This RFC proposes adding `AssetFinding` API model and its supporting logic to improve efficiency and enable aggregation.

## Background

Each finding defines some specific security details (e.g. vulnerability) discovered on an asset.
The same finding can be discovered on multiple assets by different asset scans.
This means that there is many-to-many relationship between findings, assets, and asset scans.
However, the existing [specifications](https://github.com/openclarity/vmclarity/blob/9aa03a8abe22ebddb841a9c28f7a9629f744ced7/api/openapi.yaml#L3395-L3444) 
describe findings with one-to-one relationship to assets and asset scans.
In addition, the [database logic](https://github.com/openclarity/vmclarity/blob/9aa03a8abe22ebddb841a9c28f7a9629f744ced7/pkg/apiserver/database/gorm/finding.go#L103-L105) 
treats every new finding as unique without performing any checks.
Together, this can introduce issues for multiple reasons:

- Each `Finding` is coupled with an `Asset` it was discovered on and the `AssetScan` that discovered it.
  This leads to unnecessary data duplication as each `Finding` with different association to these models will be treated as unique.
  In addition, the lack of the actual uniqueness check completely ignores the data already present in the database.
  Together, this creates performance and memory utilization overheads.
- The model differs from the existing association patterns between models compared to e.g. `Asset`, `AssetScan`, and `AssetScanEstimation`.
  This introduces complexities due to lack of proper aggregation by overloading the `Finding` data returned by the API or shown on the UI.

## Proposal

The finding-related components can be changed in several ways to address the concerns above.
This has been divided into two categories to more easily understand the scope of proposed changes.

### Addressing API changes

Findings should serve as a collection of all security details discovered so far, irrelevant of the assets they were discovered on or the scans that discovered them.
This means that _`Asset` relationship can be dropped from the `Finding` model._

To express the relationship between assets and findings, new `AssetFinding` model can be added. 
This ensures many-to-many relationship between these two models, i.e. a single finding can be discovered on multiple assets, and an asset can contain multiple findings.
Therefore, the `AssetFinding` serves as a bridge table between different assets and findings.
To provide statistical data regarding other models, the `Finding` can be expanded with `summary` property.

Similar approach can be used for asset scans and findings, although this is less relevant for this RFC.
The `AssetScan` relationship in `Finding` can be preserved, but it should be noted that it represents the first scan that discovered this finding.
If required to keep the track of all asset scans that discovered a specific finding, `AssetScanFinding` bridge table can be used.

#### Analysis
The reason to have `AssetFinding` and `AssetScanFinding` as separate tables relates to time and space complexity.
Due to current nature of these models (e.g. creating a new `Finding` or `Asset` for each version, or `AssetScan` on schedule), the size of the tables can grow rapidly.
For example, to keep track of `#assets = 100`, `#assetScans = 100`, and `#findings = 100` requires:
```
Case A: no versioning changes, no additional scheduled scans
RA_1 = R(asset, assetScan, finding) = #asset * #assetScan * #finding = 100^3 items => unified relationship table
RA_2 = R(asset, finding) = #asset * #finding = 100^2 items => AssetFinding table
RA_3 = R(assetScan, finding) = #assetScan * #finding = 100^2 items => AssetScanFinding table

Case B: 10 new versions, 10 new scans for each asset
RB_1 = 10^3 * RA_1 = 10^5 * 100^2 => grows too quickly
RB_2 = 10^2 * RA_2 = 10^2 * 100^2 => depends on assets and findings
RB_3 = 10^2 * RA_3 = 10^2 * 100^2 => scans can be scheduled to occur more often than the versioning changes on other models
```
Therefore, it makes sense to keep track of relationship in a separate tables between these models.
This is also why the `AssetScanFinding` table is omitted from the implementation.

### Addressing uniqueness

The finding database logic can **implement uniqueness check** similar to the existing logic as shown [here](https://github.com/openclarity/vmclarity/blob/9aa03a8abe22ebddb841a9c28f7a9629f744ced7/pkg/apiserver/database/gorm/asset.go#L289).
The data required for the check can be extracted directly from the actual finding.

### Non-goals

This RFC does not intend to propose changes regarding the relationship of findings to asset scans.
In the context of this proposal, it is assumed that the asset scan in a given finding represents the first scan that discovered it.
This behavior, if required, can be addressed as described in [Addressing API changes](#addressing-api-changes) section.

### Abandoned Ideas (Optional)

Adding the aggregation methods to the `uibackend` API was considered but abandoned as it does not address the data duplication issue.

---

## Implementation

### 1. Extend Findings-related API and database logic

```yaml
Finding:
  type: object
  allOf:
    - $ref: '#/components/schemas/Metadata'
    - type: object
      properties:
        id:
          type: string
        assetCount:
          type: integer
          description: List of assets that contain this finding.
          items:
            $ref: '#/components/schemas/AssetRelationship'
        foundBy:
          $ref: '#/components/schemas/AssetScanRelationship'
        foundOn:
          description: When this finding was discovered by a scan
          type: string
          format: date-time
        invalidatedOn:
          description: When this finding was invalidated by a newer scan
          type: string
          format: date-time
        findingInfo:
          anyOf:
            - $ref: '#/components/schemas/PackageFindingInfo'
            - $ref: '#/components/schemas/VulnerabilityFindingInfo'
            - $ref: '#/components/schemas/MalwareFindingInfo'
            - $ref: '#/components/schemas/SecretFindingInfo'
            - $ref: '#/components/schemas/MisconfigurationFindingInfo'
            - $ref: '#/components/schemas/RootkitFindingInfo'
            - $ref: '#/components/schemas/ExploitFindingInfo'
            - $ref: '#/components/schemas/InfoFinderFindingInfo'
          discriminator:
            propertyName: objectType
            mapping:
              Package: '#/components/schemas/PackageFindingInfo'
              Vulnerability: '#/components/schemas/VulnerabilityFindingInfo'
              Malware: '#/components/schemas/MalwareFindingInfo'
              Secret: '#/components/schemas/SecretFindingInfo'
              Misconfiguration: '#/components/schemas/MisconfigurationFindingInfo'
              Rootkit: '#/components/schemas/RootkitFindingInfo'
              Exploit: '#/components/schemas/ExploitFindingInfo'
              InfoFinder: '#/components/schemas/InfoFinderFindingInfo'
```

The API changes impact the database schema and should be handled accordingly.
In addition, the database-related logic such as bootstrapping needs to be updated to reflect these changes.

### 2. Add uniqueness checks to Findings

The uniqueness check can be added similarly to the existing implementations. An example is given below.

```go
func (s *FindingsTableHandler) checkUniqueness(finding models.Finding) (*models.Finding, error) {
  discriminator, err := finding.FindingInfo.ValueByDiscriminator()
  if err != nil {
    return nil, fmt.Errorf("failed to get value by discriminator: %w", err)
  }
  
  var filter string
  switch info := discriminator.(type) {
  case models.PackageFindingInfo:
    filter = fmt.Sprintf("uniqueness query for package finding", info)
  
  case models.VulnerabilityFindingInfo:
    filter = fmt.Sprintf("uniqueness query for vulnerability finding", info)
  
  // implementation of other cases
  
  default:
    return nil, fmt.Errorf("finding type is not supported (%T): %w", discriminator, err)
  }
  
  // implementation of the actual check
  
  return nil, nil
}
```

Once done, the database and controller logic should also be updated to handle cases when:
- the finding needs to be created or updated
- the assets need to be removed from a finding but are not present
- the assets need to be added to a finding but are already present
- any other edge cases not covered above

### 3. Update Finding-related UI components

### 4. Add AssetFinding logic

- API logic
- Database handlers
- Controllers

## UX

This RFC has no visible impacts on the UX.

## UI

This RFC changes the following UI components:
- `/findings/{findingType}` table drops the _Asset Name_ and _Asset location_ columns.
  Instead, they are replaced with a single _Asset Count_ column that shows the number of assets related to a finding.
- `/findings/{findingType}/{findingID}` removes _Asset details_ menu item.
  It is replaced with _See assets_ button on finding summary.
  The button redirects to the `/assets` page and displays filtered, finding-related assets by using OData queries.
  - Note: This can be done similarly to displaying the related _Asset Scans_ in `/scans/scans/{scanID}`.
    Alternatively, a table can be directly used instead of redirects, but might require more changes.
