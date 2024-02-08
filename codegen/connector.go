package main

import (
	"fmt"
	"go/ast"
	"log"
	"regexp"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
)

var (
	operationNameCommentRegex = regexp.MustCompile(`^@(function|procedure)(\s+([a-z][a-zA-Z0-9_]*))?`)
)

type OperationKind string

var (
	OperationFunction  OperationKind = "function"
	OperationProcedure OperationKind = "procedure"
)

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function or procedure schema
type OperationInfo struct {
	Kind        OperationKind
	Name        string
	Description string
	Parameters  []*ast.Field
	ResultType  schema.TypeEncoder
}

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function schema
type FunctionInfo OperationInfo

// Schema returns a NDC function schema
func (op FunctionInfo) Schema() schema.FunctionInfo {
	result := schema.FunctionInfo{
		Name:       op.Name,
		ResultType: op.ResultType.Encode(),
		Arguments:  schema.FunctionInfoArguments{},
	}
	if op.Description != "" {
		result.Description = &op.Description
	}
	return result
}

// ProcedureInfo represents a readable Go function info
// which can convert to a NDC procedure schema
type ProcedureInfo FunctionInfo

// Schema returns a NDC procedure schema
func (op ProcedureInfo) Schema() schema.ProcedureInfo {
	result := schema.ProcedureInfo{
		Name:       op.Name,
		ResultType: op.ResultType.Encode(),
		Arguments:  schema.ProcedureInfoArguments{},
	}
	if op.Description != "" {
		result.Description = &op.Description
	}
	return result
}

// RawConnectorSchema represents a readable Go schema object
// which can encode to NDC schema
type RawConnectorSchema struct {
	Scalars    schema.SchemaResponseScalarTypes
	Objects    []*ast.Object
	Functions  []FunctionInfo
	Procedures []ProcedureInfo
}

// Schema converts to a NDC schema
func (rcs RawConnectorSchema) Schema() *schema.SchemaResponse {
	result := &schema.SchemaResponse{
		ScalarTypes: rcs.Scalars,
		ObjectTypes: schema.SchemaResponseObjectTypes{},
		Collections: []schema.CollectionInfo{},
	}
	for _, function := range rcs.Functions {
		result.Functions = append(result.Functions, function.Schema())
	}
	for _, procedure := range rcs.Procedures {
		result.Procedures = append(result.Procedures, procedure.Schema())
	}

	return result
}

// parse raw connector schema from Go code
func parseRawConnectorSchema(packages map[string]*ast.Package) (*RawConnectorSchema, error) {
	functions := []FunctionInfo{}
	procedures := []ProcedureInfo{}
	objectTypes := []*ast.Object{}

	for _, pkg := range packages {
		for fileName, file := range pkg.Files {
			if file.Scope == nil {
				log.Printf("%s: empty scope", fileName)
				continue
			}
			for objName, obj := range file.Scope.Objects {
				switch obj.Kind {
				case ast.Typ:
					parseTypeDeclaration(obj)
					objectTypes = append(objectTypes, obj)
				case ast.Fun:
					if objName == "main" {
						continue
					}
					fnDecl, ok := obj.Decl.(*ast.FuncDecl)
					if !ok {
						continue
					}
					if fnDecl.Doc == nil {
						continue
					}

					opInfo := parseOperationInfoFromComment(objName, fnDecl.Doc.List)
					if opInfo == nil {
						continue
					}

					opInfo.Parameters = fnDecl.Type.Params.List
					if len(fnDecl.Type.Results.List) != 2 {
						return nil, fmt.Errorf("%s: invalid return result types, expect 2 results (<result>, error)", objName)
					}

					resultType, err := getInnerReturnExprType(fnDecl.Type.Results.List[0].Type, true)
					if err != nil {
						return nil, err
					}
					opInfo.ResultType = resultType
					switch opInfo.Kind {
					case OperationFunction:
						functions = append(functions, FunctionInfo(*opInfo))
					case OperationProcedure:
						procedures = append(procedures, ProcedureInfo(*opInfo))
					}
				}
			}
		}
	}

	return &RawConnectorSchema{
		Functions:  functions,
		Procedures: procedures,
	}, nil
}

func parseTypeDeclaration(object *ast.Object) (map[string]schema.ObjectType, map[string]schema.ScalarType) {
	spec, ok := object.Decl.(*ast.TypeSpec)
	if !ok {
		return nil, nil
	}
	log.Printf("name: %s, type: %+v", object.Name, spec.Type)

	var comment string
	if spec.Comment != nil {
		comment = strings.TrimSpace(spec.Comment.Text())
	}

	switch t := spec.Type.(type) {
	case *ast.StructType:
		log.Printf("obj: %+v", t.Fields)
		if t.Fields == nil || len(t.Fields.List) == 0 {
			return nil, nil
		}
		objType := schema.ObjectType{}
		if comment != "" {
			objType.Description = &comment
		}
		for _, field := range t.Fields.List {
			log.Printf("field name: %+v, doc: %+v, comment: %+v, type: %+v", field.Names, field.Doc, field.Comment, field.Type)
		}
		return nil, nil
	default:
		return nil, nil
	}
}

func mergeMap[K comparable, V any](dest map[K]V, src map[K]V) map[K]V {
	result := dest
	for k, v := range src {
		result[k] = v
	}
	return result
}
