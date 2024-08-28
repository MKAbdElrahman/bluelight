package v1

import (
	"net/http"
	"strconv"
	"strings"
)

func parseIdFromPath(r *http.Request) (int64, error) {
	idFromPath := r.PathValue("id")
	parsedId, err := strconv.ParseInt(idFromPath, 10, 64)
	return parsedId, err
}

// GetQueryParam retrieves a query parameter as a string.
func GetQueryParam(r *http.Request, key string) string {
	values := r.URL.Query()
	return values.Get(key)
}

// GetQueryParamInt retrieves a query parameter as an integer.
func GetQueryParamInt(r *http.Request, key string) (int, error) {
	values := r.URL.Query()
	if value := values.Get(key); value != "" {
		return strconv.Atoi(value)
	}
	return 0, nil
}

// GetQueryParamBool retrieves a query parameter as a boolean.
func GetQueryParamBool(r *http.Request, key string) (bool, error) {
	values := r.URL.Query()
	if value := values.Get(key); value != "" {
		return strconv.ParseBool(value)
	}
	return false, nil
}

// GetQueryParamSlice retrieves a query parameter as a slice of strings.
func GetQueryParamSlice(r *http.Request, key string, separator string) []string {
	values := r.URL.Query()
	if value := values.Get(key); value != "" {
		return strings.Split(value, separator)
	}
	return nil
}

// GetQueryParamSliceInt retrieves a query parameter as a slice of integers.
func GetQueryParamSliceInt(r *http.Request, key string, separator string) ([]int, error) {
	values := r.URL.Query()
	if value := values.Get(key); value != "" {
		strValues := strings.Split(value, separator)
		intValues := make([]int, len(strValues))
		for i, strValue := range strValues {
			intValue, err := strconv.Atoi(strValue)
			if err != nil {
				return nil, err
			}
			intValues[i] = intValue
		}
		return intValues, nil
	}
	return nil, nil
}
