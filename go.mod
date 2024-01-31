module github.com/openclarity/vmclarity

go 1.21.4

require github.com/openclarity/vmclarity/utils v0.0.0-00010101000000-000000000000

// NOTE(chrisgacsal): remove this when the following PR is merged and new helm version is released:
// https://github.com/helm/helm/pull/12310
replace helm.sh/helm/v3 => github.com/zregvart/helm/v3 v3.0.0-20240102124916-a62313e07d76

replace (
	github.com/openclarity/vmclarity/api/client => ./api/client
	github.com/openclarity/vmclarity/api/types => ./api/types
	github.com/openclarity/vmclarity/cli => ./cli
	github.com/openclarity/vmclarity/containerruntimediscovery/client => ./containerruntimediscovery/client
	github.com/openclarity/vmclarity/containerruntimediscovery/types => ./containerruntimediscovery/types
	github.com/openclarity/vmclarity/orchestrator => ./orchestrator
	github.com/openclarity/vmclarity/utils => ./utils
)

// NOTE(akijakya): replace is required for the following issue: https://github.com/mitchellh/mapstructure/issues/327,
// which has been solved in the go-viper fork.
// Remove replace if all packages using the original repo has been switched to this fork (or at least viper:
// https://github.com/spf13/viper/pull/1723)
replace github.com/mitchellh/mapstructure => github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1
