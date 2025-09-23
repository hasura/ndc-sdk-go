package schema

import (
	"encoding/json"
	"fmt"
	"slices"
)

// RelationalExpressionType represents a relational expression type enum.
type RelationalExpressionType string

const (
	RelationalExpressionTypeLiteral              RelationalExpressionType = "literal"
	RelationalExpressionTypeColumn               RelationalExpressionType = "column"
	RelationalExpressionTypeCase                 RelationalExpressionType = "case"
	RelationalExpressionTypeAnd                  RelationalExpressionType = "and"
	RelationalExpressionTypeOr                   RelationalExpressionType = "or"
	RelationalExpressionTypeNot                  RelationalExpressionType = "not"
	RelationalExpressionTypeEq                   RelationalExpressionType = "eq"
	RelationalExpressionTypeNotEq                RelationalExpressionType = "not_eq"
	RelationalExpressionTypeIsDistinctFrom       RelationalExpressionType = "is_distinct_from"
	RelationalExpressionTypeIsNotDistinctFrom    RelationalExpressionType = "is_not_distinct_from"
	RelationalExpressionTypeLt                   RelationalExpressionType = "lt"
	RelationalExpressionTypeLtEq                 RelationalExpressionType = "lt_eq"
	RelationalExpressionTypeGt                   RelationalExpressionType = "gt"
	RelationalExpressionTypeGtEq                 RelationalExpressionType = "gt_eq"
	RelationalExpressionTypeIsNotNull            RelationalExpressionType = "is_not_null"
	RelationalExpressionTypeIsNull               RelationalExpressionType = "is_null"
	RelationalExpressionTypeIsTrue               RelationalExpressionType = "is_true"
	RelationalExpressionTypeIsFalse              RelationalExpressionType = "is_false"
	RelationalExpressionTypeIsNotTrue            RelationalExpressionType = "is_not_true"
	RelationalExpressionTypeIsNotFalse           RelationalExpressionType = "is_not_false"
	RelationalExpressionTypeIn                   RelationalExpressionType = "in"
	RelationalExpressionTypeNotIn                RelationalExpressionType = "not_in"
	RelationalExpressionTypeLike                 RelationalExpressionType = "like"
	RelationalExpressionTypeNotLike              RelationalExpressionType = "not_like"
	RelationalExpressionTypeILike                RelationalExpressionType = "i_like"
	RelationalExpressionTypeNotILike             RelationalExpressionType = "not_i_like"
	RelationalExpressionTypeBetween              RelationalExpressionType = "between"
	RelationalExpressionTypeNotBetween           RelationalExpressionType = "not_between"
	RelationalExpressionTypeContains             RelationalExpressionType = "contains"
	RelationalExpressionTypeIsNaN                RelationalExpressionType = "is_na_n"
	RelationalExpressionTypeIsZero               RelationalExpressionType = "is_zero"
	RelationalExpressionTypePlus                 RelationalExpressionType = "plus"
	RelationalExpressionTypeMinus                RelationalExpressionType = "minus"
	RelationalExpressionTypeMultiply             RelationalExpressionType = "multiply"
	RelationalExpressionTypeDivide               RelationalExpressionType = "divide"
	RelationalExpressionTypeModulo               RelationalExpressionType = "modulo"
	RelationalExpressionTypeNegate               RelationalExpressionType = "negate"
	RelationalExpressionTypeCast                 RelationalExpressionType = "cast"
	RelationalExpressionTypeTryCast              RelationalExpressionType = "try_cast"
	RelationalExpressionTypeAbs                  RelationalExpressionType = "abs"
	RelationalExpressionTypeArrayElement         RelationalExpressionType = "array_element"
	RelationalExpressionTypeBTrim                RelationalExpressionType = "b_trim"
	RelationalExpressionTypeCeil                 RelationalExpressionType = "ceil"
	RelationalExpressionTypeCharacterLength      RelationalExpressionType = "character_length"
	RelationalExpressionTypeCoalesce             RelationalExpressionType = "coalesce"
	RelationalExpressionTypeConcat               RelationalExpressionType = "concat"
	RelationalExpressionTypeCos                  RelationalExpressionType = "cos"
	RelationalExpressionTypeCurrentDate          RelationalExpressionType = "current_date"
	RelationalExpressionTypeCurrentTime          RelationalExpressionType = "current_time"
	RelationalExpressionTypeCurrentTimestamp     RelationalExpressionType = "current_timestamp"
	RelationalExpressionTypeDatePart             RelationalExpressionType = "date_part"
	RelationalExpressionTypeDateTrunc            RelationalExpressionType = "date_trunc"
	RelationalExpressionTypeExp                  RelationalExpressionType = "exp"
	RelationalExpressionTypeFloor                RelationalExpressionType = "floor"
	RelationalExpressionTypeGetField             RelationalExpressionType = "get_field"
	RelationalExpressionTypeGreatest             RelationalExpressionType = "greatest"
	RelationalExpressionTypeLeast                RelationalExpressionType = "least"
	RelationalExpressionTypeLeft                 RelationalExpressionType = "left"
	RelationalExpressionTypeLn                   RelationalExpressionType = "ln"
	RelationalExpressionTypeLog                  RelationalExpressionType = "log"
	RelationalExpressionTypeLog10                RelationalExpressionType = "log10"
	RelationalExpressionTypeLog2                 RelationalExpressionType = "log2"
	RelationalExpressionTypeLPad                 RelationalExpressionType = "l_pad"
	RelationalExpressionTypeLTrim                RelationalExpressionType = "l_trim"
	RelationalExpressionTypeNullIf               RelationalExpressionType = "null_if"
	RelationalExpressionTypeNvl                  RelationalExpressionType = "nvl"
	RelationalExpressionTypePower                RelationalExpressionType = "power"
	RelationalExpressionTypeRandom               RelationalExpressionType = "random"
	RelationalExpressionTypeReplace              RelationalExpressionType = "replace"
	RelationalExpressionTypeReverse              RelationalExpressionType = "reverse"
	RelationalExpressionTypeRight                RelationalExpressionType = "right"
	RelationalExpressionTypeRound                RelationalExpressionType = "round"
	RelationalExpressionTypeRPad                 RelationalExpressionType = "r_pad"
	RelationalExpressionTypeRTrim                RelationalExpressionType = "r_trim"
	RelationalExpressionTypeSqrt                 RelationalExpressionType = "sqrt"
	RelationalExpressionTypeStrPos               RelationalExpressionType = "str_pos"
	RelationalExpressionTypeSubstr               RelationalExpressionType = "substr"
	RelationalExpressionTypeSubstrIndex          RelationalExpressionType = "substr_index"
	RelationalExpressionTypeTan                  RelationalExpressionType = "tan"
	RelationalExpressionTypeToDate               RelationalExpressionType = "to_date"
	RelationalExpressionTypeToTimestamp          RelationalExpressionType = "to_timestamp"
	RelationalExpressionTypeTrunc                RelationalExpressionType = "trunc"
	RelationalExpressionTypeToLower              RelationalExpressionType = "to_lower"
	RelationalExpressionTypeToUpper              RelationalExpressionType = "to_upper"
	RelationalExpressionTypeBinaryConcat         RelationalExpressionType = "binary_concat"
	RelationalExpressionTypeAverage              RelationalExpressionType = "average"
	RelationalExpressionTypeBoolAnd              RelationalExpressionType = "bool_and"
	RelationalExpressionTypeBoolOr               RelationalExpressionType = "bool_or"
	RelationalExpressionTypeCount                RelationalExpressionType = "count"
	RelationalExpressionTypeFirstValue           RelationalExpressionType = "first_value"
	RelationalExpressionTypeLastValue            RelationalExpressionType = "last_value"
	RelationalExpressionTypeMax                  RelationalExpressionType = "max"
	RelationalExpressionTypeMedian               RelationalExpressionType = "median"
	RelationalExpressionTypeMin                  RelationalExpressionType = "min"
	RelationalExpressionTypeStringAgg            RelationalExpressionType = "string_agg"
	RelationalExpressionTypeSum                  RelationalExpressionType = "sum"
	RelationalExpressionTypeVar                  RelationalExpressionType = "var"
	RelationalExpressionTypeStddev               RelationalExpressionType = "stddev"
	RelationalExpressionTypeStddevPop            RelationalExpressionType = "stddev_pop"
	RelationalExpressionTypeApproxPercentileCont RelationalExpressionType = "approx_percentile_cont"
	RelationalExpressionTypeArrayAgg             RelationalExpressionType = "array_agg"
	RelationalExpressionTypeApproxDistinct       RelationalExpressionType = "approx_distinct"
	RelationalExpressionTypeRowNumber            RelationalExpressionType = "row_number"
	RelationalExpressionTypeDenseRank            RelationalExpressionType = "dense_rank"
	RelationalExpressionTypeNTile                RelationalExpressionType = "n_tile"
	RelationalExpressionTypeRank                 RelationalExpressionType = "rank"
	RelationalExpressionTypeCumeDist             RelationalExpressionType = "cume_dist"
	RelationalExpressionTypePercentRank          RelationalExpressionType = "percent_rank"
)

var enumValues_RelationalExpressionType = []RelationalExpressionType{
	RelationalExpressionTypeLiteral,
	RelationalExpressionTypeColumn,
	RelationalExpressionTypeCase,
	RelationalExpressionTypeAnd,
	RelationalExpressionTypeOr,
	RelationalExpressionTypeNot,
	RelationalExpressionTypeEq,
	RelationalExpressionTypeNotEq,
	RelationalExpressionTypeIsDistinctFrom,
	RelationalExpressionTypeIsNotDistinctFrom,
	RelationalExpressionTypeLt,
	RelationalExpressionTypeLtEq,
	RelationalExpressionTypeGt,
	RelationalExpressionTypeGtEq,
	RelationalExpressionTypeIsNotNull,
	RelationalExpressionTypeIsNull,
	RelationalExpressionTypeIsTrue,
	RelationalExpressionTypeIsFalse,
	RelationalExpressionTypeIsNotTrue,
	RelationalExpressionTypeIsNotFalse,
	RelationalExpressionTypeIn,
	RelationalExpressionTypeNotIn,
	RelationalExpressionTypeLike,
	RelationalExpressionTypeNotLike,
	RelationalExpressionTypeILike,
	RelationalExpressionTypeNotILike,
	RelationalExpressionTypeBetween,
	RelationalExpressionTypeNotBetween,
	RelationalExpressionTypeContains,
	RelationalExpressionTypeIsNaN,
	RelationalExpressionTypeIsZero,
	RelationalExpressionTypePlus,
	RelationalExpressionTypeMinus,
	RelationalExpressionTypeMultiply,
	RelationalExpressionTypeDivide,
	RelationalExpressionTypeModulo,
	RelationalExpressionTypeNegate,
	RelationalExpressionTypeCast,
	RelationalExpressionTypeTryCast,
	RelationalExpressionTypeAbs,
	RelationalExpressionTypeArrayElement,
	RelationalExpressionTypeBTrim,
	RelationalExpressionTypeCeil,
	RelationalExpressionTypeCharacterLength,
	RelationalExpressionTypeCoalesce,
	RelationalExpressionTypeConcat,
	RelationalExpressionTypeCos,
	RelationalExpressionTypeCurrentDate,
	RelationalExpressionTypeCurrentTime,
	RelationalExpressionTypeCurrentTimestamp,
	RelationalExpressionTypeDatePart,
	RelationalExpressionTypeDateTrunc,
	RelationalExpressionTypeExp,
	RelationalExpressionTypeFloor,
	RelationalExpressionTypeGetField,
	RelationalExpressionTypeGreatest,
	RelationalExpressionTypeLeast,
	RelationalExpressionTypeLeft,
	RelationalExpressionTypeLn,
	RelationalExpressionTypeLog,
	RelationalExpressionTypeLog10,
	RelationalExpressionTypeLog2,
	RelationalExpressionTypeLPad,
	RelationalExpressionTypeLTrim,
	RelationalExpressionTypeNullIf,
	RelationalExpressionTypeNvl,
	RelationalExpressionTypePower,
	RelationalExpressionTypeRandom,
	RelationalExpressionTypeReplace,
	RelationalExpressionTypeReverse,
	RelationalExpressionTypeRight,
	RelationalExpressionTypeRound,
	RelationalExpressionTypeRPad,
	RelationalExpressionTypeRTrim,
	RelationalExpressionTypeSqrt,
	RelationalExpressionTypeStrPos,
	RelationalExpressionTypeSubstr,
	RelationalExpressionTypeSubstrIndex,
	RelationalExpressionTypeTan,
	RelationalExpressionTypeToDate,
	RelationalExpressionTypeToTimestamp,
	RelationalExpressionTypeTrunc,
	RelationalExpressionTypeToLower,
	RelationalExpressionTypeToUpper,
	RelationalExpressionTypeBinaryConcat,
	RelationalExpressionTypeAverage,
	RelationalExpressionTypeBoolAnd,
	RelationalExpressionTypeBoolOr,
	RelationalExpressionTypeCount,
	RelationalExpressionTypeFirstValue,
	RelationalExpressionTypeLastValue,
	RelationalExpressionTypeMax,
	RelationalExpressionTypeMedian,
	RelationalExpressionTypeMin,
	RelationalExpressionTypeStringAgg,
	RelationalExpressionTypeSum,
	RelationalExpressionTypeVar,
	RelationalExpressionTypeStddev,
	RelationalExpressionTypeStddevPop,
	RelationalExpressionTypeApproxPercentileCont,
	RelationalExpressionTypeArrayAgg,
	RelationalExpressionTypeApproxDistinct,
	RelationalExpressionTypeRowNumber,
	RelationalExpressionTypeDenseRank,
	RelationalExpressionTypeNTile,
	RelationalExpressionTypeRank,
	RelationalExpressionTypeCumeDist,
	RelationalExpressionTypePercentRank,
}

// ParseRelationalExpressionType parses a relational expression type from string.
func ParseRelationalExpressionType(input string) (RelationalExpressionType, error) {
	result := RelationalExpressionType(input)
	if !result.IsValid() {
		return RelationalExpressionType(
				"",
			), fmt.Errorf(
				"failed to parse RelationalExpressionType, expect one of %v, got %s",
				enumValues_RelationalExpressionType,
				input,
			)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j RelationalExpressionType) IsValid() bool {
	return slices.Contains(enumValues_RelationalExpressionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationalExpressionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseRelationalExpressionType(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// RelationalExpression is provided by reference to a relational expression.
type RelationalExpression struct {
	inner RelationalExpressionInner
}

// NewRelationalExpression creates a RelationalExpression instance.
func NewRelationalExpression[T RelationalExpressionInner](inner T) RelationalExpression {
	return RelationalExpression{
		inner: inner,
	}
}

// RelationalExpressionInner abstracts the interface for RelationalExpression.
type RelationalExpressionInner interface {
	Type() RelationalExpressionType
	ToMap() map[string]any
	Wrap() RelationalExpression
}

// IsEmpty checks if the inner type is empty.
func (j RelationalExpression) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j RelationalExpression) Type() RelationalExpressionType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// MarshalJSON implements json.Marshaler interface.
func (j RelationalExpression) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationalExpression) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	if err := json.Unmarshal(b, &rawType); err != nil {
		return err
	}

	ty, err := ParseRelationalExpressionType(rawType.Type)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RelationalExpression: %w", err)
	}

	switch ty {
	case RelationalExpressionTypeAbs:
		var result RelationalExpressionAbs

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionAbs: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeAnd:
		var result RelationalExpressionAnd

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionAnd: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeApproxDistinct:
		var result RelationalExpressionApproxDistinct

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionApproxDistinct: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeApproxPercentileCont:
		var result RelationalExpressionApproxPercentileCont

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal RelationalExpressionApproxPercentileCont: %w",
				err,
			)
		}

		j.inner = &result
	case RelationalExpressionTypeArrayAgg:
		var result RelationalExpressionArrayAgg

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionArrayAgg: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeArrayElement:
		var result RelationalExpressionArrayElement

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionArrayElement: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeAverage:
		var result RelationalExpressionAverage

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionAverage: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeBTrim:
		var result RelationalExpressionBTrim

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionBTrim: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeBetween:
		var result RelationalExpressionBetween

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionBetween: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLastValue:
		var result RelationalExpressionLastValue

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLastValue: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeBinaryConcat:
		var result RelationalExpressionBinaryConcat

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionBinaryConcat: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeBoolAnd:
		var result RelationalExpressionBoolAnd

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionBoolAnd: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeBoolOr:
		var result RelationalExpressionBoolOr

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionBoolOr: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCase:
		var result RelationalExpressionCase

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCase: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCast:
		var result RelationalExpressionCast

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCast: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCeil:
		var result RelationalExpressionCeil

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCeil: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCharacterLength:
		var result RelationalExpressionCharacterLength

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCharacterLength: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCoalesce:
		var result RelationalExpressionCoalesce

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCoalesce: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeColumn:
		var result RelationalExpressionColumn

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionColumn: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeConcat:
		var result RelationalExpressionConcat

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionConcat: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeContains:
		var result RelationalExpressionContains

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionContains: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCos:
		var result RelationalExpressionCos

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCos: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCount:
		var result RelationalExpressionCount

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCount: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCumeDist:
		var result RelationalExpressionCumeDist

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionCumeDist: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeCurrentDate:
		j.inner = NewRelationalExpressionCurrentDate()
	case RelationalExpressionTypeCurrentTime:
		j.inner = NewRelationalExpressionCurrentTime()
	case RelationalExpressionTypeCurrentTimestamp:
		j.inner = NewRelationalExpressionCurrentTimestamp()
	case RelationalExpressionTypeDatePart:
		var result RelationalExpressionDatePart

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionDatePart: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeDateTrunc:
		var result RelationalExpressionDateTrunc

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionDateTrunc: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeDenseRank:
		var result RelationalExpressionDenseRank

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionDenseRank: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeDivide:
		var result RelationalExpressionDivide

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionDivide: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeEq:
		var result RelationalExpressionEq

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionEq: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeExp:
		var result RelationalExpressionExp

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionExp: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeFirstValue:
		var result RelationalExpressionFirstValue

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionFirstValue: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeFloor:
		var result RelationalExpressionFloor

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionFloor: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeGetField:
		var result RelationalExpressionGetField

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionGetField: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeGreatest:
		var result RelationalExpressionGreatest

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionGreatest: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeGt:
		var result RelationalExpressionGt

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionGt: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeGtEq:
		var result RelationalExpressionGtEq

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionGtEq: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeILike:
		var result RelationalExpressionILike

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionILike: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIn:
		var result RelationalExpressionIn

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIn: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsDistinctFrom:
		var result RelationalExpressionIsDistinctFrom

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsDistinctFrom: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsFalse:
		var result RelationalExpressionIsFalse

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsFalse: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsNaN:
		var result RelationalExpressionIsNaN

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsNaN: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsNotDistinctFrom:
		var result RelationalExpressionIsDistinctFrom

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsDistinctFrom: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsNotFalse:
		var result RelationalExpressionIsNotFalse

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsNotFalse: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsNotNull:
		var result RelationalExpressionIsNotNull

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsNotNull: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsNotTrue:
		var result RelationalExpressionIsNotTrue

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsNotTrue: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsNull:
		var result RelationalExpressionIsNull

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsNull: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsTrue:
		var result RelationalExpressionIsTrue

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsTrue: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeIsZero:
		var result RelationalExpressionIsZero

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionIsZero: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLPad:
		var result RelationalExpressionLPad

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLPad: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLTrim:
		var result RelationalExpressionLTrim

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLTrim: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLeast:
		var result RelationalExpressionLeast

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLeast: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLeft:
		var result RelationalExpressionLeft

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLeft: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLike:
		var result RelationalExpressionLike

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLike: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLiteral:
		var result RelationalExpressionLiteral

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLiteral: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLn:
		var result RelationalExpressionLn

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLn: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLog:
		var result RelationalExpressionLog

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLog: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLog10:
		var result RelationalExpressionLog10

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLog10: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLog2:
		var result RelationalExpressionLog2

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLog2: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLt:
		var result RelationalExpressionLt

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLt: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeLtEq:
		var result RelationalExpressionLtEq

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionLtEq: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeMax:
		var result RelationalExpressionMax

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionMax: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeMedian:
		var result RelationalExpressionMedian

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionMedian: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeMin:
		var result RelationalExpressionMin

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionMin: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeMinus:
		var result RelationalExpressionMinus

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionMinus: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeModulo:
		var result RelationalExpressionModulo

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionModulo: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeMultiply:
		var result RelationalExpressionMultiply

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionMultiply: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNTile:
		var result RelationalExpressionNTile

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNTile: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNegate:
		var result RelationalExpressionNegate

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNegate: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNot:
		var result RelationalExpressionNot

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNot: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNotBetween:
		var result RelationalExpressionBetween

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionBetween: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNotEq:
		var result RelationalExpressionNotEq

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNotEq: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNotILike:
		var result RelationalExpressionILike

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionILike: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNotIn:
		var result RelationalExpressionNotIn

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNotIn: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNotLike:
		var result RelationalExpressionNotLike

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNotLike: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNullIf:
		var result RelationalExpressionNullIf

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNullIf: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeNvl:
		var result RelationalExpressionNvl

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionNvl: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeOr:
		var result RelationalExpressionOr

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionOr: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypePercentRank:
		var result RelationalExpressionPercentRank

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionPercentRank: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypePlus:
		var result RelationalExpressionPlus

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionPlus: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypePower:
		var result RelationalExpressionPower

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionPower: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeRPad:
		var result RelationalExpressionRPad

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionRPad: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeRTrim:
		var result RelationalExpressionRTrim

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionRTrim: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeRandom:
		j.inner = NewRelationalExpressionRandom()
	case RelationalExpressionTypeRank:
		var result RelationalExpressionRank

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionRank: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeReplace:
		var result RelationalExpressionReplace

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionReplace: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeReverse:
		var result RelationalExpressionReverse

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionReverse: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeRight:
		var result RelationalExpressionRight

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionRight: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeRound:
		var result RelationalExpressionRound

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionRound: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeRowNumber:
		var result RelationalExpressionRowNumber

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionRowNumber: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeSqrt:
		var result RelationalExpressionSqrt

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionSqrt: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeStddev:
		var result RelationalExpressionStddev

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionStddev: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeStddevPop:
		var result RelationalExpressionStddevPop

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionStddevPop: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeStrPos:
		var result RelationalExpressionStrPos

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionStrPos: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeStringAgg:
		var result RelationalExpressionStringAgg

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionStringAgg: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeSubstr:
		var result RelationalExpressionSubstr

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionSubstr: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeSubstrIndex:
		var result RelationalExpressionSubstrIndex

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionSubstrIndex: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeSum:
		var result RelationalExpressionSum

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionSum: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeTan:
		var result RelationalExpressionTan

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionTan: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeToDate:
		var result RelationalExpressionToDate

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionToDate: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeToLower:
		var result RelationalExpressionToLower

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionToLower: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeToTimestamp:
		var result RelationalExpressionToTimestamp

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionToTimestamp: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeToUpper:
		var result RelationalExpressionToUpper

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionToUpper: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeTrunc:
		var result RelationalExpressionTrunc

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionTrunc: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeTryCast:
		var result RelationalExpressionTryCast

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionTryCast: %w", err)
		}

		j.inner = &result
	case RelationalExpressionTypeVar:
		var result RelationalExpressionVar

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RelationalExpressionVar: %w", err)
		}

		j.inner = &result
	default:
		return fmt.Errorf("unsupported relational expression type: %s", ty)
	}

	return nil
}

// Interface tries to convert the instance to AggregateInner interface.
func (j RelationalExpression) Interface() RelationalExpressionInner {
	return j.inner
}

// RelationalExpressionLiteral represents a RelationalExpression with literal value.
type RelationalExpressionLiteral struct {
	Literal RelationalLiteral `json:"literal" mapstructure:"literal" yaml:"literal"`
}

// NewRelationalExpressionLiteral creates a RelationalExpressionLiteral instance.
func NewRelationalExpressionLiteral[R RelationalLiteralInner](
	literal R,
) *RelationalExpressionLiteral {
	return &RelationalExpressionLiteral{
		Literal: literal.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLiteral) Type() RelationalExpressionType {
	return RelationalExpressionTypeLiteral
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLiteral) ToMap() map[string]any {
	return map[string]any{
		"type":    j.Type(),
		"literal": j.Literal,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLiteral) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionColumn represents a RelationalExpression with column type.
type RelationalExpressionColumn struct {
	Index uint64 `json:"index" mapstructure:"index" yaml:"index"`
}

// NewRelationalExpressionColumn creates a RelationalExpressionColumn instance.
func NewRelationalExpressionColumn(index uint64) *RelationalExpressionColumn {
	return &RelationalExpressionColumn{
		Index: index,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionColumn) Type() RelationalExpressionType {
	return RelationalExpressionTypeColumn
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionColumn) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"index": j.Index,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionColumn) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCase represents a RelationalExpression with CASE type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.conditional.case`
// * During filtering: `relational_query.filter.conditional.case`
// * During sorting:`relational_query.sort.expression.conditional.case`
// * During joining: `relational_query.join.expression.conditional.case`
// * During aggregation: `relational_query.aggregate.expression.conditional.case`
// * During windowing: `relational_query.window.expression.conditional.case`.
type RelationalExpressionCase struct {
	When      []CaseWhen            `json:"when"                mapstructure:"when"                yaml:"when"`
	Scrutinee *RelationalExpression `json:"scrutinee,omitempty" mapstructure:"scrutinee,omitempty" yaml:"scrutinee,omitempty"`
	Default   *RelationalExpression `json:"default,omitempty"   mapstructure:"default,omitempty"   yaml:"default,omitempty"`
}

// NewRelationalExpressionCase creates a RelationalExpressionCase instance.
func NewRelationalExpressionCase(when []CaseWhen) *RelationalExpressionCase {
	return &RelationalExpressionCase{
		When: when,
	}
}

// WithScrutinee returns the RelationalExpressionCase with scrutinee.
func (j RelationalExpressionCase) WithScrutinee(
	scrutinee *RelationalExpression,
) *RelationalExpressionCase {
	j.Scrutinee = scrutinee

	return &j
}

// WithDefault returns the RelationalExpressionCase with default expression.
func (j RelationalExpressionCase) WithDefault(
	defaultExpr *RelationalExpression,
) *RelationalExpressionCase {
	j.Default = defaultExpr

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionCase) Type() RelationalExpressionType {
	return RelationalExpressionTypeCase
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionCase) ToMap() map[string]any {
	result := map[string]any{
		"type":    j.Type(),
		"when":    j.When,
		"default": j.Default,
	}

	if j.Scrutinee != nil {
		result["scrutinee"] = j.Scrutinee
	}

	if j.Default != nil {
		result["default"] = j.Default
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCase) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionAnd represents a RelationalExpression with the and type.
type RelationalExpressionAnd struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionAnd creates a RelationalExpressionAnd instance.
func NewRelationalExpressionAnd[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionAnd {
	return &RelationalExpressionAnd{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionAnd) Type() RelationalExpressionType {
	return RelationalExpressionTypeAnd
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionAnd) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionAnd) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionOr represents a RelationalExpression with the or type.
type RelationalExpressionOr struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionOr creates a RelationalExpressionOr instance.
func NewRelationalExpressionOr[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionOr {
	return &RelationalExpressionOr{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionOr) Type() RelationalExpressionType {
	return RelationalExpressionTypeOr
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionOr) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionOr) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNot represents a RelationalExpression with the not type.
type RelationalExpressionNot struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionNot creates a RelationalExpressionNot instance.
func NewRelationalExpressionNot[E RelationalExpressionInner](expr E) *RelationalExpressionNot {
	return &RelationalExpressionNot{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNot) Type() RelationalExpressionType {
	return RelationalExpressionTypeOr
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNot) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNot) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionEq represents a RelationalExpression with the eq type.
type RelationalExpressionEq struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionEq creates a RelationalExpressionEq instance.
func NewRelationalExpressionEq[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionEq {
	return &RelationalExpressionEq{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionEq) Type() RelationalExpressionType {
	return RelationalExpressionTypeEq
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionEq) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionEq) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNotEq represents a RelationalExpression with the not_eq type.
type RelationalExpressionNotEq struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionNotEq creates a RelationalExpressionNotEq instance.
func NewRelationalExpressionNotEq[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionNotEq {
	return &RelationalExpressionNotEq{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNotEq) Type() RelationalExpressionType {
	return RelationalExpressionTypeNotEq
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNotEq) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNotEq) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsDistinctFrom represents a RelationalExpression with the is_distinct_from type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.is_distinct_from`
// * During filtering: `relational_query.filter.comparison.is_distinct_from`
// * During sorting:`relational_query.sort.expression.comparison.is_distinct_from`
// * During joining: `relational_query.join.expression.comparison.is_distinct_from`
// * During aggregation: `relational_query.aggregate.expression.comparison.is_distinct_from`
// * During windowing: `relational_query.window.expression.comparison.is_distinct_from`.
type RelationalExpressionIsDistinctFrom struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionIsDistinctFrom creates a RelationalExpressionIsDistinctFrom instance.
func NewRelationalExpressionIsDistinctFrom[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionIsDistinctFrom {
	return &RelationalExpressionIsDistinctFrom{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsDistinctFrom) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsDistinctFrom
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsDistinctFrom) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsDistinctFrom) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsNotDistinctFrom represents a RelationalExpression with the is_not_distinct_from type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.is_not_distinct_from`
// * During filtering: `relational_query.filter.comparison.is_not_distinct_from`
// * During sorting:`relational_query.sort.expression.comparison.is_not_distinct_from`
// * During joining: `relational_query.join.expression.comparison.is_not_distinct_from`
// * During aggregation: `relational_query.aggregate.expression.comparison.is_not_distinct_from`
// * During windowing: `relational_query.window.expression.comparison.is_not_distinct_from`.
type RelationalExpressionIsNotDistinctFrom struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionIsNotDistinctFrom creates a RelationalExpressionIsDistinctFrom instance.
func NewRelationalExpressionIsNotDistinctFrom[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionIsNotDistinctFrom {
	return &RelationalExpressionIsNotDistinctFrom{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsNotDistinctFrom) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsNotDistinctFrom
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsNotDistinctFrom) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsNotDistinctFrom) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLt represents a RelationalExpression with the lt type.
type RelationalExpressionLt struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionLt creates a RelationalExpressionLt instance.
func NewRelationalExpressionLt[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionLt {
	return &RelationalExpressionLt{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLt) Type() RelationalExpressionType {
	return RelationalExpressionTypeLt
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLt) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLt) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLtEq represents a RelationalExpression with the lt_eq type.
type RelationalExpressionLtEq struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionLtEq creates a RelationalExpressionLtEq instance.
func NewRelationalExpressionLtEq[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionLtEq {
	return &RelationalExpressionLtEq{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLtEq) Type() RelationalExpressionType {
	return RelationalExpressionTypeLtEq
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLtEq) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLtEq) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionGt represents a RelationalExpression with the gt type.
type RelationalExpressionGt struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionGt creates a RelationalExpressionGt instance.
func NewRelationalExpressionGt[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionGt {
	return &RelationalExpressionGt{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionGt) Type() RelationalExpressionType {
	return RelationalExpressionTypeGt
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionGt) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionGt) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionGtEq represents a RelationalExpression with the gt_eq type.
type RelationalExpressionGtEq struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionGtEq creates a RelationalExpressionGtEq instance.
func NewRelationalExpressionGtEq[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionGtEq {
	return &RelationalExpressionGtEq{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionGtEq) Type() RelationalExpressionType {
	return RelationalExpressionTypeGtEq
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionGtEq) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionGtEq) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsNotNull represents a RelationalExpression with the is_not_null type.
type RelationalExpressionIsNotNull struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsNotNull creates a RelationalExpressionIsNotNull instance.
func NewRelationalExpressionIsNotNull[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsNotNull {
	return &RelationalExpressionIsNotNull{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsNotNull) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsNotNull
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsNotNull) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsNotNull) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsNull represents a RelationalExpression with the is_null type.
type RelationalExpressionIsNull struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsNull creates a RelationalExpressionIsNull instance.
func NewRelationalExpressionIsNull[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsNull {
	return &RelationalExpressionIsNull{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsNull) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsNull
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsNull) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsNull) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsTrue represents a RelationalExpression with the is_true type.
type RelationalExpressionIsTrue struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsTrue creates a RelationalExpressionIsTrue instance.
func NewRelationalExpressionIsTrue[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsTrue {
	return &RelationalExpressionIsTrue{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsTrue) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsTrue
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsTrue) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsTrue) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsFalse represents a RelationalExpression with the is_false type.
type RelationalExpressionIsFalse struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsFalse creates a RelationalExpressionIsFalse instance.
func NewRelationalExpressionIsFalse[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsFalse {
	return &RelationalExpressionIsFalse{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsFalse) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsFalse
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsFalse) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsFalse) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsNotTrue represents a RelationalExpression with the is_not_true type.
type RelationalExpressionIsNotTrue struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsNotTrue creates a RelationalExpressionIsNotTrue instance.
func NewRelationalExpressionIsNotTrue[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsNotTrue {
	return &RelationalExpressionIsNotTrue{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsNotTrue) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsNotTrue
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsNotTrue) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsNotTrue) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsNotFalse represents a RelationalExpression with the is_not_false type.
type RelationalExpressionIsNotFalse struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsNotFalse creates a RelationalExpressionIsNotFalse instance.
func NewRelationalExpressionIsNotFalse[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsNotFalse {
	return &RelationalExpressionIsNotFalse{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsNotFalse) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsNotFalse
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsNotFalse) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsNotFalse) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIn represents a RelationalExpression with the in type.
type RelationalExpressionIn struct {
	Expr RelationalExpression   `json:"expr" mapstructure:"expr" yaml:"expr"`
	List []RelationalExpression `json:"list" mapstructure:"list" yaml:"list"`
}

// NewRelationalExpressionIn creates a RelationalExpressionIn instance.
func NewRelationalExpressionIn[E RelationalExpressionInner](
	expr E,
	list []RelationalExpressionInner,
) *RelationalExpressionIn {
	ls := []RelationalExpression{}

	for _, item := range list {
		if item != nil {
			ls = append(ls, item.Wrap())
		}
	}

	return &RelationalExpressionIn{
		Expr: expr.Wrap(),
		List: ls,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIn) Type() RelationalExpressionType {
	return RelationalExpressionTypeIn
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIn) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
		"list": j.List,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIn) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNotIn represents a RelationalExpression with the not_in type.
type RelationalExpressionNotIn struct {
	Expr RelationalExpression   `json:"expr" mapstructure:"expr" yaml:"expr"`
	List []RelationalExpression `json:"list" mapstructure:"list" yaml:"list"`
}

// NewRelationalExpressionNotIn creates a RelationalExpressionNotIn instance.
func NewRelationalExpressionNotIn[E RelationalExpressionInner](
	expr E,
	list []RelationalExpressionInner,
) *RelationalExpressionNotIn {
	ls := []RelationalExpression{}

	for _, item := range list {
		if item != nil {
			ls = append(ls, item.Wrap())
		}
	}

	return &RelationalExpressionNotIn{
		Expr: expr.Wrap(),
		List: ls,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNotIn) Type() RelationalExpressionType {
	return RelationalExpressionTypeNotIn
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNotIn) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
		"list": j.List,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNotIn) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLike represents a RelationalExpression with the like type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.like`
// * During filtering: `relational_query.filter.comparison.like`
// * During sorting:`relational_query.sort.expression.comparison.like`
// * During joining: `relational_query.join.expression.comparison.like`
// * During aggregation: `relational_query.aggregate.expression.comparison.like`
// * During windowing: `relational_query.window.expression.comparison.like`.
type RelationalExpressionLike struct {
	Expr    RelationalExpression `json:"expr"    mapstructure:"expr"    yaml:"expr"`
	Pattern RelationalExpression `json:"pattern" mapstructure:"pattern" yaml:"pattern"`
}

// NewRelationalExpressionLike creates a RelationalExpressionLike instance.
func NewRelationalExpressionLike[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionLike {
	return &RelationalExpressionLike{
		Expr:    left.Wrap(),
		Pattern: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLike) Type() RelationalExpressionType {
	return RelationalExpressionTypeLike
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLike) ToMap() map[string]any {
	return map[string]any{
		"type":    j.Type(),
		"expr":    j.Expr,
		"pattern": j.Pattern,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLike) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNotLike represents a RelationalExpression with the not_like type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.like`
// * During filtering: `relational_query.filter.comparison.like`
// * During sorting:`relational_query.sort.expression.comparison.like`
// * During joining: `relational_query.join.expression.comparison.like`
// * During aggregation: `relational_query.aggregate.expression.comparison.like`
// * During windowing: `relational_query.window.expression.comparison.like`.
type RelationalExpressionNotLike struct {
	Expr    RelationalExpression `json:"expr"    mapstructure:"expr"    yaml:"expr"`
	Pattern RelationalExpression `json:"pattern" mapstructure:"pattern" yaml:"pattern"`
}

// NewRelationalExpressionNotLike creates a RelationalExpressionNotLike instance.
func NewRelationalExpressionNotLike[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionNotLike {
	return &RelationalExpressionNotLike{
		Expr:    left.Wrap(),
		Pattern: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNotLike) Type() RelationalExpressionType {
	return RelationalExpressionTypeNotLike
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNotLike) ToMap() map[string]any {
	return map[string]any{
		"type":    j.Type(),
		"expr":    j.Expr,
		"pattern": j.Pattern,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNotLike) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionILike represents a RelationalExpression with the i_like type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.ilike`
// * During filtering: `relational_query.filter.comparison.ilike`
// * During sorting:`relational_query.sort.expression.comparison.ilike`
// * During joining: `relational_query.join.expression.comparison.ilike`
// * During aggregation: `relational_query.aggregate.expression.comparison.ilike`
// * During windowing: `relational_query.window.expression.comparison.ilike`.
type RelationalExpressionILike struct {
	Expr    RelationalExpression `json:"expr"    mapstructure:"expr"    yaml:"expr"`
	Pattern RelationalExpression `json:"pattern" mapstructure:"pattern" yaml:"pattern"`
}

// NewRelationalExpressionILike creates a RelationalExpressionILike instance.
func NewRelationalExpressionILike[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionILike {
	return &RelationalExpressionILike{
		Expr:    left.Wrap(),
		Pattern: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionILike) Type() RelationalExpressionType {
	return RelationalExpressionTypeILike
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionILike) ToMap() map[string]any {
	return map[string]any{
		"type":    j.Type(),
		"expr":    j.Expr,
		"pattern": j.Pattern,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionILike) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNotILike represents a RelationalExpression with the not_i_like type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.ilike`
// * During filtering: `relational_query.filter.comparison.ilike`
// * During sorting:`relational_query.sort.expression.comparison.ilike`
// * During joining: `relational_query.join.expression.comparison.ilike`
// * During aggregation: `relational_query.aggregate.expression.comparison.ilike`
// * During windowing: `relational_query.window.expression.comparison.ilike`.
type RelationalExpressionNotILike struct {
	Expr    RelationalExpression `json:"expr"    mapstructure:"expr"    yaml:"expr"`
	Pattern RelationalExpression `json:"pattern" mapstructure:"pattern" yaml:"pattern"`
}

// NewRelationalExpressionNotILike creates a RelationalExpressionNotILike instance.
func NewRelationalExpressionNotILike[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionNotILike {
	return &RelationalExpressionNotILike{
		Expr:    left.Wrap(),
		Pattern: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNotILike) Type() RelationalExpressionType {
	return RelationalExpressionTypeNotILike
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNotILike) ToMap() map[string]any {
	return map[string]any{
		"type":    j.Type(),
		"expr":    j.Expr,
		"pattern": j.Pattern,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNotILike) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionBetween represents a RelationalExpression with the between type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.between`
// * During filtering: `relational_query.filter.comparison.between`
// * During sorting:`relational_query.sort.expression.comparison.between`
// * During joining: `relational_query.join.expression.comparison.between`
// * During aggregation: `relational_query.aggregate.expression.comparison.between`
// * During windowing: `relational_query.window.expression.comparison.between`.
type RelationalExpressionBetween struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
	Low  RelationalExpression `json:"low"  mapstructure:"low"  yaml:"low"`
	High RelationalExpression `json:"high" mapstructure:"high" yaml:"high"`
}

// NewRelationalExpressionBetween creates a RelationalExpressionBetween instance.
func NewRelationalExpressionBetween[E RelationalExpressionInner, L RelationalExpressionInner, H RelationalExpressionInner](
	expr E,
	low L,
	high H,
) *RelationalExpressionBetween {
	return &RelationalExpressionBetween{
		Expr: expr.Wrap(),
		Low:  low.Wrap(),
		High: high.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionBetween) Type() RelationalExpressionType {
	return RelationalExpressionTypeBetween
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionBetween) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
		"low":  j.Low,
		"high": j.High,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionBetween) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNotBetween represents a RelationalExpression with the not_between type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.between`
// * During filtering: `relational_query.filter.comparison.between`
// * During sorting:`relational_query.sort.expression.comparison.between`
// * During joining: `relational_query.join.expression.comparison.between`
// * During aggregation: `relational_query.aggregate.expression.comparison.between`
// * During windowing: `relational_query.window.expression.comparison.between`.
type RelationalExpressionNotBetween struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
	Low  RelationalExpression `json:"low"  mapstructure:"low"  yaml:"low"`
	High RelationalExpression `json:"high" mapstructure:"high" yaml:"high"`
}

// NewRelationalExpressionNotBetween creates a RelationalExpressionNotBetween instance.
func NewRelationalExpressionNotBetween[E RelationalExpressionInner, L RelationalExpressionInner, H RelationalExpressionInner](
	expr E,
	low L,
	high H,
) *RelationalExpressionNotBetween {
	return &RelationalExpressionNotBetween{
		Expr: expr.Wrap(),
		Low:  low.Wrap(),
		High: high.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNotBetween) Type() RelationalExpressionType {
	return RelationalExpressionTypeNotBetween
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNotBetween) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
		"low":  j.Low,
		"high": j.High,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNotBetween) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionContains represents a RelationalExpression with the contains type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.contains`
// * During filtering: `relational_query.filter.comparison.contains`
// * During sorting:`relational_query.sort.expression.comparison.contains`
// * During joining: `relational_query.join.expression.comparison.contains`
// * During aggregation: `relational_query.aggregate.expression.comparison.contains`
// * During windowing: `relational_query.window.expression.comparison.contains`.
type RelationalExpressionContains struct {
	Str       RelationalExpression `json:"str"        mapstructure:"str"        yaml:"str"`
	SearchStr RelationalExpression `json:"search_str" mapstructure:"search_str" yaml:"search_str"`
}

// NewRelationalExpressionContains creates a RelationalExpressionContains instance.
func NewRelationalExpressionContains[S RelationalExpressionInner, SS RelationalExpressionInner](
	str S,
	searchStr SS,
) *RelationalExpressionContains {
	return &RelationalExpressionContains{
		Str:       str.Wrap(),
		SearchStr: searchStr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionContains) Type() RelationalExpressionType {
	return RelationalExpressionTypeContains
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionContains) ToMap() map[string]any {
	return map[string]any{
		"type":       j.Type(),
		"str":        j.Str,
		"search_str": j.SearchStr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionContains) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsNaN represents a RelationalExpression with the is_na_n type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.is_nan`
// * During filtering: `relational_query.filter.comparison.is_nan`
// * During sorting:`relational_query.sort.expression.comparison.is_nan`
// * During joining: `relational_query.join.expression.comparison.is_nan`
// * During aggregation: `relational_query.aggregate.expression.comparison.is_nan`
// * During windowing: `relational_query.window.expression.comparison.is_nan`.
type RelationalExpressionIsNaN struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsNaN creates a RelationalExpressionIsNaN instance.
func NewRelationalExpressionIsNaN[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsNaN {
	return &RelationalExpressionIsNaN{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsNaN) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsNaN
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsNaN) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsNaN) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionIsZero represents a RelationalExpression with the is_zero type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.comparison.is_zero`
// * During filtering: `relational_query.filter.comparison.is_zero`
// * During sorting:`relational_query.sort.expression.comparison.is_zero`
// * During joining: `relational_query.join.expression.comparison.is_zero`
// * During aggregation: `relational_query.aggregate.expression.comparison.is_zero`
// * During windowing: `relational_query.window.expression.comparison.is_zero.
type RelationalExpressionIsZero struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionIsZero creates a RelationalExpressionIsZero instance.
func NewRelationalExpressionIsZero[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionIsZero {
	return &RelationalExpressionIsZero{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionIsZero) Type() RelationalExpressionType {
	return RelationalExpressionTypeIsZero
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionIsZero) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionIsZero) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionPlus represents a RelationalExpression with the plus type.
type RelationalExpressionPlus struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionPlus creates a RelationalExpressionPlus instance.
func NewRelationalExpressionPlus[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionPlus {
	return &RelationalExpressionPlus{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionPlus) Type() RelationalExpressionType {
	return RelationalExpressionTypePlus
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionPlus) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionPlus) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionMinus represents a RelationalExpression with the minus type.
type RelationalExpressionMinus struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionMinus creates a RelationalExpressionMinus instance.
func NewRelationalExpressionMinus[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionMinus {
	return &RelationalExpressionMinus{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionMinus) Type() RelationalExpressionType {
	return RelationalExpressionTypeMinus
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionMinus) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionMinus) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionMultiply represents a RelationalExpression with the multiply type.
type RelationalExpressionMultiply struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionMultiply creates a RelationalExpressionMultiply instance.
func NewRelationalExpressionMultiply[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionMultiply {
	return &RelationalExpressionMultiply{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionMultiply) Type() RelationalExpressionType {
	return RelationalExpressionTypeMultiply
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionMultiply) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionMultiply) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionDivide represents a RelationalExpression with the divide type.
type RelationalExpressionDivide struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionDivide creates a RelationalExpressionDivide instance.
func NewRelationalExpressionDivide[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionDivide {
	return &RelationalExpressionDivide{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionDivide) Type() RelationalExpressionType {
	return RelationalExpressionTypeDivide
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionDivide) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionDivide) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionModulo represents a RelationalExpression with the modulo type.
type RelationalExpressionModulo struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionModulo creates a RelationalExpressionModulo instance.
func NewRelationalExpressionModulo[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionModulo {
	return &RelationalExpressionModulo{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionModulo) Type() RelationalExpressionType {
	return RelationalExpressionTypeModulo
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionModulo) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionModulo) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNegate represents a RelationalExpression with the negate type.
type RelationalExpressionNegate struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionNegate creates a RelationalExpressionNegate instance.
func NewRelationalExpressionNegate[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionNegate {
	return &RelationalExpressionNegate{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNegate) Type() RelationalExpressionType {
	return RelationalExpressionTypeNegate
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNegate) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNegate) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCast represents a RelationalExpression with the cast type.
type RelationalExpressionCast struct {
	Expr   RelationalExpression `json:"expr"                mapstructure:"expr"                yaml:"expr"`
	AsType CastType             `json:"as_type"             mapstructure:"as_type"             yaml:"as_type"`
	// Optional for now, but will be required in the future
	FromType *CastType `json:"from_type,omitempty" mapstructure:"from_type,omitempty" yaml:"from_type,omitempty"`
}

// NewRelationalExpressionCast creates a RelationalExpressionCast instance.
func NewRelationalExpressionCast[E RelationalExpressionInner, C CastTypeInner](
	expr E,
	asType C,
) *RelationalExpressionCast {
	return &RelationalExpressionCast{
		Expr:   expr.Wrap(),
		AsType: asType.Wrap(),
	}
}

// WithFromType return the instance with from_type set.
func (j RelationalExpressionCast) WithFromType(fromType CastTypeInner) *RelationalExpressionCast {
	if fromType == nil {
		j.FromType = nil
	} else {
		ft := fromType.Wrap()
		j.FromType = &ft
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionCast) Type() RelationalExpressionType {
	return RelationalExpressionTypeCast
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionCast) ToMap() map[string]any {
	result := map[string]any{
		"type":    j.Type(),
		"expr":    j.Expr,
		"as_type": j.AsType,
	}

	if j.FromType != nil && !j.FromType.IsEmpty() {
		result["from_type"] = j.FromType
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCast) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionTryCast represents a RelationalExpression with the try_cast type.
type RelationalExpressionTryCast struct {
	Expr   RelationalExpression `json:"expr"                mapstructure:"expr"                yaml:"expr"`
	AsType CastType             `json:"as_type"             mapstructure:"as_type"             yaml:"as_type"`
	// Optional for now, but will be required in the future
	FromType *CastType `json:"from_type,omitempty" mapstructure:"from_type,omitempty" yaml:"from_type,omitempty"`
}

// NewRelationalExpressionTryCast creates a RelationalExpressionTryCast instance.
func NewRelationalExpressionTryCast[E RelationalExpressionInner, C CastTypeInner](
	expr E,
	asType C,
) *RelationalExpressionTryCast {
	return &RelationalExpressionTryCast{
		Expr:   expr.Wrap(),
		AsType: asType.Wrap(),
	}
}

// WithFromType return the instance with from_type set.
func (j RelationalExpressionTryCast) WithFromType(
	fromType CastTypeInner,
) *RelationalExpressionTryCast {
	if fromType == nil {
		j.FromType = nil
	} else {
		ft := fromType.Wrap()
		j.FromType = &ft
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionTryCast) Type() RelationalExpressionType {
	return RelationalExpressionTypeTryCast
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionTryCast) ToMap() map[string]any {
	result := map[string]any{
		"type":    j.Type(),
		"expr":    j.Expr,
		"as_type": j.AsType,
	}

	if j.FromType != nil && !j.FromType.IsEmpty() {
		result["from_type"] = j.FromType
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionTryCast) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionAbs represents a RelationalExpression with the abs type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.abs`
// * During filtering: `relational_query.filter.scalar.abs`
// * During sorting:`relational_query.sort.expression.scalar.abs`
// * During joining: `relational_query.join.expression.scalar.abs`
// * During aggregation: `relational_query.aggregate.expression.scalar.abs`
// * During windowing: `relational_query.window.expression.scalar.abs`.
type RelationalExpressionAbs struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionAbs creates a RelationalExpressionAbs instance.
func NewRelationalExpressionAbs[E RelationalExpressionInner](expr E) *RelationalExpressionAbs {
	return &RelationalExpressionAbs{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionAbs) Type() RelationalExpressionType {
	return RelationalExpressionTypeAbs
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionAbs) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionAbs) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionArrayElement represents a RelationalExpression with the array_element type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.array_element`
// * During filtering: `relational_query.filter.scalar.array_element`
// * During sorting:`relational_query.sort.expression.scalar.array_element`
// * During joining: `relational_query.join.expression.scalar.array_element`
// * During aggregation: `relational_query.aggregate.expression.scalar.array_element`
// * During windowing: `relational_query.window.expression.scalar.array_element`.
type RelationalExpressionArrayElement struct {
	Column RelationalExpression `json:"column" mapstructure:"column" yaml:"column"`
	Index  uint                 `json:"index"  mapstructure:"index"  yaml:"index"`
}

// NewRelationalExpressionArrayElement creates a RelationalExpressionArrayElement instance.
func NewRelationalExpressionArrayElement[E RelationalExpressionInner](
	column E,
	index uint,
) *RelationalExpressionArrayElement {
	return &RelationalExpressionArrayElement{
		Column: column.Wrap(),
		Index:  index,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionArrayElement) Type() RelationalExpressionType {
	return RelationalExpressionTypeArrayElement
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionArrayElement) ToMap() map[string]any {
	return map[string]any{
		"type":   j.Type(),
		"column": j.Column,
		"index":  j.Index,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionArrayElement) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionBTrim represents a RelationalExpression with the b_trim type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.btrim`
// * During filtering: `relational_query.filter.scalar.btrim`
// * During sorting:`relational_query.sort.expression.scalar.btrim`
// * During joining: `relational_query.join.expression.scalar.btrim`
// * During aggregation: `relational_query.aggregate.expression.scalar.btrim`
// * During windowing: `relational_query.window.expression.scalar.btrim`.
type RelationalExpressionBTrim struct {
	Str     RelationalExpression  `json:"str"                mapstructure:"str"                yaml:"str"`
	TrimStr *RelationalExpression `json:"trim_str,omitempty" mapstructure:"trim_str,omitempty" yaml:"trim_str,omitempty"`
}

// NewRelationalExpressionBTrim creates a RelationalExpressionBTrim instance.
func NewRelationalExpressionBTrim[E RelationalExpressionInner](str E) *RelationalExpressionBTrim {
	return &RelationalExpressionBTrim{
		Str: str.Wrap(),
	}
}

// WithTrimStr return the instance with new trim_str.
func (j RelationalExpressionBTrim) WithTrimStr(
	trimStr RelationalExpressionInner,
) *RelationalExpressionBTrim {
	if trimStr == nil {
		j.TrimStr = nil
	} else {
		ts := trimStr.Wrap()
		j.TrimStr = &ts
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionBTrim) Type() RelationalExpressionType {
	return RelationalExpressionTypeBTrim
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionBTrim) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"str":  j.Str,
	}

	if j.TrimStr != nil && !j.TrimStr.IsEmpty() {
		result["trim_str"] = j.TrimStr
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionBTrim) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCeil represents a RelationalExpression with the ceil type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.ceil`
// * During filtering: `relational_query.filter.scalar.ceil`
// * During sorting:`relational_query.sort.expression.scalar.ceil`
// * During joining: `relational_query.join.expression.scalar.ceil`
// * During aggregation: `relational_query.aggregate.expression.scalar.ceil`
// * During windowing: `relational_query.window.expression.scalar.ceil`.
type RelationalExpressionCeil struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionCeil creates a RelationalExpressionCeil instance.
func NewRelationalExpressionCeil[E RelationalExpressionInner](expr E) *RelationalExpressionCeil {
	return &RelationalExpressionCeil{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionCeil) Type() RelationalExpressionType {
	return RelationalExpressionTypeCeil
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionCeil) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCeil) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCharacterLength represents a RelationalExpression with the character_length type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.character_length`
// * During filtering: `relational_query.filter.scalar.character_length`
// * During sorting:`relational_query.sort.expression.scalar.character_length`
// * During joining: `relational_query.join.expression.scalar.character_length`
// * During aggregation: `relational_query.aggregate.expression.scalar.character_length`
// * During windowing: `relational_query.window.expression.scalar.character_length`.
type RelationalExpressionCharacterLength struct {
	Str RelationalExpression `json:"str" mapstructure:"str" yaml:"str"`
}

// NewRelationalExpressionCharacterLength creates a RelationalExpressionCharacterLength instance.
func NewRelationalExpressionCharacterLength[E RelationalExpressionInner](
	str E,
) *RelationalExpressionCharacterLength {
	return &RelationalExpressionCharacterLength{
		Str: str.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionCharacterLength) Type() RelationalExpressionType {
	return RelationalExpressionTypeCharacterLength
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionCharacterLength) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"str":  j.Str,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCharacterLength) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCoalesce represents a RelationalExpression with the coalesce type.
type RelationalExpressionCoalesce struct {
	Exprs []RelationalExpression `json:"exprs" mapstructure:"exprs" yaml:"exprs"`
}

// NewRelationalExpressionCoalesce creates a RelationalExpressionCoalesce instance.
func NewRelationalExpressionCoalesce(
	expressions []RelationalExpressionInner,
) *RelationalExpressionCoalesce {
	exprs := []RelationalExpression{}

	for _, expr := range expressions {
		if expr != nil {
			exprs = append(exprs, expr.Wrap())
		}
	}

	return &RelationalExpressionCoalesce{
		Exprs: exprs,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionCoalesce) Type() RelationalExpressionType {
	return RelationalExpressionTypeCoalesce
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionCoalesce) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"exprs": j.Exprs,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCoalesce) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionConcat represents a RelationalExpression with the concat type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.concat`
// * During filtering: `relational_query.filter.scalar.concat`
// * During sorting:`relational_query.sort.expression.scalar.concat`
// * During joining: `relational_query.join.expression.scalar.concat`
// * During aggregation: `relational_query.aggregate.expression.scalar.concat`
// * During windowing: `relational_query.window.expression.scalar.concat`.
type RelationalExpressionConcat struct {
	Exprs []RelationalExpression `json:"exprs" mapstructure:"exprs" yaml:"exprs"`
}

// NewRelationalExpressionConcat creates a RelationalExpressionConcat instance.
func NewRelationalExpressionConcat(
	expressions []RelationalExpressionInner,
) *RelationalExpressionConcat {
	exprs := []RelationalExpression{}

	for _, expr := range expressions {
		if expr != nil {
			exprs = append(exprs, expr.Wrap())
		}
	}

	return &RelationalExpressionConcat{
		Exprs: exprs,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionConcat) Type() RelationalExpressionType {
	return RelationalExpressionTypeConcat
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionConcat) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"exprs": j.Exprs,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionConcat) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCos represents a RelationalExpression with the cos type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.cos`
// * During filtering: `relational_query.filter.scalar.cos`
// * During sorting:`relational_query.sort.expression.scalar.cos`
// * During joining: `relational_query.join.expression.scalar.cos`
// * During aggregation: `relational_query.aggregate.expression.scalar.cos`
// * During windowing: `relational_query.window.expression.scalar.cos`.
type RelationalExpressionCos struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionCos creates a RelationalExpressionCos instance.
func NewRelationalExpressionCos[E RelationalExpressionInner](expr E) *RelationalExpressionCos {
	return &RelationalExpressionCos{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionCos) Type() RelationalExpressionType {
	return RelationalExpressionTypeCos
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionCos) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCos) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCurrentDate represents a RelationalExpression with the current_date type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.current_date`
// * During filtering: `relational_query.filter.scalar.current_date`
// * During sorting:`relational_query.sort.expression.scalar.current_date`
// * During joining: `relational_query.join.expression.scalar.current_date`
// * During aggregation: `relational_query.aggregate.expression.scalar.current_date`
// * During windowing: `relational_query.window.expression.scalar.current_date`.
type RelationalExpressionCurrentDate struct{}

// NewRelationalExpressionCurrentDate creates a RelationalExpressionCurrentDate instance.
func NewRelationalExpressionCurrentDate() *RelationalExpressionCurrentDate {
	return &RelationalExpressionCurrentDate{}
}

// Type return the type name of the instance.
func (j RelationalExpressionCurrentDate) Type() RelationalExpressionType {
	return RelationalExpressionTypeCurrentDate
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionCurrentDate) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCurrentDate) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCurrentTime represents a RelationalExpression with the current_time type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.current_time`
// * During filtering: `relational_query.filter.scalar.current_time`
// * During sorting:`relational_query.sort.expression.scalar.current_time`
// * During joining: `relational_query.join.expression.scalar.current_time`
// * During aggregation: `relational_query.aggregate.expression.scalar.current_time`
// * During windowing: `relational_query.window.expression.scalar.current_time`.
type RelationalExpressionCurrentTime struct{}

// NewRelationalExpressionCurrentTime creates a RelationalExpressionCurrentTime instance.
func NewRelationalExpressionCurrentTime() *RelationalExpressionCurrentTime {
	return &RelationalExpressionCurrentTime{}
}

// Type return the type name of the instance.
func (j RelationalExpressionCurrentTime) Type() RelationalExpressionType {
	return RelationalExpressionTypeCurrentTime
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionCurrentTime) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCurrentTime) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCurrentTimestamp represents a RelationalExpression with the current_timestamp type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.current_timestamp`
// * During filtering: `relational_query.filter.scalar.current_timestamp`
// * During sorting:`relational_query.sort.expression.scalar.current_timestamp`
// * During joining: `relational_query.join.expression.scalar.current_timestamp`
// * During aggregation: `relational_query.aggregate.expression.scalar.current_timestamp`
// * During windowing: `relational_query.window.expression.scalar.current_timestamp`.
type RelationalExpressionCurrentTimestamp struct{}

// NewRelationalExpressionCurrentTimestamp creates a RelationalExpressionCurrentTime instance.
func NewRelationalExpressionCurrentTimestamp() *RelationalExpressionCurrentTimestamp {
	return &RelationalExpressionCurrentTimestamp{}
}

// Type return the type name of the instance.
func (j RelationalExpressionCurrentTimestamp) Type() RelationalExpressionType {
	return RelationalExpressionTypeCurrentTimestamp
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionCurrentTimestamp) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCurrentTimestamp) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionDatePart represents a RelationalExpression with the date_part type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.date_part`
// * During filtering: `relational_query.filter.scalar.date_part`
// * During sorting:`relational_query.sort.expression.scalar.date_part`
// * During joining: `relational_query.join.expression.scalar.date_part`
// * During aggregation: `relational_query.aggregate.expression.scalar.date_part`
// * During windowing: `relational_query.window.expression.scalar.date_part`.
type RelationalExpressionDatePart struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
	Part DatePartUnit         `json:"part" mapstructure:"part" yaml:"part"`
}

// NewRelationalExpressionDatePart creates a RelationalExpressionDatePart instance.
func NewRelationalExpressionDatePart[E RelationalExpressionInner](
	expr E,
	part DatePartUnit,
) *RelationalExpressionDatePart {
	return &RelationalExpressionDatePart{
		Expr: expr.Wrap(),
		Part: part,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionDatePart) Type() RelationalExpressionType {
	return RelationalExpressionTypeDatePart
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionDatePart) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
		"part": j.Part,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionDatePart) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionDateTrunc represents a RelationalExpression with the date_trunc type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.date_trunc`
// * During filtering: `relational_query.filter.scalar.date_trunc`
// * During sorting:`relational_query.sort.expression.scalar.date_trunc`
// * During joining: `relational_query.join.expression.scalar.date_trunc`
// * During aggregation: `relational_query.aggregate.expression.scalar.date_trunc`
// * During windowing: `relational_query.window.expression.scalar.date_trunc`.
type RelationalExpressionDateTrunc struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
	Part RelationalExpression `json:"part" mapstructure:"part" yaml:"part"`
}

// NewRelationalExpressionDateTrunc creates a RelationalExpressionDateTrunc instance.
func NewRelationalExpressionDateTrunc[E RelationalExpressionInner, R RelationalExpressionInner](
	expr E,
	part R,
) *RelationalExpressionDateTrunc {
	return &RelationalExpressionDateTrunc{
		Expr: expr.Wrap(),
		Part: part.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionDateTrunc) Type() RelationalExpressionType {
	return RelationalExpressionTypeDateTrunc
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionDateTrunc) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
		"part": j.Part,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionDateTrunc) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionExp represents a RelationalExpression with the exp type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.exp`
// * During filtering: `relational_query.filter.scalar.exp`
// * During sorting:`relational_query.sort.expression.scalar.exp`
// * During joining: `relational_query.join.expression.scalar.exp`
// * During aggregation: `relational_query.aggregate.expression.scalar.exp`
// * During windowing: `relational_query.window.expression.scalar.exp`.
type RelationalExpressionExp struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionExp creates a RelationalExpressionExp instance.
func NewRelationalExpressionExp[E RelationalExpressionInner](expr E) *RelationalExpressionExp {
	return &RelationalExpressionExp{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionExp) Type() RelationalExpressionType {
	return RelationalExpressionTypeExp
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionExp) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionExp) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionFloor represents a RelationalExpression with the floor type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.floor`
// * During filtering: `relational_query.filter.scalar.floor`
// * During sorting:`relational_query.sort.expression.scalar.floor`
// * During joining: `relational_query.join.expression.scalar.floor`
// * During aggregation: `relational_query.aggregate.expression.scalar.floor`
// * During windowing: `relational_query.window.expression.scalar.floor`.
type RelationalExpressionFloor struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionFloor creates a RelationalExpressionFloor instance.
func NewRelationalExpressionFloor[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionFloor {
	return &RelationalExpressionFloor{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionFloor) Type() RelationalExpressionType {
	return RelationalExpressionTypeFloor
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionFloor) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionFloor) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionGetField represents a RelationalExpression with the get_field type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.get_field`
// * During filtering: `relational_query.filter.scalar.get_field`
// * During sorting:`relational_query.sort.expression.scalar.get_field`
// * During joining: `relational_query.join.expression.scalar.get_field`
// * During aggregation: `relational_query.aggregate.expression.scalar.get_field`
// * During windowing: `relational_query.window.expression.scalar.get_field`.
type RelationalExpressionGetField struct {
	Column RelationalExpression `json:"column" mapstructure:"column" yaml:"column"`
	Field  string               `json:"field"  mapstructure:"field"  yaml:"field"`
}

// NewRelationalExpressionGetField creates a RelationalExpressionGetField instance.
func NewRelationalExpressionGetField[E RelationalExpressionInner](
	column E,
	field string,
) *RelationalExpressionGetField {
	return &RelationalExpressionGetField{
		Column: column.Wrap(),
		Field:  field,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionGetField) Type() RelationalExpressionType {
	return RelationalExpressionTypeGetField
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionGetField) ToMap() map[string]any {
	return map[string]any{
		"type":   j.Type(),
		"column": j.Column,
		"field":  j.Field,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionGetField) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionGreatest represents a RelationalExpression with the greatest type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.greatest`
// * During filtering: `relational_query.filter.scalar.greatest`
// * During sorting:`relational_query.sort.expression.scalar.greatest`
// * During joining: `relational_query.join.expression.scalar.greatest`
// * During aggregation: `relational_query.aggregate.expression.scalar.greatest`
// * During windowing: `relational_query.window.expression.scalar.greatest`.
type RelationalExpressionGreatest struct {
	Exprs []RelationalExpression `json:"exprs" mapstructure:"exprs" yaml:"exprs"`
}

// NewRelationalExpressionGreatest creates a RelationalExpressionGreatest instance.
func NewRelationalExpressionGreatest(
	expressions []RelationalExpressionInner,
) *RelationalExpressionGreatest {
	exprs := []RelationalExpression{}

	for _, expr := range expressions {
		if expr != nil {
			exprs = append(exprs, expr.Wrap())
		}
	}

	return &RelationalExpressionGreatest{
		Exprs: exprs,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionGreatest) Type() RelationalExpressionType {
	return RelationalExpressionTypeGreatest
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionGreatest) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"exprs": j.Exprs,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionGreatest) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLeast represents a RelationalExpression with the least type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.least`
// * During filtering: `relational_query.filter.scalar.least`
// * During sorting:`relational_query.sort.expression.scalar.least`
// * During joining: `relational_query.join.expression.scalar.least`
// * During aggregation: `relational_query.aggregate.expression.scalar.least`
// * During windowing: `relational_query.window.expression.scalar.least`.
type RelationalExpressionLeast struct {
	Exprs []RelationalExpression `json:"least" mapstructure:"least" yaml:"least"`
}

// NewRelationalExpressionLeast creates a RelationalExpressionLeast instance.
func NewRelationalExpressionLeast(
	expressions []RelationalExpressionInner,
) *RelationalExpressionLeast {
	exprs := []RelationalExpression{}

	for _, expr := range expressions {
		if expr != nil {
			exprs = append(exprs, expr.Wrap())
		}
	}

	return &RelationalExpressionLeast{
		Exprs: exprs,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLeast) Type() RelationalExpressionType {
	return RelationalExpressionTypeLeast
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionLeast) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"exprs": j.Exprs,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLeast) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLeft represents a RelationalExpression with the left type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.left`
// * During filtering: `relational_query.filter.scalar.left`
// * During sorting:`relational_query.sort.expression.scalar.left`
// * During joining: `relational_query.join.expression.scalar.left`
// * During aggregation: `relational_query.aggregate.expression.scalar.left`
// * During windowing: `relational_query.window.expression.scalar.left`.
type RelationalExpressionLeft struct {
	Str RelationalExpression `json:"str" mapstructure:"str" yaml:"str"`
	N   RelationalExpression `json:"n"   mapstructure:"n"   yaml:"n"`
}

// NewRelationalExpressionLeft creates a RelationalExpressionLeft instance.
func NewRelationalExpressionLeft[E RelationalExpressionInner, N RelationalExpressionInner](
	str E,
	n N,
) *RelationalExpressionLeft {
	return &RelationalExpressionLeft{
		Str: str.Wrap(),
		N:   n.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLeft) Type() RelationalExpressionType {
	return RelationalExpressionTypeLeft
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLeft) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"str":  j.Str,
		"n":    j.N,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLeft) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLn represents a RelationalExpression with the ln type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.ln`
// * During filtering: `relational_query.filter.scalar.ln`
// * During sorting:`relational_query.sort.expression.scalar.ln`
// * During joining: `relational_query.join.expression.scalar.ln`
// * During aggregation: `relational_query.aggregate.expression.scalar.ln`
// * During windowing: `relational_query.window.expression.scalar.ln`.
type RelationalExpressionLn struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionLn creates a RelationalExpressionLn instance.
func NewRelationalExpressionLn[E RelationalExpressionInner](expr E) *RelationalExpressionLn {
	return &RelationalExpressionLn{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLn) Type() RelationalExpressionType {
	return RelationalExpressionTypeLn
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLn) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLn) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLog represents a RelationalExpression with the log type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.log`
// * During filtering: `relational_query.filter.scalar.log`
// * During sorting:`relational_query.sort.expression.scalar.log`
// * During joining: `relational_query.join.expression.scalar.log`
// * During aggregation: `relational_query.aggregate.expression.scalar.log`
// * During windowing: `relational_query.window.expression.scalar.log`.
type RelationalExpressionLog struct {
	Expr RelationalExpression  `json:"expr"           mapstructure:"expr"           yaml:"expr"`
	Base *RelationalExpression `json:"base,omitempty" mapstructure:"base,omitempty" yaml:"base,omitempty"`
}

// NewRelationalExpressionLog creates a RelationalExpressionLog instance.
func NewRelationalExpressionLog[E RelationalExpressionInner](expr E) *RelationalExpressionLog {
	return &RelationalExpressionLog{
		Expr: expr.Wrap(),
	}
}

// WithBase return the instance with a new base set.
func (j RelationalExpressionLog) WithBase(
	base RelationalExpressionInner,
) *RelationalExpressionLog {
	if base == nil {
		j.Base = nil
	} else {
		b := base.Wrap()
		j.Base = &b
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionLog) Type() RelationalExpressionType {
	return RelationalExpressionTypeLog
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLog) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}

	if j.Base != nil && !j.Base.IsEmpty() {
		result["base"] = j.Base
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLog) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLog10 represents a RelationalExpression with the log10 type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.log10`
// * During filtering: `relational_query.filter.scalar.log10`
// * During sorting:`relational_query.sort.expression.scalar.log10`
// * During joining: `relational_query.join.expression.scalar.log10`
// * During aggregation: `relational_query.aggregate.expression.scalar.log10`
// * During windowing: `relational_query.window.expression.scalar.log10`.
type RelationalExpressionLog10 struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionLog10 creates a RelationalExpressionLog10 instance.
func NewRelationalExpressionLog10[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionLog10 {
	return &RelationalExpressionLog10{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLog10) Type() RelationalExpressionType {
	return RelationalExpressionTypeLog10
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLog10) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLog10) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLog2 represents a RelationalExpression with the log2 type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.log2`
// * During filtering: `relational_query.filter.scalar.log2`
// * During sorting:`relational_query.sort.expression.scalar.log2`
// * During joining: `relational_query.join.expression.scalar.log2`
// * During aggregation: `relational_query.aggregate.expression.scalar.log2`
// * During windowing: `relational_query.window.expression.scalar.log2`.
type RelationalExpressionLog2 struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionLog2 creates a RelationalExpressionLog2 instance.
func NewRelationalExpressionLog2[E RelationalExpressionInner](expr E) *RelationalExpressionLog2 {
	return &RelationalExpressionLog2{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLog2) Type() RelationalExpressionType {
	return RelationalExpressionTypeLog2
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLog2) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLog2) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLPad represents a RelationalExpression with the l_pad type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.lpad`
// * During filtering: `relational_query.filter.scalar.lpad`
// * During sorting:`relational_query.sort.expression.scalar.lpad`
// * During joining: `relational_query.join.expression.scalar.lpad`
// * During aggregation: `relational_query.aggregate.expression.scalar.lpad`
// * During windowing: `relational_query.window.expression.scalar.lpad`.
type RelationalExpressionLPad struct {
	Str        RelationalExpression  `json:"str"                   mapstructure:"str"                   yaml:"str"`
	N          RelationalExpression  `json:"n"                     mapstructure:"n"                     yaml:"n"`
	PaddingStr *RelationalExpression `json:"padding_str,omitempty" mapstructure:"padding_str,omitempty" yaml:"padding_str,omitempty"`
}

// NewRelationalExpressionLPad creates a RelationalExpressionLPad instance.
func NewRelationalExpressionLPad[E RelationalExpressionInner, N RelationalExpressionInner](
	str E,
	n N,
) *RelationalExpressionLPad {
	return &RelationalExpressionLPad{
		Str: str.Wrap(),
		N:   n.Wrap(),
	}
}

// WithPaddingStr returns the new instance with padding str set.
func (j RelationalExpressionLPad) WithPaddingStr(
	paddingStr RelationalExpressionInner,
) *RelationalExpressionLPad {
	if paddingStr == nil {
		j.PaddingStr = nil
	} else {
		ps := paddingStr.Wrap()
		j.PaddingStr = &ps
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionLPad) Type() RelationalExpressionType {
	return RelationalExpressionTypeLPad
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLPad) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"str":  j.Str,
		"n":    j.N,
	}

	if j.PaddingStr != nil && !j.PaddingStr.IsEmpty() {
		result["padding_str"] = j.PaddingStr
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLPad) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLTrim represents a RelationalExpression with the l_trim type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.ltrim`
// * During filtering: `relational_query.filter.scalar.ltrim`
// * During sorting:`relational_query.sort.expression.scalar.ltrim`
// * During joining: `relational_query.join.expression.scalar.ltrim`
// * During aggregation: `relational_query.aggregate.expression.scalar.ltrim`
// * During windowing: `relational_query.window.expression.scalar.ltrim`.
type RelationalExpressionLTrim struct {
	Str     RelationalExpression  `json:"str"                mapstructure:"str"                yaml:"str"`
	TrimStr *RelationalExpression `json:"trim_str,omitempty" mapstructure:"trim_str,omitempty" yaml:"trim_str,omitempty"`
}

// NewRelationalExpressionLTrim creates a RelationalExpressionLTrim instance.
func NewRelationalExpressionLTrim[E RelationalExpressionInner](str E) *RelationalExpressionLTrim {
	return &RelationalExpressionLTrim{
		Str: str.Wrap(),
	}
}

// WithTrimStr returns the new instance with trim str set.
func (j RelationalExpressionLTrim) WithTrimStr(
	trimStr RelationalExpressionInner,
) *RelationalExpressionLTrim {
	if trimStr == nil {
		j.TrimStr = nil
	} else {
		ts := trimStr.Wrap()
		j.TrimStr = &ts
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionLTrim) Type() RelationalExpressionType {
	return RelationalExpressionTypeLTrim
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLTrim) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"str":  j.Str,
	}

	if j.TrimStr != nil && !j.TrimStr.IsEmpty() {
		result["trim_str"] = j.TrimStr
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLTrim) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNullIf represents a relational null_if expression type.
type RelationalExpressionNullIf struct {
	Expr1 RelationalExpression `json:"expr1" mapstructure:"expr1" yaml:"expr1"`
	Expr2 RelationalExpression `json:"expr2" mapstructure:"expr2" yaml:"expr2"`
}

// NewRelationalExpressionNullIf creates a RelationalExpressionNullIf instance.
func NewRelationalExpressionNullIf[L RelationalExpressionInner, R RelationalExpressionInner](
	expr1 L,
	expr2 R,
) *RelationalExpressionNullIf {
	return &RelationalExpressionNullIf{
		Expr1: expr1.Wrap(),
		Expr2: expr2.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNullIf) Type() RelationalExpressionType {
	return RelationalExpressionTypeNullIf
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNullIf) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"expr1": j.Expr1,
		"expr2": j.Expr2,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNullIf) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNvl represents a relational nvl expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.nvl`
// * During filtering: `relational_query.filter.scalar.nvl`
// * During sorting:`relational_query.sort.expression.scalar.nvl`
// * During joining: `relational_query.join.expression.scalar.nvl`
// * During aggregation: `relational_query.aggregate.expression.scalar.nvl`
// * During windowing: `relational_query.window.expression.scalar.nvl`.
type RelationalExpressionNvl struct {
	Expr1 RelationalExpression `json:"expr1" mapstructure:"expr1" yaml:"expr1"`
	Expr2 RelationalExpression `json:"expr2" mapstructure:"expr2" yaml:"expr2"`
}

// NewRelationalExpressionNvl creates a RelationalExpressionNvl instance.
func NewRelationalExpressionNvl[L RelationalExpressionInner, R RelationalExpressionInner](
	expr1 L,
	expr2 R,
) *RelationalExpressionNvl {
	return &RelationalExpressionNvl{
		Expr1: expr1.Wrap(),
		Expr2: expr2.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNvl) Type() RelationalExpressionType {
	return RelationalExpressionTypeNvl
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNvl) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"expr1": j.Expr1,
		"expr2": j.Expr2,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNvl) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionPower represents a relational power expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.power`
// * During filtering: `relational_query.filter.scalar.power`
// * During sorting:`relational_query.sort.expression.scalar.power`
// * During joining: `relational_query.join.expression.scalar.power`
// * During aggregation: `relational_query.aggregate.expression.scalar.power`
// * During windowing: `relational_query.window.expression.scalar.power`.
type RelationalExpressionPower struct {
	Base RelationalExpression `json:"base" mapstructure:"base" yaml:"base"`
	Exp  RelationalExpression `json:"exp"  mapstructure:"exp"  yaml:"exp"`
}

// NewRelationalExpressionPower creates a RelationalExpressionPower instance.
func NewRelationalExpressionPower[L RelationalExpressionInner, R RelationalExpressionInner](
	base L,
	exp R,
) *RelationalExpressionPower {
	return &RelationalExpressionPower{
		Base: base.Wrap(),
		Exp:  exp.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionPower) Type() RelationalExpressionType {
	return RelationalExpressionTypePower
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionPower) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"base": j.Base,
		"exp":  j.Exp,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionPower) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionRandom represents a relational current_date expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.random`
// * During filtering: `relational_query.filter.scalar.random`
// * During sorting:`relational_query.sort.expression.scalar.random`
// * During joining: `relational_query.join.expression.scalar.random`
// * During aggregation: `relational_query.aggregate.expression.scalar.random`
// * During windowing: `relational_query.window.expression.scalar.random`.
type RelationalExpressionRandom struct{}

// NewRelationalExpressionRandom creates a RelationalExpressionRandom instance.
func NewRelationalExpressionRandom() *RelationalExpressionRandom {
	return &RelationalExpressionRandom{}
}

// Type return the type name of the instance.
func (j RelationalExpressionRandom) Type() RelationalExpressionType {
	return RelationalExpressionTypeRandom
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionRandom) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionRandom) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionReplace represents a relational replace expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.replace`
// * During filtering: `relational_query.filter.scalar.replace`
// * During sorting:`relational_query.sort.expression.scalar.replace`
// * During joining: `relational_query.join.expression.scalar.replace`
// * During aggregation: `relational_query.aggregate.expression.scalar.replace`
// * During windowing: `relational_query.window.expression.scalar.replace`.
type RelationalExpressionReplace struct {
	Str         RelationalExpression `json:"str"         mapstructure:"str"         yaml:"str"`
	Substr      RelationalExpression `json:"substr"      mapstructure:"substr"      yaml:"substr"`
	Replacement RelationalExpression `json:"replacement" mapstructure:"replacement" yaml:"replacement"`
}

// NewRelationalExpressionReplace creates a RelationalExpressionReplace instance.
func NewRelationalExpressionReplace[S RelationalExpressionInner, SS RelationalExpressionInner, R RelationalExpressionInner](
	str S,
	substr SS,
	replacement R,
) *RelationalExpressionReplace {
	return &RelationalExpressionReplace{
		Str:         str.Wrap(),
		Substr:      substr.Wrap(),
		Replacement: replacement.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionReplace) Type() RelationalExpressionType {
	return RelationalExpressionTypeReplace
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionReplace) ToMap() map[string]any {
	return map[string]any{
		"type":        j.Type(),
		"str":         j.Str,
		"substr":      j.Substr,
		"replacement": j.Replacement,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionReplace) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionReverse represents a relational reverse expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.reverse`
// * During filtering: `relational_query.filter.scalar.reverse`
// * During sorting:`relational_query.sort.expression.scalar.reverse`
// * During joining: `relational_query.join.expression.scalar.reverse`
// * During aggregation: `relational_query.aggregate.expression.scalar.reverse`
// * During windowing: `relational_query.window.expression.scalar.reverse`.
type RelationalExpressionReverse struct {
	Str RelationalExpression `json:"str" mapstructure:"str" yaml:"str"`
}

// NewRelationalExpressionReverse creates a RelationalExpressionReverse instance.
func NewRelationalExpressionReverse[S RelationalExpressionInner](
	str S,
) *RelationalExpressionReverse {
	return &RelationalExpressionReverse{
		Str: str.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionReverse) Type() RelationalExpressionType {
	return RelationalExpressionTypeReverse
}

// ToMap converts the instance to raw map.
func (j RelationalExpressionReverse) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"str":  j.Str,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionReverse) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionRight represents a RelationalExpression with the right type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.right`
// * During filtering: `relational_query.filter.scalar.right`
// * During sorting:`relational_query.sort.expression.scalar.right`
// * During joining: `relational_query.join.expression.scalar.right`
// * During aggregation: `relational_query.aggregate.expression.scalar.right`
// * During windowing: `relational_query.window.expression.scalar.right`.
type RelationalExpressionRight struct {
	Str RelationalExpression `json:"str" mapstructure:"str" yaml:"str"`
	N   RelationalExpression `json:"n"   mapstructure:"n"   yaml:"n"`
}

// NewRelationalExpressionRight creates a RelationalExpressionRight instance.
func NewRelationalExpressionRight[E RelationalExpressionInner, N RelationalExpressionInner](
	str E,
	n N,
) *RelationalExpressionRight {
	return &RelationalExpressionRight{
		Str: str.Wrap(),
		N:   n.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionRight) Type() RelationalExpressionType {
	return RelationalExpressionTypeRight
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionRight) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"str":  j.Str,
		"n":    j.N,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionRight) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionRound represents a RelationalExpression with the round type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.round`
// * During filtering: `relational_query.filter.scalar.round`
// * During sorting:`relational_query.sort.expression.scalar.round`
// * During joining: `relational_query.join.expression.scalar.round`
// * During aggregation: `relational_query.aggregate.expression.scalar.round`
// * During windowing: `relational_query.window.expression.scalar.round`.
type RelationalExpressionRound struct {
	Expr RelationalExpression  `json:"expr"           mapstructure:"expr"           yaml:"expr"`
	Prec *RelationalExpression `json:"prec,omitempty" mapstructure:"prec,omitempty" yaml:"prec,omitempty"`
}

// NewRelationalExpressionRound creates a RelationalExpressionRound instance.
func NewRelationalExpressionRound[E RelationalExpressionInner](str E) *RelationalExpressionRound {
	return &RelationalExpressionRound{
		Expr: str.Wrap(),
	}
}

// WithPrec returns the new instance with prec set.
func (j RelationalExpressionRound) WithPrec(
	prec RelationalExpressionInner,
) *RelationalExpressionRound {
	if prec == nil {
		j.Prec = nil
	} else {
		ts := prec.Wrap()
		j.Prec = &ts
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionRound) Type() RelationalExpressionType {
	return RelationalExpressionTypeRound
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionRound) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}

	if j.Prec != nil && !j.Prec.IsEmpty() {
		result["prec"] = j.Prec
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionRound) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionRPad represents a RelationalExpression with the r_pad type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.rpad`
// * During filtering: `relational_query.filter.scalar.rpad`
// * During sorting:`relational_query.sort.expression.scalar.rpad`
// * During joining: `relational_query.join.expression.scalar.rpad`
// * During aggregation: `relational_query.aggregate.expression.scalar.rpad`
// * During windowing: `relational_query.window.expression.scalar.rpad`.
type RelationalExpressionRPad struct {
	Str        RelationalExpression  `json:"str"                   mapstructure:"str"                   yaml:"str"`
	N          RelationalExpression  `json:"n"                     mapstructure:"n"                     yaml:"n"`
	PaddingStr *RelationalExpression `json:"padding_str,omitempty" mapstructure:"padding_str,omitempty" yaml:"padding_str,omitempty"`
}

// NewRelationalExpressionRPad creates a RelationalExpressionRPad instance.
func NewRelationalExpressionRPad[E RelationalExpressionInner, N RelationalExpressionInner](
	str E,
	n N,
) *RelationalExpressionRPad {
	return &RelationalExpressionRPad{
		Str: str.Wrap(),
		N:   n.Wrap(),
	}
}

// WithPaddingStr returns the new instance with padding str set.
func (j RelationalExpressionRPad) WithPaddingStr(
	paddingStr RelationalExpressionInner,
) *RelationalExpressionRPad {
	if paddingStr == nil {
		j.PaddingStr = nil
	} else {
		ps := paddingStr.Wrap()
		j.PaddingStr = &ps
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionRPad) Type() RelationalExpressionType {
	return RelationalExpressionTypeLPad
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionRPad) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"str":  j.Str,
		"n":    j.N,
	}

	if j.PaddingStr != nil && !j.PaddingStr.IsEmpty() {
		result["padding_str"] = j.PaddingStr
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionRPad) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionRTrim represents a RelationalExpression with the r_trim type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.rtrim`
// * During filtering: `relational_query.filter.scalar.rtrim`
// * During sorting:`relational_query.sort.expression.scalar.rtrim`
// * During joining: `relational_query.join.expression.scalar.rtrim`
// * During aggregation: `relational_query.aggregate.expression.scalar.rtrim`
// * During windowing: `relational_query.window.expression.scalar.rtrim`.
type RelationalExpressionRTrim struct {
	Str     RelationalExpression  `json:"str"                mapstructure:"str"                yaml:"str"`
	TrimStr *RelationalExpression `json:"trim_str,omitempty" mapstructure:"trim_str,omitempty" yaml:"trim_str,omitempty"`
}

// NewRelationalExpressionRTrim creates a RelationalExpressionRTrim instance.
func NewRelationalExpressionRTrim[E RelationalExpressionInner](str E) *RelationalExpressionRTrim {
	return &RelationalExpressionRTrim{
		Str: str.Wrap(),
	}
}

// WithTrimStr returns the new instance with trim str set.
func (j RelationalExpressionRTrim) WithTrimStr(
	trimStr RelationalExpressionInner,
) *RelationalExpressionRTrim {
	if trimStr == nil {
		j.TrimStr = nil
	} else {
		ts := trimStr.Wrap()
		j.TrimStr = &ts
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionRTrim) Type() RelationalExpressionType {
	return RelationalExpressionTypeRTrim
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionRTrim) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"str":  j.Str,
	}

	if j.TrimStr != nil && !j.TrimStr.IsEmpty() {
		result["trim_str"] = j.TrimStr
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionRTrim) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionSqrt represents a RelationalExpression with the sqrt type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.sqrt`
// * During filtering: `relational_query.filter.scalar.sqrt`
// * During sorting:`relational_query.sort.expression.scalar.sqrt`
// * During joining: `relational_query.join.expression.scalar.sqrt`
// * During aggregation: `relational_query.aggregate.expression.scalar.sqrt`
// * During windowing: `relational_query.window.expression.scalar.sqrt`.
type RelationalExpressionSqrt struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionSqrt creates a RelationalExpressionSqrt instance.
func NewRelationalExpressionSqrt[E RelationalExpressionInner](expr E) *RelationalExpressionSqrt {
	return &RelationalExpressionSqrt{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionSqrt) Type() RelationalExpressionType {
	return RelationalExpressionTypeSqrt
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionSqrt) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionSqrt) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionStrPos represents a relational str_pos expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.strpos`
// * During filtering: `relational_query.filter.scalar.strpos`
// * During sorting:`relational_query.sort.expression.scalar.strpos`
// * During joining: `relational_query.join.expression.scalar.strpos`
// * During aggregation: `relational_query.aggregate.expression.scalar.strpos`
// * During windowing: `relational_query.window.expression.scalar.strpos`.
type RelationalExpressionStrPos struct {
	Str    RelationalExpression `json:"str"    mapstructure:"str"    yaml:"str"`
	Substr RelationalExpression `json:"substr" mapstructure:"substr" yaml:"substr"`
}

// NewRelationalExpressionStrPos creates a RelationalExpressionStrPos instance.
func NewRelationalExpressionStrPos[S RelationalExpressionInner, SS RelationalExpressionInner](
	str S,
	substr SS,
) *RelationalExpressionStrPos {
	return &RelationalExpressionStrPos{
		Str:    str.Wrap(),
		Substr: substr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionStrPos) Type() RelationalExpressionType {
	return RelationalExpressionTypeStrPos
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionStrPos) ToMap() map[string]any {
	return map[string]any{
		"type":   j.Type(),
		"str":    j.Str,
		"substr": j.Substr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionStrPos) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionSubstr represents a relational substr expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.substr`
// * During filtering: `relational_query.filter.scalar.substr`
// * During sorting:`relational_query.sort.expression.scalar.substr`
// * During joining: `relational_query.join.expression.scalar.substr`
// * During aggregation: `relational_query.aggregate.expression.scalar.substr`
// * During windowing: `relational_query.window.expression.scalar.substr`.
type RelationalExpressionSubstr struct {
	Str      RelationalExpression  `json:"str"           mapstructure:"str"           yaml:"str"`
	StartPos RelationalExpression  `json:"start_pos"     mapstructure:"start_pos"     yaml:"start_pos"`
	Len      *RelationalExpression `json:"len,omitempty" mapstructure:"len,omitempty" yaml:"len,omitempty"`
}

// NewRelationalExpressionSubstr creates a RelationalExpressionSubstr instance.
func NewRelationalExpressionSubstr[S RelationalExpressionInner, SP RelationalExpressionInner](
	str S,
	startPos SP,
) *RelationalExpressionSubstr {
	return &RelationalExpressionSubstr{
		Str:      str.Wrap(),
		StartPos: startPos.Wrap(),
	}
}

// WithLen returns the new instance with len set.
func (j RelationalExpressionSubstr) WithLen(
	length RelationalExpressionInner,
) *RelationalExpressionSubstr {
	if length == nil {
		j.Len = nil
	} else {
		l := length.Wrap()
		j.Len = &l
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionSubstr) Type() RelationalExpressionType {
	return RelationalExpressionTypeSubstr
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionSubstr) ToMap() map[string]any {
	result := map[string]any{
		"type":      j.Type(),
		"str":       j.Str,
		"start_pos": j.StartPos,
	}

	if j.Len != nil && !j.Len.IsEmpty() {
		result["len"] = j.Len
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionSubstr) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionSubstrIndex represents a relational substr_index expression type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.substr_index`
// * During filtering: `relational_query.filter.scalar.substr_index`
// * During sorting:`relational_query.sort.expression.scalar.substr_index`
// * During joining: `relational_query.join.expression.scalar.substr_index`
// * During aggregation: `relational_query.aggregate.expression.scalar.substr_index`
// * During windowing: `relational_query.window.expression.scalar.substr_index`.
type RelationalExpressionSubstrIndex struct {
	Str   RelationalExpression `json:"str"   mapstructure:"str"   yaml:"str"`
	Delim RelationalExpression `json:"delim" mapstructure:"delim" yaml:"delim"`
	Count RelationalExpression `json:"count" mapstructure:"count" yaml:"count"`
}

// NewRelationalExpressionSubstrIndex creates a RelationalExpressionSubstrIndex instance.
func NewRelationalExpressionSubstrIndex[S RelationalExpressionInner, D RelationalExpressionInner, C RelationalExpressionInner](
	str S,
	delim D,
	count C,
) *RelationalExpressionSubstrIndex {
	return &RelationalExpressionSubstrIndex{
		Str:   str.Wrap(),
		Delim: delim.Wrap(),
		Count: count.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionSubstrIndex) Type() RelationalExpressionType {
	return RelationalExpressionTypeSubstrIndex
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionSubstrIndex) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"str":   j.Str,
		"delim": j.Delim,
		"count": j.Count,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionSubstrIndex) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionTan represents a RelationalExpression with the tan type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.tan`
// * During filtering: `relational_query.filter.scalar.tan`
// * During sorting:`relational_query.sort.expression.scalar.tan`
// * During joining: `relational_query.join.expression.scalar.tan`
// * During aggregation: `relational_query.aggregate.expression.scalar.tan`
// * During windowing: `relational_query.window.expression.scalar.tan`.
type RelationalExpressionTan struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionTan creates a RelationalExpressionTan instance.
func NewRelationalExpressionTan[E RelationalExpressionInner](expr E) *RelationalExpressionTan {
	return &RelationalExpressionTan{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionTan) Type() RelationalExpressionType {
	return RelationalExpressionTypeTan
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionTan) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionTan) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionToDate represents a RelationalExpression with the to_date type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.to_date`
// * During filtering: `relational_query.filter.scalar.to_date`
// * During sorting:`relational_query.sort.expression.scalar.to_date`
// * During joining: `relational_query.join.expression.scalar.to_date`
// * During aggregation: `relational_query.aggregate.expression.scalar.to_date`
// * During windowing: `relational_query.window.expression.scalar.to_date`.
type RelationalExpressionToDate struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionToDate creates a RelationalExpressionToDate instance.
func NewRelationalExpressionToDate[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionToDate {
	return &RelationalExpressionToDate{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionToDate) Type() RelationalExpressionType {
	return RelationalExpressionTypeToDate
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionToDate) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionToDate) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionToTimestamp represents a RelationalExpression with the to_timestamp type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.to_timestamp`
// * During filtering: `relational_query.filter.scalar.to_timestamp`
// * During sorting:`relational_query.sort.expression.scalar.to_timestamp`
// * During joining: `relational_query.join.expression.scalar.to_timestamp`
// * During aggregation: `relational_query.aggregate.expression.scalar.to_timestamp`
// * During windowing: `relational_query.window.expression.scalar.to_timestamp`.
type RelationalExpressionToTimestamp struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionToTimestamp creates a RelationalExpressionToTimestamp instance.
func NewRelationalExpressionToTimestamp[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionToTimestamp {
	return &RelationalExpressionToTimestamp{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionToTimestamp) Type() RelationalExpressionType {
	return RelationalExpressionTypeToTimestamp
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionToTimestamp) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionToTimestamp) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionTrunc represents a RelationalExpression with the trunc type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.trunc`
// * During filtering: `relational_query.filter.scalar.trunc`
// * During sorting:`relational_query.sort.expression.scalar.trunc`
// * During joining: `relational_query.join.expression.scalar.trunc`
// * During aggregation: `relational_query.aggregate.expression.scalar.trunc`
// * During windowing: `relational_query.window.expression.scalar.trunc`.
type RelationalExpressionTrunc struct {
	Expr RelationalExpression  `json:"expr"           mapstructure:"expr"           yaml:"expr"`
	Prec *RelationalExpression `json:"prec,omitempty" mapstructure:"prec,omitempty" yaml:"prec,omitempty"`
}

// NewRelationalExpressionTrunc creates a RelationalExpressionTrunc instance.
func NewRelationalExpressionTrunc[E RelationalExpressionInner](str E) *RelationalExpressionTrunc {
	return &RelationalExpressionTrunc{
		Expr: str.Wrap(),
	}
}

// WithPrec returns the new instance with prec set.
func (j RelationalExpressionTrunc) WithPrec(
	prec RelationalExpressionInner,
) *RelationalExpressionTrunc {
	if prec == nil {
		j.Prec = nil
	} else {
		ts := prec.Wrap()
		j.Prec = &ts
	}

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionTrunc) Type() RelationalExpressionType {
	return RelationalExpressionTypeTrunc
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionTrunc) ToMap() map[string]any {
	result := map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}

	if j.Prec != nil && !j.Prec.IsEmpty() {
		result["prec"] = j.Prec
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionTrunc) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionToLower represents a RelationalExpression with the to_lower type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.to_lower`
// * During filtering: `relational_query.filter.scalar.to_lower`
// * During sorting:`relational_query.sort.expression.scalar.to_lower`
// * During joining: `relational_query.join.expression.scalar.to_lower`
// * During aggregation: `relational_query.aggregate.expression.scalar.to_lower`
// * During windowing: `relational_query.window.expression.scalar.to_lower`.
type RelationalExpressionToLower struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionToLower creates a RelationalExpressionToLower instance.
func NewRelationalExpressionToLower[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionToLower {
	return &RelationalExpressionToLower{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionToLower) Type() RelationalExpressionType {
	return RelationalExpressionTypeToLower
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionToLower) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionToLower) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionToUpper represents a RelationalExpression with the to_upper type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.to_upper`
// * During filtering: `relational_query.filter.scalar.to_upper`
// * During sorting:`relational_query.sort.expression.scalar.to_upper`
// * During joining: `relational_query.join.expression.scalar.to_upper`
// * During aggregation: `relational_query.aggregate.expression.scalar.to_upper`
// * During windowing: `relational_query.window.expression.scalar.to_upper`.
type RelationalExpressionToUpper struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionToUpper creates a RelationalExpressionToUpper instance.
func NewRelationalExpressionToUpper[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionToUpper {
	return &RelationalExpressionToUpper{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionToUpper) Type() RelationalExpressionType {
	return RelationalExpressionTypeToUpper
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionToUpper) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionToUpper) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionBinaryConcat represents a RelationalExpression with the binary_concat type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.scalar.binary_concat`
// * During filtering: `relational_query.filter.scalar.binary_concat`
// * During sorting:`relational_query.sort.expression.scalar.binary_concat`
// * During joining: `relational_query.join.expression.scalar.binary_concat`
// * During aggregation: `relational_query.aggregate.expression.scalar.binary_concat`
// * During windowing: `relational_query.window.expression.scalar.binary_concat`.
type RelationalExpressionBinaryConcat struct {
	Left  RelationalExpression `json:"left"  mapstructure:"left"  yaml:"left"`
	Right RelationalExpression `json:"right" mapstructure:"right" yaml:"right"`
}

// NewRelationalExpressionBinaryConcat creates a RelationalExpressionBinaryConcat instance.
func NewRelationalExpressionBinaryConcat[L RelationalExpressionInner, R RelationalExpressionInner](
	left L,
	right R,
) *RelationalExpressionBinaryConcat {
	return &RelationalExpressionBinaryConcat{
		Left:  left.Wrap(),
		Right: right.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionBinaryConcat) Type() RelationalExpressionType {
	return RelationalExpressionTypeBinaryConcat
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionBinaryConcat) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"left":  j.Left,
		"right": j.Right,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionBinaryConcat) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionAverage represents a RelationalExpression with the average type.
type RelationalExpressionAverage struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionAverage creates a RelationalExpressionAverage instance.
func NewRelationalExpressionAverage[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionAverage {
	return &RelationalExpressionAverage{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionAverage) Type() RelationalExpressionType {
	return RelationalExpressionTypeAverage
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionAverage) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionAverage) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionBoolAnd represents a RelationalExpression with the bool_and type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.bool_and`
// * During filtering: `relational_query.filter.aggregate.bool_and`
// * During sorting:`relational_query.sort.expression.aggregate.bool_and`
// * During joining: `relational_query.join.expression.aggregate.bool_and`
// * During aggregation: `relational_query.aggregate.expression.aggregate.bool_and`
// * During windowing: `relational_query.window.expression.aggregate.bool_and`.
type RelationalExpressionBoolAnd struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionBoolAnd creates a RelationalExpressionBoolAnd instance.
func NewRelationalExpressionBoolAnd[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionBoolAnd {
	return &RelationalExpressionBoolAnd{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionBoolAnd) Type() RelationalExpressionType {
	return RelationalExpressionTypeBoolAnd
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionBoolAnd) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionBoolAnd) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionBoolOr represents a RelationalExpression with the bool_or type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.bool_or`
// * During filtering: `relational_query.filter.aggregate.bool_or`
// * During sorting:`relational_query.sort.expression.aggregate.bool_or`
// * During joining: `relational_query.join.expression.aggregate.bool_or`
// * During aggregation: `relational_query.aggregate.expression.aggregate.bool_or`
// * During windowing: `relational_query.window.expression.aggregate.bool_or`.
type RelationalExpressionBoolOr struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionBoolOr creates a RelationalExpressionBoolOr instance.
func NewRelationalExpressionBoolOr[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionBoolOr {
	return &RelationalExpressionBoolOr{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionBoolOr) Type() RelationalExpressionType {
	return RelationalExpressionTypeBoolOr
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionBoolOr) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionBoolOr) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCount represents a RelationalExpression with the count type.
type RelationalExpressionCount struct {
	Expr RelationalExpression `json:"expr"     mapstructure:"expr"     yaml:"expr"`
	// Only used when in specific contexts where the appropriate capability is supported:
	// * During projection: `relational_query.project.expression.aggregate.count.distinct`
	// * During filtering: `relational_query.filter.aggregate.count.distinct`
	// * During sorting:`relational_query.sort.expression.aggregate.count.distinct`
	// * During joining: `relational_query.join.expression.aggregate.count.distinct`
	// * During aggregation: `relational_query.aggregate.expression.aggregate.count.distinct`
	// * During windowing: `relational_query.window.expression.aggregate.count.distinct`
	Distinct bool `json:"distinct" mapstructure:"distinct" yaml:"distinct"`
}

// NewRelationalExpressionCount creates a RelationalExpressionCount instance.
func NewRelationalExpressionCount[E RelationalExpressionInner](
	expr E,
	distinct bool,
) *RelationalExpressionCount {
	return &RelationalExpressionCount{
		Expr:     expr.Wrap(),
		Distinct: distinct,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionCount) Type() RelationalExpressionType {
	return RelationalExpressionTypeCount
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionCount) ToMap() map[string]any {
	return map[string]any{
		"type":     j.Type(),
		"expr":     j.Expr,
		"distinct": j.Distinct,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCount) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionFirstValue represents a RelationalExpression with the first_value type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.first_value`
// * During filtering: `relational_query.filter.aggregate.first_value`
// * During sorting:`relational_query.sort.expression.aggregate.first_value`
// * During joining: `relational_query.join.expression.aggregate.first_value`
// * During aggregation: `relational_query.aggregate.expression.aggregate.first_value`
// * During windowing: `relational_query.window.expression.aggregate.first_value`.
type RelationalExpressionFirstValue struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionFirstValue creates a RelationalExpressionFirstValue instance.
func NewRelationalExpressionFirstValue[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionFirstValue {
	return &RelationalExpressionFirstValue{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionFirstValue) Type() RelationalExpressionType {
	return RelationalExpressionTypeFirstValue
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionFirstValue) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionFirstValue) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionLastValue represents a RelationalExpression with the last_value type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.last_value`
// * During filtering: `relational_query.filter.aggregate.last_value`
// * During sorting:`relational_query.sort.expression.aggregate.last_value`
// * During joining: `relational_query.join.expression.aggregate.last_value`
// * During aggregation: `relational_query.aggregate.expression.aggregate.last_value`
// * During windowing: `relational_query.window.expression.aggregate.last_value`.
type RelationalExpressionLastValue struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionLastValue creates a RelationalExpressionLastValue instance.
func NewRelationalExpressionLastValue[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionLastValue {
	return &RelationalExpressionLastValue{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionLastValue) Type() RelationalExpressionType {
	return RelationalExpressionTypeLastValue
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionLastValue) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionLastValue) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionMax represents a RelationalExpression with the max type.
type RelationalExpressionMax struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionMax creates a RelationalExpressionMax instance.
func NewRelationalExpressionMax[E RelationalExpressionInner](expr E) *RelationalExpressionMax {
	return &RelationalExpressionMax{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionMax) Type() RelationalExpressionType {
	return RelationalExpressionTypeMax
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionMax) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionMax) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionMedian represents a RelationalExpression with the median type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.median`
// * During filtering: `relational_query.filter.aggregate.median`
// * During sorting:`relational_query.sort.expression.aggregate.median`
// * During joining: `relational_query.join.expression.aggregate.median`
// * During aggregation: `relational_query.aggregate.expression.aggregate.median`
// * During windowing: `relational_query.window.expression.aggregate.median`.
type RelationalExpressionMedian struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionMedian creates a RelationalExpressionMedian instance.
func NewRelationalExpressionMedian[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionMedian {
	return &RelationalExpressionMedian{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionMedian) Type() RelationalExpressionType {
	return RelationalExpressionTypeMedian
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionMedian) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionMedian) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionMin represents a RelationalExpression with the min type.
type RelationalExpressionMin struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionMin creates a RelationalExpressionMin instance.
func NewRelationalExpressionMin[E RelationalExpressionInner](expr E) *RelationalExpressionMin {
	return &RelationalExpressionMin{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionMin) Type() RelationalExpressionType {
	return RelationalExpressionTypeMin
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionMin) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionMin) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionStringAgg represents a RelationalExpression with the string_agg type.
type RelationalExpressionStringAgg struct {
	Expr      RelationalExpression `json:"expr"               mapstructure:"expr"               yaml:"expr"`
	Separator string               `json:"separator"          mapstructure:"separator"          yaml:"separator"`
	// Only used when in specific contexts where the appropriate capability is supported:
	// * During projection: `relational_query.project.expression.aggregate.string_agg.distinct`
	// * During filtering: `relational_query.filter.aggregate.string_agg.distinct`
	// * During sorting:`relational_query.sort.expression.aggregate.string_agg.distinct`
	// * During joining: `relational_query.join.expression.aggregate.string_agg.distinct`
	// * During aggregation: `relational_query.aggregate.expression.aggregate.string_agg.distinct` * During windowing: `relational_query.window.expression.aggregate.string_agg.distinct`
	Distinct bool `json:"distinct"           mapstructure:"distinct"           yaml:"distinct"`
	// Only used when in specific contexts where the appropriate capability is supported:
	// * During projection: `relational_query.project.expression.aggregate.string_agg.order_by`
	// * During filtering: `relational_query.filter.aggregate.string_agg.order_by`
	// * During sorting:`relational_query.sort.expression.aggregate.string_agg.order_by`
	// * During joining: `relational_query.join.expression.aggregate.string_agg.order_by`
	// * During aggregation: `relational_query.aggregate.expression.aggregate.string_agg.order_by`
	// * During windowing: `relational_query.window.expression.aggregate.string_agg.order_by`
	OrderBy []Sort `json:"order_by,omitempty" mapstructure:"order_by,omitempty" yaml:"order_by,omitempty"`
}

// NewRelationalExpressionStringAgg creates a RelationalExpressionStringAgg instance.
func NewRelationalExpressionStringAgg[E RelationalExpressionInner](
	expr E,
	separator string,
	distinct bool,
) *RelationalExpressionStringAgg {
	return &RelationalExpressionStringAgg{
		Expr:      expr.Wrap(),
		Separator: separator,
		Distinct:  distinct,
	}
}

// WithOrderBy returns the RelationalExpressionStringAgg with order by.
func (j RelationalExpressionStringAgg) WithOrderBy(orderBy []Sort) *RelationalExpressionStringAgg {
	j.OrderBy = orderBy

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionStringAgg) Type() RelationalExpressionType {
	return RelationalExpressionTypeStringAgg
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionStringAgg) ToMap() map[string]any {
	result := map[string]any{
		"type":      j.Type(),
		"expr":      j.Expr,
		"separator": j.Separator,
		"distinct":  j.Distinct,
	}

	if j.OrderBy != nil {
		result["order_by"] = j.OrderBy
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionStringAgg) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionSum represents a RelationalExpression with the sum type.
type RelationalExpressionSum struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionSum creates a RelationalExpressionSum instance.
func NewRelationalExpressionSum[E RelationalExpressionInner](expr E) *RelationalExpressionSum {
	return &RelationalExpressionSum{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionSum) Type() RelationalExpressionType {
	return RelationalExpressionTypeSum
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionSum) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionSum) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionVar represents a RelationalExpression with the var type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.var`
// * During filtering: `relational_query.filter.aggregate.var`
// * During sorting:`relational_query.sort.expression.aggregate.var`
// * During joining: `relational_query.join.expression.aggregate.var`
// * During aggregation: `relational_query.aggregate.expression.aggregate.var`
// * During windowing: `relational_query.window.expression.aggregate.var`.
type RelationalExpressionVar struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionVar creates a RelationalExpressionVar instance.
func NewRelationalExpressionVar[E RelationalExpressionInner](expr E) *RelationalExpressionVar {
	return &RelationalExpressionVar{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionVar) Type() RelationalExpressionType {
	return RelationalExpressionTypeVar
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionVar) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionVar) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionStddev represents a RelationalExpression with the stddev type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.stddev`
// * During filtering: `relational_query.filter.aggregate.stddev`
// * During sorting:`relational_query.sort.expression.aggregate.stddev`
// * During joining: `relational_query.join.expression.aggregate.stddev`
// * During aggregation: `relational_query.aggregate.expression.aggregate.stddev`
// * During windowing: `relational_query.window.expression.aggregate.stddev`.
type RelationalExpressionStddev struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionStddev creates a RelationalExpressionStddev instance.
func NewRelationalExpressionStddev[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionStddev {
	return &RelationalExpressionStddev{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionStddev) Type() RelationalExpressionType {
	return RelationalExpressionTypeStddev
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionStddev) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionStddev) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionStddevPop represents a RelationalExpression with the stddev_pop type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.stddev_pop`
// * During filtering: `relational_query.filter.aggregate.stddev_pop`
// * During sorting:`relational_query.sort.expression.aggregate.stddev_pop`
// * During joining: `relational_query.join.expression.aggregate.stddev_pop`
// * During aggregation: `relational_query.aggregate.expression.aggregate.stddev_pop`
// * During windowing: `relational_query.window.expression.aggregate.stddev_pop`.
type RelationalExpressionStddevPop struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionStddevPop creates a RelationalExpressionStddevPop instance.
func NewRelationalExpressionStddevPop[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionStddevPop {
	return &RelationalExpressionStddevPop{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionStddevPop) Type() RelationalExpressionType {
	return RelationalExpressionTypeStddevPop
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionStddevPop) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionStddevPop) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionApproxPercentileCont represents a RelationalExpression with the approx_percentile_cont type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.approx_percentile_cont`
// * During filtering: `relational_query.filter.aggregate.approx_percentile_cont`
// * During sorting:`relational_query.sort.expression.aggregate.approx_percentile_cont`
// * During joining: `relational_query.join.expression.aggregate.approx_percentile_cont`
// * During aggregation: `relational_query.aggregate.expression.aggregate.approx_percentile_cont`
// * During windowing: `relational_query.window.expression.aggregate.approx_percentile_cont`.
type RelationalExpressionApproxPercentileCont struct {
	Expr       RelationalExpression `json:"expr"       mapstructure:"expr"       yaml:"expr"`
	Percentile float64              `json:"percentile" mapstructure:"percentile" yaml:"percentile"`
}

// NewRelationalExpressionApproxPercentileCont creates a RelationalExpressionApproxPercentileCont instance.
func NewRelationalExpressionApproxPercentileCont[E RelationalExpressionInner](
	expr E,
	percentile float64,
) *RelationalExpressionApproxPercentileCont {
	return &RelationalExpressionApproxPercentileCont{
		Expr:       expr.Wrap(),
		Percentile: percentile,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionApproxPercentileCont) Type() RelationalExpressionType {
	return RelationalExpressionTypeApproxPercentileCont
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionApproxPercentileCont) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionApproxPercentileCont) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionArrayAgg represents a RelationalExpression with the array_agg type.
type RelationalExpressionArrayAgg struct {
	Expr RelationalExpression `json:"expr"               mapstructure:"expr"               yaml:"expr"`
	// Only used when in specific contexts where the appropriate capability is supported:
	// * During projection: `relational_query.project.expression.aggregate.array_agg.distinct`
	// * During filtering: `relational_query.filter.aggregate.array_agg.distinct`
	// * During sorting:`relational_query.sort.expression.aggregate.array_agg.distinct`
	// * During joining: `relational_query.join.expression.aggregate.array_agg.distinct`
	// * During aggregation: `relational_query.aggregate.expression.aggregate.array_agg.distinct`
	// * During windowing: `relational_query.window.expression.aggregate.array_agg.distinct`
	Distinct bool `json:"distinct"           mapstructure:"distinct"           yaml:"distinct"`
	// Only used when in specific contexts where the appropriate capability is supported:
	// * During projection: `relational_query.project.expression.aggregate.array_agg.order_by`
	// * During filtering: `relational_query.filter.aggregate.array_agg.order_by`
	// * During sorting:`relational_query.sort.expression.aggregate.array_agg.order_by`
	// * During joining: `relational_query.join.expression.aggregate.array_agg.order_by`
	// * During aggregation: `relational_query.aggregate.expression.aggregate.array_agg.order_by`
	// * During windowing: `relational_query.window.expression.aggregate.array_agg.order_by`
	OrderBy []Sort `json:"order_by,omitempty" mapstructure:"order_by,omitempty" yaml:"order_by,omitempty"`
}

// NewRelationalExpressionArrayAgg creates a RelationalExpressionArrayAgg instance.
func NewRelationalExpressionArrayAgg[E RelationalExpressionInner](
	expr E,
	distinct bool,
) *RelationalExpressionArrayAgg {
	return &RelationalExpressionArrayAgg{
		Expr:     expr.Wrap(),
		Distinct: distinct,
	}
}

// WithOrderBy returns the RelationalExpressionStringAgg with order by.
func (j RelationalExpressionArrayAgg) WithOrderBy(orderBy []Sort) *RelationalExpressionArrayAgg {
	j.OrderBy = orderBy

	return &j
}

// Type return the type name of the instance.
func (j RelationalExpressionArrayAgg) Type() RelationalExpressionType {
	return RelationalExpressionTypeArrayAgg
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionArrayAgg) ToMap() map[string]any {
	result := map[string]any{
		"type":     j.Type(),
		"expr":     j.Expr,
		"distinct": j.Distinct,
	}

	if j.OrderBy != nil {
		result["order_by"] = j.OrderBy
	}

	return result
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionArrayAgg) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionApproxDistinct represents a RelationalExpression with the approx_distinct type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.aggregate.approx_distinct`
// * During filtering: `relational_query.filter.aggregate.approx_distinct`
// * During sorting:`relational_query.sort.expression.aggregate.approx_distinct`
// * During joining: `relational_query.join.expression.aggregate.approx_distinct`
// * During aggregation: `relational_query.aggregate.expression.aggregate.approx_distinct`
// * During windowing: `relational_query.window.expression.aggregate.approx_distinct`.
type RelationalExpressionApproxDistinct struct {
	Expr RelationalExpression `json:"expr" mapstructure:"expr" yaml:"expr"`
}

// NewRelationalExpressionApproxDistinct creates a RelationalExpressionApproxDistinct instance.
func NewRelationalExpressionApproxDistinct[E RelationalExpressionInner](
	expr E,
) *RelationalExpressionApproxDistinct {
	return &RelationalExpressionApproxDistinct{
		Expr: expr.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionApproxDistinct) Type() RelationalExpressionType {
	return RelationalExpressionTypeApproxDistinct
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionApproxDistinct) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
		"expr": j.Expr,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionApproxDistinct) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionRowNumber represents a RelationalExpression with the row_number type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.window.row_number`
// * During filtering: `relational_query.filter.window.row_number`
// * During sorting:`relational_query.sort.expression.window.row_number`
// * During joining: `relational_query.join.expression.window.row_number`
// * During aggregation: `relational_query.window.row_number`
// * During windowing: `relational_query.window.expression.window.row_number`.
type RelationalExpressionRowNumber struct {
	OrderBy     []Sort                 `json:"order_by"     mapstructure:"order_by"     yaml:"order_by"`
	PartitionBy []RelationalExpression `json:"partition_by" mapstructure:"partition_by" yaml:"partition_by"`
}

// NewRelationalExpressionRowNumber creates a RelationalExpressionRowNumber instance.
func NewRelationalExpressionRowNumber(
	orderBy []Sort,
	partitionBy []RelationalExpressionInner,
) *RelationalExpressionRowNumber {
	rb := []RelationalExpression{}

	for _, r := range partitionBy {
		if r != nil {
			rb = append(rb, r.Wrap())
		}
	}

	return &RelationalExpressionRowNumber{
		OrderBy:     orderBy,
		PartitionBy: rb,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionRowNumber) Type() RelationalExpressionType {
	return RelationalExpressionTypeRowNumber
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionRowNumber) ToMap() map[string]any {
	return map[string]any{
		"type":         j.Type(),
		"order_by":     j.OrderBy,
		"partition_by": j.PartitionBy,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionRowNumber) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionDenseRank represents a RelationalExpression with the dense_rank type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.window.dense_rank`
// * During filtering: `relational_query.filter.window.dense_rank`
// * During sorting:`relational_query.sort.expression.window.dense_rank`
// * During joining: `relational_query.join.expression.window.dense_rank`
// * During aggregation: `relational_query.window.dense_rank`
// * During windowing: `relational_query.window.expression.window.dense_rank`.
type RelationalExpressionDenseRank struct {
	OrderBy     []Sort                 `json:"order_by"     mapstructure:"order_by"     yaml:"order_by"`
	PartitionBy []RelationalExpression `json:"partition_by" mapstructure:"partition_by" yaml:"partition_by"`
}

// NewRelationalExpressionDenseRank creates a RelationalExpressionDenseRank instance.
func NewRelationalExpressionDenseRank(
	orderBy []Sort,
	partitionBy []RelationalExpressionInner,
) *RelationalExpressionDenseRank {
	rb := []RelationalExpression{}

	for _, r := range partitionBy {
		if r != nil {
			rb = append(rb, r.Wrap())
		}
	}

	return &RelationalExpressionDenseRank{
		OrderBy:     orderBy,
		PartitionBy: rb,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionDenseRank) Type() RelationalExpressionType {
	return RelationalExpressionTypeDenseRank
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionDenseRank) ToMap() map[string]any {
	return map[string]any{
		"type":         j.Type(),
		"order_by":     j.OrderBy,
		"partition_by": j.PartitionBy,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionDenseRank) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionNTile represents a RelationalExpression with the ntile type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.window.ntile`
// * During filtering: `relational_query.filter.window.ntile`
// * During sorting:`relational_query.sort.expression.window.ntile`
// * During joining: `relational_query.join.expression.window.ntile`
// * During aggregation: `relational_query.window.ntile`
// * During windowing: `relational_query.window.expression.window.ntile`.
type RelationalExpressionNTile struct {
	OrderBy     []Sort                 `json:"order_by"     mapstructure:"order_by"     yaml:"order_by"`
	PartitionBy []RelationalExpression `json:"partition_by" mapstructure:"partition_by" yaml:"partition_by"`
	N           int64                  `json:"n"            mapstructure:"n"            yaml:"n"`
}

// NewRelationalExpressionNTile creates a RelationalExpressionNTile instance.
func NewRelationalExpressionNTile(
	orderBy []Sort,
	partitionBy []RelationalExpressionInner,
	n int64,
) *RelationalExpressionNTile {
	rb := []RelationalExpression{}

	for _, r := range partitionBy {
		if r != nil {
			rb = append(rb, r.Wrap())
		}
	}

	return &RelationalExpressionNTile{
		OrderBy:     orderBy,
		PartitionBy: rb,
		N:           n,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionNTile) Type() RelationalExpressionType {
	return RelationalExpressionTypeNTile
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionNTile) ToMap() map[string]any {
	return map[string]any{
		"type":         j.Type(),
		"order_by":     j.OrderBy,
		"partition_by": j.PartitionBy,
		"n":            j.N,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionNTile) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionRank represents a RelationalExpression with the rank type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.window.rank`
// * During filtering: `relational_query.filter.window.rank`
// * During sorting:`relational_query.sort.expression.window.rank`
// * During joining: `relational_query.join.expression.window.rank`
// * During aggregation: `relational_query.window.rank`
// * During windowing: `relational_query.window.expression.window.rank`.
type RelationalExpressionRank struct {
	OrderBy     []Sort                 `json:"order_by"     mapstructure:"order_by"     yaml:"order_by"`
	PartitionBy []RelationalExpression `json:"partition_by" mapstructure:"partition_by" yaml:"partition_by"`
}

// NewRelationalExpressionRank creates a RelationalExpressionRank instance.
func NewRelationalExpressionRank(
	orderBy []Sort,
	partitionBy []RelationalExpressionInner,
) *RelationalExpressionRank {
	rb := []RelationalExpression{}

	for _, r := range partitionBy {
		if r != nil {
			rb = append(rb, r.Wrap())
		}
	}

	return &RelationalExpressionRank{
		OrderBy:     orderBy,
		PartitionBy: rb,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionRank) Type() RelationalExpressionType {
	return RelationalExpressionTypeRank
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionRank) ToMap() map[string]any {
	return map[string]any{
		"type":         j.Type(),
		"order_by":     j.OrderBy,
		"partition_by": j.PartitionBy,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionRank) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionCumeDist represents a RelationalExpression with the cume_dist type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.window.cume_dist`
// * During filtering: `relational_query.filter.window.cume_dist`
// * During sorting:`relational_query.sort.expression.window.cume_dist`
// * During joining: `relational_query.join.expression.window.cume_dist`
// * During aggregation: `relational_query.window.cume_dist`
// * During windowing: `relational_query.window.expression.window.cume_dist`.
type RelationalExpressionCumeDist struct {
	OrderBy     []Sort                 `json:"order_by"     mapstructure:"order_by"     yaml:"order_by"`
	PartitionBy []RelationalExpression `json:"partition_by" mapstructure:"partition_by" yaml:"partition_by"`
}

// NewRelationalExpressionCumeDist creates a RelationalExpressionCumeDist instance.
func NewRelationalExpressionCumeDist(
	orderBy []Sort,
	partitionBy []RelationalExpressionInner,
) *RelationalExpressionCumeDist {
	rb := []RelationalExpression{}

	for _, r := range partitionBy {
		if r != nil {
			rb = append(rb, r.Wrap())
		}
	}

	return &RelationalExpressionCumeDist{
		OrderBy:     orderBy,
		PartitionBy: rb,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionCumeDist) Type() RelationalExpressionType {
	return RelationalExpressionTypeCumeDist
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionCumeDist) ToMap() map[string]any {
	return map[string]any{
		"type":         j.Type(),
		"order_by":     j.OrderBy,
		"partition_by": j.PartitionBy,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionCumeDist) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}

// RelationalExpressionPercentRank represents a RelationalExpression with the percent_rank type.
// Only used when in specific contexts where the appropriate capability is supported:
// * During projection: `relational_query.project.expression.window.percent_rank`
// * During filtering: `relational_query.filter.window.percent_rank`
// * During sorting:`relational_query.sort.expression.window.percent_rank`
// * During joining: `relational_query.join.expression.window.percent_rank`
// * During aggregation: `relational_query.window.percent_rank`
// * During windowing: `relational_query.window.expression.window.percent_rank`.
type RelationalExpressionPercentRank struct {
	OrderBy     []Sort                 `json:"order_by"     mapstructure:"order_by"     yaml:"order_by"`
	PartitionBy []RelationalExpression `json:"partition_by" mapstructure:"partition_by" yaml:"partition_by"`
}

// NewRelationalExpressionPercentRank creates a RelationalExpressionPercentRank instance.
func NewRelationalExpressionPercentRank(
	orderBy []Sort,
	partitionBy []RelationalExpressionInner,
) *RelationalExpressionPercentRank {
	rb := []RelationalExpression{}

	for _, r := range partitionBy {
		if r != nil {
			rb = append(rb, r.Wrap())
		}
	}

	return &RelationalExpressionPercentRank{
		OrderBy:     orderBy,
		PartitionBy: rb,
	}
}

// Type return the type name of the instance.
func (j RelationalExpressionPercentRank) Type() RelationalExpressionType {
	return RelationalExpressionTypePercentRank
}

// ToMap converts the instance to raw Field.
func (j RelationalExpressionPercentRank) ToMap() map[string]any {
	return map[string]any{
		"type":         j.Type(),
		"order_by":     j.OrderBy,
		"partition_by": j.PartitionBy,
	}
}

// Wrap returns the relation wrapper.
func (j RelationalExpressionPercentRank) Wrap() RelationalExpression {
	return NewRelationalExpression(&j)
}
