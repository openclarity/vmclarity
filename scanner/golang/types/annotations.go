package types

func AnnotationsAsMap(a *Annotations) map[string]string {
	if a == nil {
		return nil
	}

	slice := *a
	result := make(map[string]string)
	for idx := range *a {
		if slice[idx].Key != nil && slice[idx].Value != nil {
			result[*slice[idx].Key] = *slice[idx].Value
		}
	}
	return result
}
