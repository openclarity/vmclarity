package utils

import "encoding/json"

// JsonPrettyPrint used for testing.
func JsonPrettyPrint(got any) string {
	jsonResults, err := json.MarshalIndent(got, "", "    ")
	panic(err)
	return string(jsonResults)
}
