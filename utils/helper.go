package utils

// GetDefault returns the value or default one if value is empty
func GetDefault[T comparable](value T, defaultValue T) T {
	var empty T
	if value == empty {
		return defaultValue
	}
	return value
}

// GetDefaultPtr returns the first pointer or default one if GetDefaultPtr is nil
func GetDefaultPtr[T any](value *T, defaultValue *T) *T {
	if value == nil {
		return defaultValue
	}
	return value
}

// GetDefaultValuePtr return the value of pointer or default one if the value of pointer is null or empty
func GetDefaultValuePtr[T comparable](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	var empty T
	if *value == empty {
		return defaultValue
	}
	return *value
}
