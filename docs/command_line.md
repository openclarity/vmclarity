# Initiate scan using the cli:

## Reporting results into file:
```
./cli/bin/vmclarity-cli scan --config ~/testConf.yaml -o outputfile
```

If we want to report results to the VMClarity backend, we need to create asset and asset scan object before scan because it requires asset-scan-id

## Reporting results to VMClarity backand:

```
ASSET_ID=$(./cli/bin/vmclarity-cli asset-create --from-json-file assets/dir-asset.json --server http://localhost:8888/api)
ASSET_SCAN_ID=$(./cli/bin/vmclarity-cli asset-scan-create --asset-id $ASSET_ID --server http://localhost:8888/api)
./cli/bin/vmclarity-cli scan --config ~/testConf.yaml --server http://localhost:8888/api --asset-scan-id $ASSET_SCAN_ID
```

Using one-liner:
```
./cli/bin/vmclarity-cli asset-create --from-json-file assets/dir-asset.json --server http://localhost:8888/api | xargs -I{} ./cli/bin/vmclarity-cli asset-scan-create --asset-id {} --server http://localhost:8888/api | xargs -I{} ./cli/bin/vmclarity-cli scan --config ~/testConf.yaml --server http://localhost:8888/api --asset-scan-id {}
```