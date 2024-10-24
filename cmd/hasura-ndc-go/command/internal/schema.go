package internal

import (
	"fmt"
	"go/types"

	"github.com/hasura/ndc-sdk-go/schema"
)

type OperationKind string

const (
	OperationFunction  OperationKind = "Function"
	OperationProcedure OperationKind = "Procedure"
)

type TypeKind string

// TypeInfo represents the serialization information of a type.
type TypeInfo struct {
	Name                 string
	TypeParameters       []TypeInfo
	SchemaName           string
	Description          *string
	PackagePath          string
	PackageName          string
	Embedded             bool
	IsScalar             bool
	ScalarRepresentation schema.TypeRepresentation
	TypeFragments        []string
	TypeAST              types.Type
	Schema               schema.TypeEncoder
}

// IsNullable checks if the current type is nullable.
func (ti *TypeInfo) IsNullable() bool {
	return isNullableFragments(ti.TypeFragments)
}

// IsArray checks if the current type is an array.
func (ti *TypeInfo) IsArray() bool {
	return isArrayFragments(ti.TypeFragments)
}

// CanMethod checks whether generating decoding methods for this type.
func (ti *TypeInfo) CanMethod() bool {
	return len(ti.TypeParameters) == 0
}

// GetArgumentName returns the argument name.
func (ti *TypeInfo) GetArgumentName(packagePath string) string {
	name := ti.Name
	if ti.PackagePath != "" && packagePath != ti.PackagePath {
		name = fmt.Sprintf("%s.%s", ti.PackageName, ti.Name)
	}

	paramLen := len(ti.TypeParameters)
	if paramLen > 0 {
		name += "["
		for i, param := range ti.TypeParameters {
			name += param.GetArgumentName(packagePath)
			if i < paramLen-1 {
				name += ", "
			}
		}
		name += "]"
	}

	return name
}

// String implements the fmt.Stringer interface.
func (ti *TypeInfo) String() string {
	return ti.GetArgumentName(ti.PackagePath)
}

func isNullableFragment(fragment string) bool {
	return fragment == "*"
}

func isNullableFragments(fragments []string) bool {
	return len(fragments) > 0 && isNullableFragment(fragments[0])
}

func isArrayFragment(fragment string) bool {
	return fragment == "[]"
}

func isArrayFragments(fragments []string) bool {
	return len(fragments) > 0 && isArrayFragment(fragments[0])
}

// ObjectField represents the serialization information of an object field.
type ObjectField struct {
	Name        string
	Description *string
	Type        *TypeInfo
}

// ObjectInfo represents the serialization information of an object type.
type ObjectInfo struct {
	IsAnonymous bool
	Type        *TypeInfo
	Fields      map[string]ObjectField
}

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function or procedure schema.
type OperationInfo struct {
	Kind          OperationKind
	Name          string
	OriginName    string
	PackageName   string
	PackagePath   string
	Description   *string
	ArgumentsType *TypeInfo
	Arguments     map[string]schema.ArgumentInfo
	ResultType    *TypeInfo
}

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function schema.
type FunctionInfo OperationInfo

// Schema returns a NDC function schema.
func (op FunctionInfo) Schema() schema.FunctionInfo {
	result := schema.FunctionInfo{
		Name:        op.Name,
		Description: op.Description,
		ResultType:  op.ResultType.Schema.Encode(),
		Arguments:   op.Arguments,
	}
	return result
}

// ProcedureInfo represents a readable Go function info
// which can convert to a NDC procedure schema.
type ProcedureInfo FunctionInfo

// Schema returns a NDC procedure schema.
func (op ProcedureInfo) Schema() schema.ProcedureInfo {
	result := schema.ProcedureInfo{
		Name:        op.Name,
		Description: op.Description,
		ResultType:  op.ResultType.Schema.Encode(),
		Arguments:   schema.ProcedureInfoArguments(op.Arguments),
	}
	return result
}

// RawConnectorSchema represents a readable Go schema object
// which can encode to NDC schema.
type RawConnectorSchema struct {
	StateType         *TypeInfo
	Imports           map[string]bool
	CustomScalars     map[string]*TypeInfo
	ScalarSchemas     schema.SchemaResponseScalarTypes
	Objects           map[string]ObjectInfo
	ObjectSchemas     schema.SchemaResponseObjectTypes
	Functions         []FunctionInfo
	FunctionArguments map[string]ObjectInfo
	Procedures        []ProcedureInfo
}

// NewRawConnectorSchema creates an empty RawConnectorSchema instance.
func NewRawConnectorSchema() *RawConnectorSchema {
	return &RawConnectorSchema{
		Imports:           make(map[string]bool),
		CustomScalars:     make(map[string]*TypeInfo),
		ScalarSchemas:     make(schema.SchemaResponseScalarTypes),
		Objects:           make(map[string]ObjectInfo),
		ObjectSchemas:     make(schema.SchemaResponseObjectTypes),
		Functions:         []FunctionInfo{},
		FunctionArguments: make(map[string]ObjectInfo),
		Procedures:        []ProcedureInfo{},
	}
}

// Schema converts to a NDC schema.
func (rcs RawConnectorSchema) Schema() *schema.SchemaResponse {
	result := &schema.SchemaResponse{
		ScalarTypes: rcs.ScalarSchemas,
		ObjectTypes: rcs.ObjectSchemas,
		Collections: []schema.CollectionInfo{},
		Functions:   []schema.FunctionInfo{},
		Procedures:  []schema.ProcedureInfo{},
	}
	for _, function := range rcs.Functions {
		result.Functions = append(result.Functions, function.Schema())
	}
	for _, procedure := range rcs.Procedures {
		result.Procedures = append(result.Procedures, procedure.Schema())
	}

	return result
}

func (rcs RawConnectorSchema) setFunctionArgument(info ObjectInfo) {
	key := info.Type.String()
	if _, ok := rcs.FunctionArguments[key]; ok {
		return
	}
	rcs.FunctionArguments[key] = info
}
