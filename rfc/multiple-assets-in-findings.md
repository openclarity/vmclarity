# [RFC] Allow multiple assets in findings

*Note: this RFC template follows HashiCrop RFC format described [here](https://works.hashicorp.com/articles/rfc-template)*


|               |                                            |
|---------------|--------------------------------------------|
| **Created**   | 2023-01-15                                 |
| **Status**    | WIP\| **InReview** \| Approved \| Obsolete |
| **Owner**     | *ramizpolic*                               |
| **Approvers** | *github handles*                           |

---

This RFC proposes the API extensions by allowing multiple Assets to be referenced in Findings in order to improve efficiency and achieve parity with the existing features.

## Background

The existing [API specifications](https://github.com/openclarity/vmclarity/blob/9aa03a8abe22ebddb841a9c28f7a9629f744ced7/api/openapi.yaml#L3395-L3444) only allow a single `Asset` to be referenced for a given `Finding`.
In principle, each `Finding` is described using `(findingInfo, AssetRelationship)` pair.
Additionally, the [database logic](https://github.com/openclarity/vmclarity/blob/9aa03a8abe22ebddb841a9c28f7a9629f744ced7/pkg/apiserver/database/gorm/finding.go#L103-L105) treats every finding that needs to be created as unique.
Together, this can introduce suboptimal behavior for multiple reasons:

- The lack of uniqueness check completely ignores the actual finding contents.
- The way findings API is structured creates an unwanted dependency against assets.
  This also leads to data duplication resulting in performance and memory utilization overheads.
- The lack of aggregation can introduce usage and understanding complexities, e.g. by overloading the data shown on the UI.

## Proposal

The finding-related components can be changed in several ways to address the concerns above.
This has been divided into two categories to more easily understand the scope of proposed changes.

#### Addressing uniqueness

The finding database logic can **implement uniqueness check** similar to the existing logic as exemplified [here](https://github.com/openclarity/vmclarity/blob/9aa03a8abe22ebddb841a9c28f7a9629f744ced7/pkg/apiserver/database/gorm/asset.go#L289).
Moreover, the data required for the check can be extracted directly from the underlying object such as `VulnerabilityFindingInfo`.
In turn, this also enables validation and support for downstream operations such as custom flows within the controller.

The relationship between findings and assets can be fully ignored in checks to avoid creating dependencies.
However, this creates an issue when trying to create a `Finding` with the same underlying data but a different `Asset`.
This is addressed in the following section.

#### Addressing API changes

Findings can be extended to use a **list of assets** instead of referencing a single one as defined in the [API specs](https://github.com/openclarity/vmclarity/blob/2681efa7b5bd1009e9cf740d430587ef7f06ebb7/api/openapi.yaml#L3412).
This ensures that the same finding can be discovered on multiple assets without having to duplicate the data.
Combined, these changes address the performance and memory utilization issues while also enabling aggregation methods.

Alternatively, the `Finding` model can also be extended by adding `AssetFinding` and `AssetFindingRelationship`.
This addresses the issue when many assets contain the same finding.
The implementation, although slightly more complex, would be more efficient as the number of assets grows for a given finding.
_This remains an open question on how to address the `Asset-Finding` relationship._

#### Non-goals

This RFC does not intend to propose changes regarding the relationship of findings to other models like asset scans.
In the context of this proposal, it is assumed that the `AssetScanRelationship` specified in the `Finding.foundBy` property denotes the first `AssetScan` that discovered a given `Finding`.
This behavior, if required, can be addressed later.

### Abandoned Ideas (Optional)

Adding the aggregation methods to the `uibackend` API was considered but abandoned as it does not address the data duplication issue.

---

## Implementation

### 1. Update `Finding` API model with the proposed changes

_Option 1_ - simple finding extension

```yaml
Finding:
  type: object
  allOf:
    - $ref: '#/components/schemas/Metadata'
    - type: object
      properties:
        id:
          type: string
        assets:
          type: array
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

_Option 1_ - more verbose models

Alternatively, depending on the selected approach, `AssetFinding` and `AssetFindingRelationship` can be implemented.
The model should then use a list of `AssetFindingRelationship` to express the relationship of `Assets` for a given `Finding`.

### 2. Handle database-schema changes

The API changes impact the database schema defined in `pkg/apiserver/database/gorm/odata.go` and should be handled accordingly.
In addition, the database-related logic such as bootstrapping needs to be updated to reflect these changes.

### 3. Add uniqueness checks and update related logic

The uniqueness check can be added similarly to the existing implementations. An example is given below.

```go
// file: pkg/apiserver/database/gorm/finding.go

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

### 4. Update related `Finding` UI components

## UX

This RFC has no visible impacts on the UX.

## UI

This RFC changes the following UI components:
- `/findings/{findingType}` table drops the _Asset Name_ and _Asset location_ columns.
  Instead, they are replaced with a single _Assets_ column that shows the number of assets related to a finding.
- `/findings/{findingType}/{findingID}` removes _Asset details_ menu item.
  It is replaced with _See assets_ button on finding summary.
  The button redirects to the `/assets` page and displays filtered, finding-related assets by using OData queries.
  - Note: This can be done similarly to displaying the related _Asset Scans_ in `/scans/scans/{scanID}`.
    Alternatively, a table can be directly used instead of redirects, but might require more changes.
