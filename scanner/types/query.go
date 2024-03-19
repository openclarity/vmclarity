package types

import (
	"fmt"
	"strings"
)

const MetaSelectorSeparator = "="

func MetaSelectorsToMap(m *MetaSelectors) *map[string]string {
	if m == nil {
		return nil
	}

	result := make(map[string]string)
	for _, selector := range *m {
		parts := strings.SplitN(selector, MetaSelectorSeparator, 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		result[key] = value
	}

	return &result
}

func MetaSelectorsFromMap(data map[string]string) *MetaSelectors {
	if len(data) == 0 {
		return nil
	}

	result := MetaSelectors{}
	for key, value := range data {
		key, value := key, value
		result = append(result,
			fmt.Sprintf("%s%s%s", key, MetaSelectorSeparator, value))
	}

	return &result
}
