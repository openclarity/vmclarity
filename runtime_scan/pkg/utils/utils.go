package utils

func StringPtr(val string) *string {
	ret := val
	return &ret
}

func BoolPtr(val bool) *bool {
	ret := val
	return &ret
}

func Int32Ptr(val int32) *int32 {
	ret := val
	return &ret
}
