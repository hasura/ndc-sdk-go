package internal

import (
	"context"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/trace"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"golang.org/x/tools/go/packages"
)

type SchemaParser struct {
	context      context.Context
	moduleName   string
	rawSchema    *RawConnectorSchema
	packages     []*packages.Package
	packageIndex int
	namingStyle  OperationNamingStyle
}

// GetCurrentPackage gets the current evaluating package.
func (sp SchemaParser) GetCurrentPackage() *packages.Package {
	return sp.packages[sp.packageIndex]
}

// FindPackageByPath finds the package by package path.
func (sp SchemaParser) FindPackageByPath(input string) *packages.Package {
	for _, p := range sp.packages {
		if p.ID == input {
			return p
		}
	}
	return nil
}

func parseRawConnectorSchemaFromGoCode(ctx context.Context, moduleName string, filePath string, args *ConnectorGenerationArguments) (*RawConnectorSchema, error) {
	var err error
	namingStyle := StyleCamelCase
	if args.Style != "" {
		namingStyle, err = ParseOperationNamingStyle(args.Style)
		if err != nil {
			return nil, err
		}
	}
	rawSchema := NewRawConnectorSchema()

	tempDirs := args.Directories
	if len(args.Directories) == 0 {
		// recursively walk directories if the user don't explicitly specify target folders
		entries, err := os.ReadDir(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read subdirectories of %s: %w", filePath, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			tempDirs = append(tempDirs, entry.Name())
		}
	}
	var directories []string
	for _, dir := range tempDirs {
		for _, globPath := range []string{path.Join(filePath, dir, "*.go"), path.Join(filePath, dir, "**", "*.go")} {
			goFiles, err := filepath.Glob(globPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read subdirectories of %s/%s: %w", filePath, dir, err)
			}
			// cleanup types.generated.go files
			fileCount := 0
			for _, fp := range goFiles {
				if !strings.HasSuffix(fp, typeMethodsOutputFile) {
					fileCount++
					continue
				}
				if err := os.Remove(fp); err != nil {
					return nil, fmt.Errorf("failed to delete %s: %w", fp, err)
				}
			}
			if fileCount > 0 {
				directories = append(directories, dir)
				break
			}
		}
	}

	if len(directories) > 0 {
		log.Info().Interface("directories", directories).Msgf("parsing connector schema...")

		var packageList []*packages.Package
		fset := token.NewFileSet()
		for _, folder := range directories {
			_, parseCodeTask := trace.NewTask(ctx, fmt.Sprintf("parse_%s_code", folder))
			folderPath := path.Join(filePath, folder)
			cfg := &packages.Config{
				Mode: packages.NeedSyntax | packages.NeedTypes,
				Dir:  folderPath,
				Fset: fset,
			}
			pkgList, err := packages.Load(cfg, flag.Args()...)
			parseCodeTask.End()
			if err != nil {
				return nil, err
			}
			packageList = append(packageList, pkgList...)
		}

		for i := range packageList {
			parseSchemaCtx, parseSchemaTask := trace.NewTask(ctx, "parse_schema_"+packageList[i].ID)
			sp := &SchemaParser{
				context:      parseSchemaCtx,
				moduleName:   moduleName,
				packages:     packageList,
				packageIndex: i,
				rawSchema:    rawSchema,
				namingStyle:  namingStyle,
			}

			err = sp.parseRawConnectorSchema(packageList[i].Types)
			parseSchemaTask.End()
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Info().Msgf("no subdirectory in %s", filePath)
	}

	if rawSchema.StateType != nil {
		rawSchema.Imports[rawSchema.StateType.PackagePath] = true
	} else {
		pkgPathTypes, err := evalPackageTypesLocation(moduleName, filePath, args.ConnectorDir)
		if err != nil {
			return nil, err
		}
		rawSchema.StateType = &TypeInfo{
			Name:        "State",
			PackagePath: pkgPathTypes,
			PackageName: "types",
		}
		rawSchema.Imports[rawSchema.StateType.PackagePath] = true
	}
	return rawSchema, nil
}

// parse raw connector schema from Go code.
func (sp *SchemaParser) parseRawConnectorSchema(pkg *types.Package) error {
	for _, name := range pkg.Scope().Names() {
		_, task := trace.NewTask(sp.context, fmt.Sprintf("parse_%s_schema_%s", sp.GetCurrentPackage().Name, name))
		err := sp.parsePackageScope(pkg, name)
		task.End()
		if err != nil {
			return err
		}
	}

	return nil
}

func (sp *SchemaParser) parsePackageScope(pkg *types.Package, name string) error {
	switch obj := pkg.Scope().Lookup(name).(type) {
	case *types.Func:
		// only parse public functions
		if !obj.Exported() {
			return nil
		}
		opInfo := sp.parseOperationInfo(obj)
		if opInfo == nil {
			return nil
		}
		opInfo.PackageName = pkg.Name()
		opInfo.PackagePath = pkg.Path()
		var resultTuple *types.Tuple
		var params *types.Tuple
		switch sig := obj.Type().(type) {
		case *types.Signature:
			params = sig.Params()
			resultTuple = sig.Results()
		default:
			return fmt.Errorf("expected function signature, got: %s", sig.String())
		}

		if params == nil || (params.Len() < 2 || params.Len() > 3) {
			return fmt.Errorf("%s: expect 2 or 3 parameters only (ctx context.Context, state *types.State, arguments *[ArgumentType]), got %s", opInfo.OriginName, params)
		}

		if resultTuple == nil || resultTuple.Len() != 2 {
			return fmt.Errorf("%s: expect result tuple ([type], error), got %s", opInfo.OriginName, resultTuple)
		}

		if sp.rawSchema.StateType == nil {
			ty := sp.getNamedType(params.At(1).Type())
			if ty != nil {
				so := ty.Obj()
				if so != nil {
					objPkg := so.Pkg()
					if objPkg != nil {
						sp.rawSchema.StateType = &TypeInfo{
							Name:        so.Name(),
							PackageName: objPkg.Name(),
							PackagePath: objPkg.Path(),
						}
					}
				}
			}
		}

		// parse arguments in the function if exists
		// ignore 2 first parameters (context and state)
		if params.Len() == 3 {
			arg := params.At(2)
			argumentParser := NewTypeParser(sp, &Field{}, arg.Type(), &opInfo.Kind)
			argumentInfo, err := argumentParser.ParseArgumentTypes([]string{})
			if err != nil {
				return err
			}
			opInfo.ArgumentsType = argumentInfo.Type
			if opInfo.Kind == OperationFunction {
				sp.rawSchema.setFunctionArgument(*argumentInfo)
			}
			// convert argument schema
			for k, a := range argumentInfo.SchemaFields {
				opInfo.Arguments[k] = schema.ArgumentInfo{
					Description: a.Description,
					Type:        a.Type,
				}
			}
		}

		typeParser := NewTypeParser(sp, &Field{}, resultTuple.At(0).Type(), nil)
		resultType, err := typeParser.Parse([]string{})
		if err != nil {
			return err
		}
		opInfo.ResultType = resultType

		switch opInfo.Kind {
		case OperationProcedure:
			sp.rawSchema.Procedures = append(sp.rawSchema.Procedures, ProcedureInfo(*opInfo))
		case OperationFunction:
			sp.rawSchema.Functions = append(sp.rawSchema.Functions, FunctionInfo(*opInfo))
		}
	default:
	}
	return nil
}

func (sp *SchemaParser) getNamedType(ty types.Type) *types.Named {
	switch t := ty.(type) {
	case *types.Pointer:
		return sp.getNamedType(t.Elem())
	case *types.Named:
		return t
	case *types.Slice:
		return sp.getNamedType(t.Elem())
	case *types.Array:
		return sp.getNamedType(t.Elem())
	default:
		return nil
	}
}

// format operation name with style.
func (sp SchemaParser) formatOperationName(name string) string {
	switch sp.namingStyle {
	case StyleSnakeCase:
		return strcase.ToSnake(name)
	default:
		return strcase.ToLowerCamel(name)
	}
}

func (sp *SchemaParser) parseOperationInfo(fn *types.Func) *OperationInfo {
	functionName := fn.Name()
	result := OperationInfo{
		OriginName: functionName,
		Arguments:  make(map[string]schema.ArgumentInfo),
	}

	var descriptions []string
	commentGroup := findCommentsFromPos(sp.GetCurrentPackage(), fn.Scope(), functionName)
	if commentGroup != nil {
		for i, comment := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(comment.Text, "/"))

			// trim the function name in the first line if exists
			if i == 0 {
				text = strings.TrimPrefix(text, functionName+" ")
			}
			matches := ndcOperationCommentRegex.FindStringSubmatch(text)
			matchesLen := len(matches)
			if matchesLen > 1 {
				switch matches[1] {
				case strings.ToLower(string(OperationFunction)):
					result.Kind = OperationFunction
				case strings.ToLower(string(OperationProcedure)):
					result.Kind = OperationProcedure
				default:
					log.Debug().Msgf("unsupported operation kind: %s", matches)
				}

				if matchesLen > 3 && strings.TrimSpace(matches[3]) != "" {
					result.Name = strings.TrimSpace(matches[3])
				} else {
					result.Name = sp.formatOperationName(functionName)
				}
			} else {
				descriptions = append(descriptions, text)
			}
		}
	}

	// try to parse function with following prefixes:
	// - FunctionXxx as a query function
	// - ProcedureXxx as a mutation procedure
	if result.Kind == "" {
		operationNameResults := ndcOperationNameRegex.FindStringSubmatch(functionName)
		if len(operationNameResults) < 3 {
			return nil
		}
		result.Kind = OperationKind(operationNameResults[1])
		result.Name = sp.formatOperationName(operationNameResults[2])
	}

	desc := strings.TrimSpace(strings.Join(descriptions, " "))
	if desc != "" {
		result.Description = &desc
	}

	return &result
}

func findCommentsFromPos(pkg *packages.Package, scope *types.Scope, name string) *ast.CommentGroup {
	if pkg == nil {
		return nil
	}

	for _, f := range pkg.Syntax {
		for _, cg := range f.Comments {
			if len(cg.List) == 0 {
				continue
			}
			exp := regexp.MustCompile(fmt.Sprintf(`^//\s+%s(\s|$)`, name))
			if !exp.MatchString(cg.List[0].Text) {
				continue
			}
			if _, obj := scope.LookupParent(name, cg.Pos()); obj != nil {
				return cg
			}
		}
	}
	return nil
}

func evalPackageTypesLocation(moduleName string, filePath string, connectorDir string) (string, error) {
	matches, err := filepath.Glob(path.Join(filePath, "types", "*.go"))
	if err == nil && len(matches) > 0 {
		return moduleName + "/types", nil
	}

	if connectorDir != "" && !strings.HasPrefix(connectorDir, ".") {
		matches, err = filepath.Glob(path.Join(filePath, connectorDir, "types", "*.go"))
		if err == nil && len(matches) > 0 {
			return fmt.Sprintf("%s/%s/types", moduleName, connectorDir), nil
		}
	}
	return "", fmt.Errorf("the `types` package where the State struct is in must be placed in root or connector directory, %w", err)
}
