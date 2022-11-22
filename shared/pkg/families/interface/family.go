package _interface

type Family interface {
	Run(ResultsGetter) (IsResults, error)
}
