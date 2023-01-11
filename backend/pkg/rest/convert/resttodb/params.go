package resttodb

import (
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

func ConvertGetTargetsParams(params models.GetTargetsParams) database.GetTargetsParams {
	return database.GetTargetsParams{
		Filter:   params.Filter,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
}

func ConvertGetScanResultsParams(params models.GetScanResultsParams) database.GetScanResultsParams {
	return database.GetScanResultsParams{
		Filter:   params.Filter,
		Select:   params.Select,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
}

func ConvertGetScanResultsScanResultIDParams(params models.GetScanResultsScanResultIDParams) database.GetScanResultsScanResultIDParams {
	return database.GetScanResultsScanResultIDParams{
		Select: params.Select,
	}
}

func ConvertGetScansParams(params models.GetScansParams) database.GetScansParams {
	return database.GetScansParams{
		Filter:   params.Filter,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
}

func ConvertGetScanConfigsParams(params models.GetScanConfigsParams) database.GetScanConfigsParams {
	return database.GetScanConfigsParams{
		Filter:   params.Filter,
		Page:     params.Page,
		PageSize: params.PageSize,
	}
}
