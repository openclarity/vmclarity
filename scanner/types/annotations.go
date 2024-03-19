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

func AnnotationsFromMap(data map[string]string) *Annotations {
	if len(data) == 0 {
		return nil
	}

	result := Annotations{}
	for k, v := range data {
		k, v := k, v
		result = append(result, Annotation{
			Key:   &k,
			Value: &v,
		})
	}

	return &result
}
