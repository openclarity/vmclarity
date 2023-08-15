package assetscanestimationwatcher

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

func (w *Watcher) getLatestAssetScanStats(ctx context.Context, asset *models.Asset) (models.AssetScanStats, error) {
	var stats models.AssetScanStats

	filterTmpl := "asset/id eq %s and status/general/state eq 'Done' and status/general/errors eq null and scanFamiliesConfig/%s/enabled eq true"
	// get the latest asset scan with exploits enabled
	params := models.GetAssetScansParams{
		Filter:  utils.PointerTo(fmt.Sprintf(filterTmpl, *asset.Id, "exploits")),
		Top:     utils.PointerTo(1),
		OrderBy: utils.PointerTo("status/general/lastTransitionTime DESC"),
	}
	res, err := w.backend.GetAssetScans(ctx, params)
	if err != nil {
		logrus.Errorf("Failed to get asset scans for exploits. Ommiting stats: %v", err)
	} else {
		stats.Exploits = (*res.Items)[0].Stats.Exploits
	}

	// get the latest asset scan with sbom enabled
	params.Filter = utils.PointerTo(fmt.Sprintf(filterTmpl, *asset.Id, "sbom"))
	res, err = w.backend.GetAssetScans(ctx, params)
	if err != nil {
		logrus.Errorf("Failed to get asset scans for sbom. Ommiting stats: %v", err)
	} else {
		stats.Sbom = (*res.Items)[0].Stats.Sbom
	}

	// get the latest asset scan with vulnerability enabled
	params.Filter = utils.PointerTo(fmt.Sprintf(filterTmpl, *asset.Id, "vulnerabilities"))
	res, err = w.backend.GetAssetScans(ctx, params)
	if err != nil {
		logrus.Errorf("Failed to get asset scans for vulnerabilities. Ommiting stats: %v", err)
	} else {
		stats.Vulnerabilities = (*res.Items)[0].Stats.Vulnerabilities
	}

	// get the latest asset scan with malware enabled
	params.Filter = utils.PointerTo(fmt.Sprintf(filterTmpl, *asset.Id, "malware"))
	res, err = w.backend.GetAssetScans(ctx, params)
	if err != nil {
		logrus.Errorf("Failed to get asset scans for malware. Ommiting stats: %v", err)
	} else {
		stats.Malware = (*res.Items)[0].Stats.Malware
	}

	// get the latest asset scan with misconfiguration enabled
	params.Filter = utils.PointerTo(fmt.Sprintf(filterTmpl, *asset.Id, "misconfigurations"))
	res, err = w.backend.GetAssetScans(ctx, params)
	if err != nil {
		logrus.Errorf("Failed to get asset scans for misconfigurations. Ommiting stats: %v", err)
	} else {
		stats.Misconfigurations = (*res.Items)[0].Stats.Misconfigurations
	}

	// get the latest asset scan with rootkits enabled
	params.Filter = utils.PointerTo(fmt.Sprintf(filterTmpl, *asset.Id, "rootkits"))
	res, err = w.backend.GetAssetScans(ctx, params)
	if err != nil {
		logrus.Errorf("Failed to get asset scans for rootkits. Ommiting stats: %v", err)
	} else {
		stats.Rootkits = (*res.Items)[0].Stats.Rootkits
	}

	// get the latest asset scan with secrets enabled
	params.Filter = utils.PointerTo(fmt.Sprintf(filterTmpl, *asset.Id, "secrets"))
	res, err = w.backend.GetAssetScans(ctx, params)
	if err != nil {
		logrus.Errorf("Failed to get asset scans for secrets. Ommiting stats: %v", err)
	} else {
		stats.Secrets = (*res.Items)[0].Stats.Secrets
	}

	return stats, nil
}
