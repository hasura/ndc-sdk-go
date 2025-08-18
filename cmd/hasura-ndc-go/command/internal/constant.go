package internal

import (
	"errors"
	"regexp"

	"github.com/hasura/ndc-sdk-go/v2/schema"
)

type ScalarName string

const (
	ScalarBoolean        ScalarName = "Boolean"
	ScalarString         ScalarName = "String"
	ScalarInt8           ScalarName = "Int8"
	ScalarInt16          ScalarName = "Int16"
	ScalarInt32          ScalarName = "Int32"
	ScalarInt64          ScalarName = "Int64"
	ScalarFloat32        ScalarName = "Float32"
	ScalarFloat64        ScalarName = "Float64"
	ScalarBigInt         ScalarName = "BigInt"
	ScalarBigDecimal     ScalarName = "BigDecimal"
	ScalarUUID           ScalarName = "UUID"
	ScalarDate           ScalarName = "Date"
	ScalarDuration       ScalarName = "Duration"
	ScalarDurationString ScalarName = "DurationString"
	ScalarDurationInt64  ScalarName = "DurationInt64"
	ScalarTimestamp      ScalarName = "Timestamp"
	ScalarTimestampTZ    ScalarName = "TimestampTZ"
	ScalarGeography      ScalarName = "Geography"
	ScalarBytes          ScalarName = "Bytes"
	ScalarJSON           ScalarName = "JSON"
	// ScalarRawJSON is a special scalar for raw json data serialization.
	// The underlying Go type for this scalar is json.RawMessage.
	// Note: we don't recommend to use this scalar for function arguments
	// because the decoder will re-encode the value to []byte that isn't performance-wise.
	ScalarRawJSON ScalarName = "RawJSON"
	ScalarURL     ScalarName = "URL"
)

var defaultScalarTypes = map[ScalarName]schema.ScalarType{
	ScalarBoolean: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationBoolean().Encode(),
	},
	ScalarString: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationString().Encode(),
	},
	ScalarInt8: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationInt8().Encode(),
	},
	ScalarInt16: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationInt16().Encode(),
	},
	ScalarInt32: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationInt32().Encode(),
	},
	ScalarInt64: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationInt64().Encode(),
	},
	ScalarFloat32: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationFloat32().Encode(),
	},
	ScalarFloat64: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationFloat64().Encode(),
	},
	ScalarBigInt: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationBigInteger().Encode(),
	},
	ScalarBigDecimal: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationBigDecimal().Encode(),
	},
	ScalarUUID: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationUUID().Encode(),
	},
	ScalarDate: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationDate().Encode(),
	},
	ScalarDuration: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationJSON().Encode(),
	},
	ScalarDurationString: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationString().Encode(),
	},
	ScalarDurationInt64: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationInt64().Encode(),
	},
	ScalarTimestamp: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationTimestamp().Encode(),
	},
	ScalarTimestampTZ: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationTimestampTZ().Encode(),
	},
	ScalarGeography: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationGeography().Encode(),
	},
	ScalarBytes: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationBytes().Encode(),
	},
	ScalarJSON: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationJSON().Encode(),
	},
	ScalarRawJSON: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationJSON().Encode(),
	},
	ScalarURL: {
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
		Representation:      schema.NewTypeRepresentationString().Encode(),
	},
}

var (
	fieldNameRegex           = regexp.MustCompile(`[^\w]`)
	ndcOperationNameRegex    = regexp.MustCompile(`^(Function|Procedure)([A-Z][A-Za-z0-9]*)$`)
	ndcOperationCommentRegex = regexp.MustCompile(`^@(function|procedure)(\s+([A-Za-z]\w*))?`)
	ndcScalarNameRegex       = regexp.MustCompile(`^Scalar([A-Z]\w*)$`)
	ndcScalarCommentRegex    = regexp.MustCompile(`^@scalar(\s+(\w+))?(\s+([a-z]+))?\.?$`)
	ndcEnumCommentRegex      = regexp.MustCompile(`^@enum\s+([\w-.,!@#$%^&*()+=~\s\t]+)$`)
)

var (
	errUnsupportedTypeDuration = errors.New(
		"unsupported time.Duration. Use github.com/hasura/ndc-sdk-go/v2/scalar.Duration instead",
	)
	errMustUseEnumTag = errors.New("use @enum tag with values instead")
)

const (
	packageSDKUtils  = "github.com/hasura/ndc-sdk-go/v2/utils"
	packageSDKSchema = "github.com/hasura/ndc-sdk-go/v2/schema"
)
