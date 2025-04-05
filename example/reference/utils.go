package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

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

	errInvalidType := schema.InternalServerError(
		fmt.Sprintf("cannot compare values with different types: %s <> %s", kindV1, kindV2),
		nil,
	)
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

func evalExtraction(functionName string, value any) (any, error) {
	if functionName == "" {
		return value, nil
	}

	valueDateTime, err := utils.DecodeDateTime(value)
	if err != nil {
		return nil, schema.UnprocessableContentError(
			"failed to extract value to date time: "+err.Error(),
			nil,
		)
	}

	switch functionName {
	case "year":
		return valueDateTime.Year(), nil
	case "month":
		return valueDateTime.Month(), nil
	case "day":
		return valueDateTime.Day(), nil
	default:
		return nil, schema.UnprocessableContentError(
			"unsupported extraction function name "+functionName,
			nil,
		)
	}
}

func evalRowFieldPath(fieldPath []string, row map[string]any) (map[string]any, error) {
	if len(fieldPath) == 0 {
		return row, nil
	}

	for _, name := range fieldPath {
		child, ok := row[name]
		if !ok {
			return nil, fmt.Errorf("invalid row field path %s", name)
		}

		row, ok = child.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object when navigating row field path, got %v", child)
		}
	}

	return row, nil
}

func evalComparisonOperator( //nolint:gocyclo,cyclop,funlen
	operator string,
	leftVal any,
	rightValues []any,
) (bool, error) {
	switch operator {
	case "eq":
		for _, rightVal := range rightValues {
			if isEqual(leftVal, rightVal) {
				return true, nil
			}
		}

		return false, nil
	case "in":
		for _, rightVal := range rightValues {
			items, ok := rightVal.([]any)
			if !ok {
				return false, fmt.Errorf("expected array, got %v", rightVal)
			}

			for _, rv := range items {
				if isEqual(leftVal, rv) {
					return true, nil
				}
			}
		}

		return false, nil
	case "lt", "lte", "gt", "gte": //nolint:goconst
		if leftStr, ok := leftVal.(string); ok {
			for _, rightVal := range rightValues {
				rightString, err := utils.DecodeString(rightVal)
				if err != nil {
					return false, err
				}

				var isValid bool

				switch operator {
				case "lt":
					isValid = strings.Compare(leftStr, rightString) < 0
				case "lte":
					isValid = strings.Compare(leftStr, rightString) <= 0
				case "gt":
					isValid = strings.Compare(leftStr, rightString) > 0
				case "gte":
					isValid = strings.Compare(leftStr, rightString) >= 0
				}

				if isValid {
					return true, nil
				}
			}

			return false, nil
		}

		leftNum, err := utils.DecodeFloat[float64](leftVal)
		if err != nil {
			return false, err
		}

		for _, rightVal := range rightValues {
			rightNum, err := utils.DecodeFloat[float64](rightVal)
			if err != nil {
				return false, err
			}

			var isValid bool

			switch operator {
			case "lt":
				isValid = leftNum < rightNum
			case "lte":
				isValid = leftNum <= rightNum
			case "gt":
				isValid = leftNum > rightNum
			case "gte":
				isValid = leftNum >= rightNum
			}

			if isValid {
				return true, nil
			}
		}

		return false, nil
	case "like":
		columnStr, ok := leftVal.(string)
		if !ok {
			return false, fmt.Errorf(
				"failed to compare text values, expected string, got: %v",
				leftVal,
			)
		}

		for _, rawRegex := range rightValues {
			regexStr, ok := rawRegex.(string)
			if !ok {
				return false, schema.UnprocessableContentError(
					fmt.Sprintf("invalid regular expression, got %+v", rawRegex),
					nil,
				)
			}

			regex, err := regexp.Compile(regexStr)
			if err != nil {
				return false, schema.UnprocessableContentError(
					fmt.Sprintf("invalid regular expression: %s", err),
					nil,
				)
			}

			if regex.MatchString(columnStr) {
				return true, nil
			}
		}

		return false, nil
	case "contains", "icontains", "starts_with", "istarts_with", "ends_with", "iends_with":
		leftStr, ok := leftVal.(string)
		if !ok {
			return false, fmt.Errorf(
				"comparison operator %s is only supported on strings, got: %v",
				operator,
				leftVal,
			)
		}

		for _, rightVal := range rightValues {
			rightStr, ok := rightVal.(string)
			if !ok {
				return false, schema.UnprocessableContentError(
					fmt.Sprintf("value is not a string, got %+v", rightVal),
					nil,
				)
			}

			var isValid bool

			switch operator {
			case "contains":
				isValid = strings.Contains(leftStr, rightStr)
			case "icontains":
				isValid = strings.Contains(strings.ToLower(leftStr), strings.ToLower(rightStr))
			case "starts_with":
				isValid = strings.HasPrefix(leftStr, rightStr)
			case "istarts_with":
				isValid = strings.HasPrefix(strings.ToLower(leftStr), strings.ToLower(rightStr))
			case "ends_with":
				isValid = strings.HasSuffix(leftStr, rightStr)
			case "iends_with":
				isValid = strings.HasSuffix(strings.ToLower(leftStr), strings.ToLower(rightStr))
			}

			if isValid {
				return true, nil
			}
		}

		return false, nil
	default:
		return false, schema.UnprocessableContentError(
			"invalid comparison operator: "+operator,
			nil,
		)
	}
}

func evalAggregateFunction(function string, values []any) (*int, error) {
	if len(values) == 0 {
		return nil, nil
	}

	var intValues []int

	for _, value := range values {
		switch v := value.(type) {
		case int:
			intValues = append(intValues, v)
		case int16:
			intValues = append(intValues, int(v))
		case int32:
			intValues = append(intValues, int(v))
		case int64:
			intValues = append(intValues, int(v))
		default:
			return nil, schema.UnprocessableContentError(fmt.Sprintf("%s: column is not an integer, got %+v", function, reflect.ValueOf(v).Kind()), nil)
		}
	}

	sort.Ints(intValues)

	switch function {
	case "min":
		return &intValues[0], nil
	case "max":
		return &intValues[len(intValues)-1], nil
	default:
		return nil, schema.UnprocessableContentError(function+": invalid aggregation function", nil)
	}
}

func isEqual(leftVal, rightVal any) bool {
	if leftVal == rightVal || (leftVal == nil && rightVal == nil) {
		return true
	}

	leftReflectValue, leftNotNull := utils.UnwrapPointerFromAnyToReflectValue(leftVal)
	rightReflectValue, rightNotNull := utils.UnwrapPointerFromAnyToReflectValue(rightVal)

	if !leftNotNull && !rightNotNull {
		return true
	}

	if !leftNotNull || !rightNotNull {
		return false
	}

	leftKind := leftReflectValue.Kind()
	rightKind := rightReflectValue.Kind()

	if leftKind == rightKind {
		return reflect.DeepEqual(leftReflectValue.Interface(), rightReflectValue.Interface())
	}

	switch leftKind {
	case reflect.String:
		rightStr, err := utils.DecodeString(rightReflectValue.Interface())

		return err == nil && leftReflectValue.String() == rightStr
	case reflect.Bool:
		rightBool, err := utils.DecodeBoolean(rightReflectValue.Interface())

		return err == nil && leftReflectValue.Bool() == rightBool
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		leftInt, _ := utils.DecodeInt[int64](leftReflectValue.Interface())
		rightInt, err := utils.DecodeInt[int64](rightReflectValue.Interface())

		return err == nil && leftInt == rightInt
	case reflect.Float32, reflect.Float64:
		leftInt := leftReflectValue.Float()
		rightInt, err := utils.DecodeFloat[float64](rightReflectValue.Interface())

		return err == nil && leftInt == rightInt
	default:
		return false
	}
}
