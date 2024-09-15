package internal

import (
	_ "embed"
	"fmt"
	"regexp"
	"text/template"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/iancoleman/strcase"
)

const (
	connectorOutputFile   = "connector.generated.go"
	schemaOutputJSONFile  = "schema.generated.json"
	schemaOutputGoFile    = "schema.generated.go"
	typeMethodsOutputFile = "types.generated.go"
)

//go:embed templates/connector/connector.go.tmpl
var connectorTemplateStr string
var connectorTemplate *template.Template

func init() {
	var err error
	connectorTemplate, err = template.New(connectorOutputFile).Parse(connectorTemplateStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse connector template: %s", err))
	}

	strcase.ConfigureAcronym("API", "Api")
	strcase.ConfigureAcronym("REST", "Rest")
	strcase.ConfigureAcronym("HTTP", "Http")
	strcase.ConfigureAcronym("SQL", "sql")
}

type ScalarName string

const (
	ScalarBoolean     ScalarName = "Boolean"
	ScalarString      ScalarName = "String"
	ScalarInt8        ScalarName = "Int8"
	ScalarInt16       ScalarName = "Int16"
	ScalarInt32       ScalarName = "Int32"
	ScalarInt64       ScalarName = "Int64"
	ScalarFloat32     ScalarName = "Float32"
	ScalarFloat64     ScalarName = "Float64"
	ScalarBigInt      ScalarName = "BigInt"
	ScalarBigDecimal  ScalarName = "BigDecimal"
	ScalarUUID        ScalarName = "UUID"
	ScalarDate        ScalarName = "Date"
	ScalarTimestamp   ScalarName = "Timestamp"
	ScalarTimestampTZ ScalarName = "TimestampTZ"
	ScalarGeography   ScalarName = "Geography"
	ScalarBytes       ScalarName = "Bytes"
	ScalarJSON        ScalarName = "JSON"
	// ScalarRawJSON is a special scalar for raw json data serialization.
	// The underlying Go type for this scalar is json.RawMessage.
	// Note: we don't recommend to use this scalar for function arguments
	// because the decoder will re-encode the value to []byte that isn't performance-wise.
	ScalarRawJSON ScalarName = "RawJSON"
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
}

var ndcOperationNameRegex = regexp.MustCompile(`^(Function|Procedure)([A-Z][A-Za-z0-9]*)$`)
var ndcOperationCommentRegex = regexp.MustCompile(`^@(function|procedure)(\s+([A-Za-z]\w*))?`)
var ndcScalarNameRegex = regexp.MustCompile(`^Scalar([A-Z]\w*)$`)
var ndcScalarCommentRegex = regexp.MustCompile(`^@scalar(\s+(\w+))?(\s+([a-z]+))?$`)
var ndcEnumCommentRegex = regexp.MustCompile(`^@enum\s+([\w-.,!@#$%^&*()+=~\s\t]+)$`)

type nativeScalarPackageConfig struct {
	PackageName string
	Pattern     *regexp.Regexp
}

var fieldNameRegex = regexp.MustCompile(`[^\w]`)

var nativeScalarPackages map[string]nativeScalarPackageConfig = map[string]nativeScalarPackageConfig{
	"scalar": {
		PackageName: "github.com/hasura/ndc-sdk-go/scalar",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*github\.com\/hasura\/ndc-sdk-go\/scalar\.`),
	},
	"json": {
		PackageName: "encoding/json",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*encoding\/json\.`),
	},
	"uuid": {
		PackageName: "github.com/google/uuid",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*github\.com\/google\/uuid\.`),
	},
	"time": {
		PackageName: "time",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*time\.`),
	},
}

const textBlockErrorCheck = `
    if err != nil {
		  return err
    }
`

const textBlockErrorCheck2 = `
    if err != nil {
      return nil, err
    }
`
