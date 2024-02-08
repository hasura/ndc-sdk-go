package main

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
)

var defaultScalarTypes = schema.SchemaResponseScalarTypes{
	"String": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
	"Int": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
	"Float": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
	"Boolean": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
}

// get scalar type name go native type name
func getScalarTypeNameFromNativeTypeName(name string) (string, bool) {
	switch name {
	case "bool":
		return "Boolean", true
	case "string":
		return "String", true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
		return "Int", true
	case "float32", "float64", "complex64", "complex128":
		return "Float", true
	default:
		return "", false
	}
}

// check if the input type name is a scalar
func isTypeNameScalar(input string) bool {
	_, ok := getScalarTypeNameFromNativeTypeName(input)
	return ok
}

// get scalar type name go native type
func getScalarTypeNameFromNativeType(input any) (string, bool) {
	switch input.(type) {
	case bool, *bool:
		return "Boolean", true
	case string, *string:
		return "String", true
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *uintptr:
		return "Int", true
	case float32, float64, complex64, complex128, *float32, *float64, *complex64, *complex128:
		return "Float", true
	default:
		return "", false
	}
}

// check if the input type is a scalar
func isScalar(input any) bool {
	_, ok := getScalarTypeNameFromNativeType(input)
	return ok
}

func getInnerReturnExprType(ty ast.Expr, skipNullable bool) (schema.TypeEncoder, error) {
	switch t := ty.(type) {
	case *ast.StarExpr:
		if !skipNullable {
			innerType, err := getInnerReturnExprType(t.X, false)
			if err != nil {
				return nil, err
			}
			return schema.NewNullableType(innerType), nil
		}
		return getInnerReturnExprType(t.X, false)
	case *ast.Ident:
		log.Printf("obj: %+v", *t.Obj)
		return schema.NewNamedType(t.Name), nil
	default:
		return nil, fmt.Errorf("unhandled expr %+v: %+v", ty, t)
	}
}

func parseOperationInfoFromComment(functionName string, comments []*ast.Comment) *OperationInfo {
	var result OperationInfo
	var descriptions []string
	for _, comment := range comments {
		text := strings.TrimSpace(strings.TrimLeft(comment.Text, "/"))
		matches := operationNameCommentRegex.FindStringSubmatch(text)
		matchesLen := len(matches)
		if matchesLen > 1 {
			switch matches[1] {
			case string(OperationFunction):
				result.Kind = OperationFunction
			case string(OperationProcedure):
				result.Kind = OperationProcedure
			default:
				log.Println("unsupported operation kind:", matches[0])
			}

			if matchesLen > 3 && strings.TrimSpace(matches[3]) != "" {
				result.Name = strings.TrimSpace(matches[3])
			} else {
				result.Name = strings.ToLower(functionName[:1]) + functionName[1:]
			}
		} else {
			descriptions = append(descriptions, text)
		}
	}

	if result.Kind == "" {
		return nil
	}

	result.Description = strings.TrimSpace(strings.Join(descriptions, " "))
	return &result
}
