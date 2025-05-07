package utils

import "net/url"

func GetOrDefault(query url.Values, key string, defaultVal string) string {
	val := query.Get(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func GetOrNil(query url.Values, key string) *string {
	val := query.Get(key)
	if val == "" {
		return nil
	}
	return &val
}
