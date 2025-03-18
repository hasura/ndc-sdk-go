package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

func evalColumn(variables map[string]any, row map[string]any, columnName string, arguments map[string]schema.Argument) (any, error) {
	column, ok := row[columnName]
	if !ok {
		return nil, fmt.Errorf("invalid column name %s", columnName)
	}

	if utils.IsNil(column) {
		return nil, nil
	}

	values, ok := column.([]any)
	if !ok {
		return column, nil
	}

	if len(values) == 0 {
		return values, nil
	}

	limitArgument, ok := arguments["limit"]
	if !ok {
		return nil, schema.UnprocessableContentError(fmt.Sprintf("Expected argument 'limit' in column %s", columnName), nil)
	}

	rawLimitArg, err := utils.ResolveArgument(limitArgument, variables)
	if err != nil {
		return nil, err
	}

	limitArg, err := utils.DecodeNullableInt[int](rawLimitArg)
	if err != nil {
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}

	if limitArg != nil && *limitArg < len(values) {
		values = values[:*limitArg]
	}

	return values, nil
}

func evalColumnFieldPath(variables map[string]any, row map[string]any, columnName string, fieldPath []string, arguments map[string]schema.Argument) (any, error) {
	columnValue, err := evalColumn(variables, row, columnName, arguments)
	if err != nil {
		return nil, err
	}

	return evalFieldPath(fieldPath, columnValue)
}

func evalFieldPath(path []string, value any) (any, error) {
	if len(path) == 0 || utils.IsNil(value) {
		return value, nil
	}

	values, ok := value.([]any)
	if ok {
		results := make([]any, len(values))

		for i, v := range values {
			result, err := evalFieldPath(path, v)
			if err != nil {
				return nil, err
			}

			results[i] = result
		}

		return results, nil
	}

	obj, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("invalid field path")
	}

	childValue, ok := obj[path[0]]
	if !ok {
		return nil, errors.New("invalid field path")
	}

	return evalFieldPath(path[1:], childValue)
}

func paginate[R any](collection []R, limit *int, offset *int) []R {
	var start int
	if offset != nil {
		start = *offset
	}
	collectionLength := len(collection)
	if collectionLength <= start {
		return nil
	}
	if limit == nil {
		return collection[start:]
	}
	return collection[start:int(math.Min(float64(collectionLength), float64(start+*limit)))]
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func compare(v1 any, v2 any) (int, error) {
	if v1 == v2 || (v1 == nil && v2 == nil) {
		return 0, nil
	}
	if v1 == nil {
		return -1, nil
	}
	if v2 == nil {
		return 1, nil
	}

	value1 := reflect.ValueOf(v1)
	kindV1 := value1.Kind()
	value2 := reflect.ValueOf(v2)
	kindV2 := value2.Kind()

	errInvalidType := schema.InternalServerError(fmt.Sprintf("cannot compare values with different types: %s <> %s", kindV1, kindV2), nil)
	if kindV1 != kindV2 {
		return 0, errInvalidType
	}

	if kindV1 == reflect.Pointer {
		return compare(value1.Elem().Interface(), value2.Elem().Interface())
	}

	switch value1 := v1.(type) {
	case bool:
		value2, ok := v2.(bool)
		if !ok {
			return 0, errInvalidType
		}
		return boolToInt(value1) - boolToInt(value2), nil
	case int:
		value2, ok := v2.(int)
		if !ok {
			return 0, errInvalidType
		}
		return value1 - value2, nil
	case int8:
		value2, ok := v2.(int8)
		if !ok {
			return 0, errInvalidType
		}
		return int(value1 - value2), nil
	case int16:
		value2, ok := v2.(int16)
		if !ok {
			return 0, errInvalidType
		}
		return int(value1 - value2), nil
	case int32:
		value2, ok := v2.(int32)
		if !ok {
			return 0, errInvalidType
		}
		return int(value1 - value2), nil
	case int64:
		value2, ok := v2.(int64)

		if !ok {
			return 0, errInvalidType
		}
		return int(value1 - value2), nil
	case string:
		value2, ok := v2.(string)
		if !ok {
			return 0, errInvalidType
		}
		return strings.Compare(value1, value2), nil
	default:
		rawV1, err := json.Marshal(v1)
		if err != nil {
			return 0, errInvalidType
		}
		return 0, schema.InternalServerError(fmt.Sprintf("cannot compare values with type: %s, value: %s", kindV1, string(rawV1)), nil)
	}
}
