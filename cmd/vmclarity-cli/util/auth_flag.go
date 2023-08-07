package util

import "github.com/spf13/pflag"

var BearerTokenEnvVar = ""

func RegisterBearerTokenEnvVarFlag(flags *pflag.FlagSet) {
	flags.StringVar(&BearerTokenEnvVar,
		"bearer-token-env",
		"",
		"Set API server authentication bearer token env variable to use for auth data request injection")
}
