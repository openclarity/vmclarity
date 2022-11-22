package _interface

import "github.com/openclarity/vmclarity/shared/pkg/families/types"

type IsResults interface {
	IsResults()
}

type ResultsGetter interface {
	GetResults(familyType types.FamilyType) IsResults
}
