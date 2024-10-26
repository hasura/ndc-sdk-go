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

type Scalar struct {
	Schema     schema.ScalarType
	NativeType *TypeInfo
}

type Type interface {
	Kind() schema.TypeEnum
	Schema() schema.TypeEncoder
	SchemaName() string
	FullName() string
	String() string
}

type NullableType struct {
	UnderlyingType Type
}

var _ Type = &NullableType{}

func NewNullableType(input Type) *NullableType {
	return &NullableType{input}
}

func (t *NullableType) Kind() schema.TypeEnum {
	return schema.TypeNullable
}

func (t NullableType) SchemaName() string {
	return t.UnderlyingType.SchemaName()
}

func (t *NullableType) Schema() schema.TypeEncoder {
	if t.UnderlyingType.Kind() == schema.TypeNullable {
		return t.UnderlyingType.Schema()
	}
	return schema.NewNullableType(t.UnderlyingType.Schema())
}

func (t NullableType) FullName() string {
	return "*" + t.UnderlyingType.FullName()
}

func (t NullableType) String() string {
	return "*" + t.UnderlyingType.String()
}

type ArrayType struct {
	ElementType Type
}

var _ Type = &ArrayType{}

func NewArrayType(input Type) *ArrayType {
	return &ArrayType{input}
}

func (t *ArrayType) Kind() schema.TypeEnum {
	return schema.TypeArray
}

func (t *ArrayType) Schema() schema.TypeEncoder {
	return schema.NewArrayType(t.ElementType.Schema())
}

func (t ArrayType) SchemaName() string {
	return t.ElementType.SchemaName()
}

func (t ArrayType) FullName() string {
	return "[]" + t.ElementType.FullName()
}

func (t ArrayType) String() string {
	return "[]" + t.ElementType.String()
}

type NamedType struct {
	Name       string
	NativeType *TypeInfo
}

var _ Type = &NamedType{}

func NewNamedType(name string, info *TypeInfo) *NamedType {
	return &NamedType{name, info}
}

func (t *NamedType) Kind() schema.TypeEnum {
	return schema.TypeNamed
}

func (t *NamedType) Schema() schema.TypeEncoder {
	return schema.NewNamedType(t.Name)
}

func (t NamedType) SchemaName() string {
	return t.Name
}

func (t NamedType) FullName() string {
	return t.NativeType.GetAbsoluteName()
}

func (t *NamedType) String() string {
	return t.NativeType.String()
}

// TypeInfo represents the serialization information of a type.
type TypeInfo struct {
	Name           string
	IsAnonymous    bool
	TypeParameters []Type
	SchemaName     string
	Description    *string
	PackagePath    string
	PackageName    string
	TypeAST        types.Type
}

// CanMethod checks whether generating decoding methods for this type.
func (ti *TypeInfo) CanMethod() bool {
	return len(ti.TypeParameters) == 0
}

// GetArgumentName returns the argument name.
func (ti *TypeInfo) GetArgumentName(packagePath string) string {
	return ti.getArgumentName(packagePath, false)
}

// GetAbsoluteName return the type name with absolute package paths.
func (ti *TypeInfo) GetAbsoluteName() string {
	return ti.getArgumentName(ti.PackagePath, true)
}

// String implements the fmt.Stringer interface.
func (ti *TypeInfo) String() string {
	return ti.GetArgumentName(ti.PackagePath)
}

// String implements the fmt.Stringer interface.
func (ti *TypeInfo) GetPackagePaths(currentPackagePath string) []string {
	results := make([]string, 0)

	if ti.PackagePath != "" && ti.PackagePath != currentPackagePath {
		results = append(results, ti.PackagePath)
	}

	for _, param := range ti.TypeParameters {
		results = append(results, getTypePackagePaths(param, currentPackagePath)...)
	}

	return results
}

func (ti *TypeInfo) getArgumentName(packagePath string, isAbsolute bool) string {
	name := ti.Name
	if isAbsolute {
		if ti.PackagePath != "" {
			name = ti.PackagePath + "." + ti.Name
		}
	} else if ti.PackagePath != "" && packagePath != ti.PackagePath {
		name = ti.PackageName + "." + ti.Name
	}

	paramLen := len(ti.TypeParameters)
	if paramLen > 0 {
		name += "["
		for i, param := range ti.TypeParameters {
			name += getTypeArgumentName(param, packagePath, isAbsolute)
			if i < paramLen-1 {
				name += ", "
			}
		}
		name += "]"
	}

	return name
}

// Field represents the serialization information of a field.
type Field struct {
	Name        string
	Description *string
	Embedded    bool
	Type        Type
}

// ObjectInfo represents the serialization information of an object type.
type ObjectInfo struct {
	Description  *string
	Type         *TypeInfo
	Fields       map[string]Field
	SchemaFields schema.ObjectTypeFields
}

func (oi ObjectInfo) Schema() *schema.ObjectType {
	result := &schema.ObjectType{
		Description: oi.Description,
		Fields:      oi.SchemaFields,
	}
	return result
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
	ResultType    *Field
}

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function schema.
type FunctionInfo OperationInfo

// Schema returns a NDC function schema.
func (op FunctionInfo) Schema() schema.FunctionInfo {
	result := schema.FunctionInfo{
		Name:        op.Name,
		Description: op.Description,
		ResultType:  op.ResultType.Type.Schema().Encode(),
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
		ResultType:  op.ResultType.Type.Schema().Encode(),
		Arguments:   schema.ProcedureInfoArguments(op.Arguments),
	}
	return result
}

// RawConnectorSchema represents a readable Go schema object
// which can encode to NDC schema.
type RawConnectorSchema struct {
	StateType         *TypeInfo
	Imports           map[string]bool
	Scalars           map[string]Scalar
	Objects           map[string]ObjectInfo
	Functions         []FunctionInfo
	FunctionArguments map[string]ObjectInfo
	Procedures        []ProcedureInfo
}

// NewRawConnectorSchema creates an empty RawConnectorSchema instance.
func NewRawConnectorSchema() *RawConnectorSchema {
	return &RawConnectorSchema{
		Imports:           make(map[string]bool),
		Scalars:           make(map[string]Scalar),
		Objects:           make(map[string]ObjectInfo),
		Functions:         []FunctionInfo{},
		FunctionArguments: make(map[string]ObjectInfo),
		Procedures:        []ProcedureInfo{},
	}
}

// Schema converts to a NDC schema.
func (rcs RawConnectorSchema) Schema() *schema.SchemaResponse {
	result := &schema.SchemaResponse{
		ScalarTypes: schema.SchemaResponseScalarTypes{},
		ObjectTypes: schema.SchemaResponseObjectTypes{},
		Collections: []schema.CollectionInfo{},
		Functions:   []schema.FunctionInfo{},
		Procedures:  []schema.ProcedureInfo{},
	}

	for key, item := range rcs.Scalars {
		result.ScalarTypes[key] = item.Schema
	}
	for _, obj := range rcs.Objects {
		result.ObjectTypes[obj.Type.SchemaName] = *obj.Schema()
	}
	for _, function := range rcs.Functions {
		result.Functions = append(result.Functions, function.Schema())
	}
	for _, procedure := range rcs.Procedures {
		result.Procedures = append(result.Procedures, procedure.Schema())
	}

	return result
}

func (rcs *RawConnectorSchema) SetScalar(name string, value Scalar) {
	if rcs.Scalars == nil {
		rcs.Scalars = map[string]Scalar{}
	}
	_, ok := rcs.Scalars[name]
	if !ok {
		rcs.Scalars[name] = value
	}
}

func (rcs *RawConnectorSchema) setFunctionArgument(info ObjectInfo) {
	key := info.Type.String()
	if _, ok := rcs.FunctionArguments[key]; ok {
		return
	}
	rcs.FunctionArguments[key] = info
}

func (rcs RawConnectorSchema) GetScalarFromType(ty Type) *Scalar {
	switch t := ty.(type) {
	case *NullableType:
		return rcs.GetScalarFromType(t.UnderlyingType)
	case *ArrayType:
		return rcs.GetScalarFromType(t.ElementType)
	case *NamedType:
		result, ok := rcs.Scalars[t.Name]
		if ok {
			return &result
		}
	}
	return nil
}

func getTypeArgumentName(input Type, packagePath string, isAbsolute bool) string {
	switch t := input.(type) {
	case *NullableType:
		return "*" + getTypeArgumentName(t.UnderlyingType, packagePath, isAbsolute)
	case *ArrayType:
		return "[]" + getTypeArgumentName(t.ElementType, packagePath, isAbsolute)
	case *NamedType:
		return t.NativeType.getArgumentName(packagePath, isAbsolute)
	default:
		panic(fmt.Errorf("getTypeArgumentName: invalid type %v", input))
	}
}

func getTypePackagePaths(input Type, currentPackagePath string) []string {
	switch t := input.(type) {
	case *NullableType:
		return getTypePackagePaths(t.UnderlyingType, currentPackagePath)
	case *ArrayType:
		return getTypePackagePaths(t.ElementType, currentPackagePath)
	case *NamedType:
		return t.NativeType.GetPackagePaths(currentPackagePath)
	default:
		panic(fmt.Errorf("getTypePackagePaths: invalid type %v", input))
	}
}

func unwrapNullableType(input Type) (Type, bool) {
	switch t := input.(type) {
	case *NullableType:
		result, _ := unwrapNullableType(t.UnderlyingType)
		return result, true
	case *ArrayType:
		return t, false
	case *NamedType:
		return t, false
	default:
		panic(fmt.Errorf("unwrapNullableType: invalid type %v", input))
	}
}
