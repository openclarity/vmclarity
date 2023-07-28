package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"k8s.io/client-go/util/jsonpath"
)

func PrintJSONData(data interface{}, fields string) error {
	// If jsonpath is not set it will print the whole data as json format.
	if fields == "" {
		dataB, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %v", err)
		}
		fmt.Printf("asset: %s", string(dataB))
		return nil
	}
	j := jsonpath.New("parser")
	if err := j.Parse(fields); err != nil {
		return fmt.Errorf("failed to parse jsonpath: %v", err)
	}
	err := j.Execute(os.Stdout, data)
	if err != nil {
		return fmt.Errorf("failed to execute jsonpath: %v", err)
	}
	return nil
}
