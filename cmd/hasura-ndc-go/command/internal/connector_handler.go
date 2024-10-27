package internal

import (
	"fmt"
	"strings"
)

const (
	functionEnumsName  = "enumValues_FunctionName"
	procedureEnumsName = "enumValues_ProcedureName"
)

type connectorHandlerBuilder struct {
	RawSchema  *RawConnectorSchema
	Functions  []FunctionInfo
	Procedures []ProcedureInfo
	Builder    *connectorTypeBuilder
}

func (chb connectorHandlerBuilder) Render() {
	if len(chb.Functions) == 0 && len(chb.Procedures) == 0 {
		return
	}

	bs := chb.Builder
	bs.imports["context"] = ""
	bs.imports["log/slog"] = ""
	bs.imports["slices"] = ""
	bs.imports["github.com/hasura/ndc-sdk-go/connector"] = ""
	bs.imports["github.com/hasura/ndc-sdk-go/schema"] = ""
	bs.imports["go.opentelemetry.io/otel/trace"] = ""
	if chb.RawSchema.StateType != nil && bs.packagePath != chb.RawSchema.StateType.PackagePath {
		bs.imports[chb.RawSchema.StateType.PackagePath] = ""
	}

	_, _ = bs.builder.WriteString(`
// DataConnectorHandler implements the data connector handler 
type DataConnectorHandler struct{}
`)
	chb.writeQuery(bs.builder)
	chb.writeMutation(bs.builder)

	bs.builder.WriteString(`    
func connector_addSpanEvent(span trace.Span, logger *slog.Logger, name string, data map[string]any, options ...trace.EventOption) {
  logger.Debug(name, slog.Any("data", data))
  attrs := utils.DebugJSONAttributes(data, utils.IsDebug(logger))
  span.AddEvent(name, append(options, trace.WithAttributes(attrs...))...)
}`)
}

func (chb connectorHandlerBuilder) writeOperationNameEnums(sb *strings.Builder, name string, values []string) {
	sb.WriteString("var ")
	sb.WriteString(name)
	sb.WriteString(" = []string{")
	for i, enum := range values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteRune('"')
		sb.WriteString(enum)
		sb.WriteRune('"')
	}
	sb.WriteRune('}')
}

func (chb connectorHandlerBuilder) writeStateArgumentName() string {
	if chb.RawSchema.StateType == nil {
		return "State"
	}
	return chb.RawSchema.StateType.GetArgumentName(chb.Builder.packagePath)
}

func (chb connectorHandlerBuilder) writeQuery(sb *strings.Builder) {
	if len(chb.Functions) == 0 {
		return
	}
	stateArgument := chb.writeStateArgumentName()
	_, _ = sb.WriteString(`
// QueryExists check if the query name exists
func (dch DataConnectorHandler) QueryExists(name string) bool {
  return slices.Contains(`)
	_, _ = sb.WriteString(functionEnumsName)
	_, _ = sb.WriteString(`, name)
}`)

	_, _ = sb.WriteString(`
func (dch DataConnectorHandler) Query(ctx context.Context, state *`)
	_, _ = sb.WriteString(stateArgument)
	_, _ = sb.WriteString(`, request *schema.QueryRequest, rawArgs map[string]any) (*schema.RowSet, error) {
  if !dch.QueryExists(request.Collection) {
    return nil, utils.ErrHandlerNotfound
  }
  queryFields, err := utils.EvalFunctionSelectionFieldValue(request)
  if err != nil {
    return nil, schema.UnprocessableContentError(err.Error(), nil)
  }

  result, err := dch.execQuery(ctx, state, request, queryFields, rawArgs)
  if err != nil {
    return nil, err
  }
  
  return &schema.RowSet{
    Aggregates: schema.RowSetAggregates{},
    Rows: []map[string]any{
      {
        "__value": result,
      },
    },
  }, nil
}
  
func (dch DataConnectorHandler) execQuery(ctx context.Context, state *`)
	_, _ = sb.WriteString(stateArgument)
	_, _ = sb.WriteString(`, request *schema.QueryRequest, queryFields schema.NestedField, rawArgs map[string]any) (any, error) {
  span := trace.SpanFromContext(ctx)
  logger := connector.GetLogger(ctx)
  switch request.Collection {`)

	functionKeys := make([]string, len(chb.Functions))
	for i, fn := range chb.Functions {
		functionKeys[i] = fn.Name
		op := OperationInfo(fn)
		resultType, isNullable := unwrapNullableType(op.ResultType.Type)
		chb.writeOperationValidation(sb, &op, "queryFields", resultType)

		var argumentParamStr string
		if fn.ArgumentsType != nil {
			argName := fn.ArgumentsType.GetArgumentName(chb.Builder.packagePath)
			if fn.ArgumentsType.PackagePath != "" && fn.ArgumentsType.PackagePath != chb.Builder.packagePath {
				chb.Builder.imports[fn.ArgumentsType.PackagePath] = ""
			}

			sb.WriteString(`
    var args `)
			sb.WriteString(argName)
			sb.WriteString("\n    if parseErr := ")
			if fn.ArgumentsType.CanMethod() {
				sb.WriteString("args.FromValue(rawArgs)")
			} else {
				sb.WriteString("connector_Decoder.DecodeObject(&args, rawArgs)")
			}
			sb.WriteString(`; parseErr != nil {
      return nil, schema.UnprocessableContentError("failed to resolve arguments", map[string]any{
        "cause": parseErr.Error(),
      })
    }
    
    connector_addSpanEvent(span, logger, "execute_function", map[string]any{
      "arguments": args,
    })`)
			argumentParamStr = ", &args"
		}

		switch t := resultType.(type) {
		case *ArrayType:
			chb.writeOperationResult(sb, fn.OriginName, OperationFunction, argumentParamStr, isNullable)
			sb.WriteString("\n    result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)")
			writeErrorCheck(sb, 2, 4)
			sb.WriteString("    return result, nil\n")
		case *NamedType:
			if _, ok := chb.RawSchema.Scalars[t.NativeType.SchemaName]; ok {
				sb.WriteString(fmt.Sprintf("\n    return %s(ctx, state%s)\n", fn.OriginName, argumentParamStr))
			} else {
				chb.writeOperationResult(sb, fn.OriginName, OperationFunction, argumentParamStr, isNullable)
				sb.WriteString("\n    result, err := utils.EvalNestedColumnObject(selection, rawResult)")
				writeErrorCheck(sb, 2, 4)
				sb.WriteString("    return result, nil\n")
			}
		}
	}

	_, _ = sb.WriteString(`
  default:
    return nil, utils.ErrHandlerNotfound
  }
}
`)
	chb.writeOperationNameEnums(sb, functionEnumsName, functionKeys)
}

func (chb connectorHandlerBuilder) writeMutation(sb *strings.Builder) {
	if len(chb.Procedures) == 0 {
		return
	}
	stateArgument := chb.writeStateArgumentName()
	chb.Builder.imports["encoding/json"] = ""

	_, _ = sb.WriteString(`
// MutationExists check if the mutation name exists
func (dch DataConnectorHandler) MutationExists(name string) bool {
  return slices.Contains(`)
	_, _ = sb.WriteString(procedureEnumsName)
	_, _ = sb.WriteString(`, name)
}`)

	_, _ = sb.WriteString(`
func (dch DataConnectorHandler) Mutation(ctx context.Context, state *`)
	_, _ = sb.WriteString(stateArgument)
	_, _ = sb.WriteString(`, operation *schema.MutationOperation) (schema.MutationOperationResults, error) {
  span := trace.SpanFromContext(ctx)  
  logger := connector.GetLogger(ctx)
  connector_addSpanEvent(span, logger, "validate_request", map[string]any{
    "operations_name": operation.Name,
  })
  
  switch operation.Name {`)

	procedureKeys := make([]string, len(chb.Procedures))
	for i, fn := range chb.Procedures {
		procedureKeys[i] = fn.Name
		op := OperationInfo(fn)
		resultType, isNullable := unwrapNullableType(op.ResultType.Type)
		chb.writeOperationValidation(sb, &op, "operation.Fields", resultType)

		var argumentParamStr string
		if fn.ArgumentsType != nil {
			argName := fn.ArgumentsType.GetArgumentName(chb.Builder.packagePath)
			if fn.ArgumentsType.PackagePath != "" && fn.ArgumentsType.PackagePath != chb.Builder.packagePath {
				chb.Builder.imports[fn.ArgumentsType.PackagePath] = ""
			}

			argumentStr := fmt.Sprintf(`
    var args %s
    if err := json.Unmarshal(operation.Arguments, &args); err != nil {
      return nil, schema.UnprocessableContentError("failed to decode arguments", map[string]any{
        "cause": err.Error(),
      })
    }`, argName)
			sb.WriteString(argumentStr)
			argumentParamStr = ", &args"
		}

		sb.WriteString("\n    span.AddEvent(\"execute_procedure\")")

		switch t := resultType.(type) {
		case *ArrayType:
			chb.writeOperationResult(sb, fn.OriginName, OperationProcedure, argumentParamStr, isNullable)
			sb.WriteString("\n    result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)\n")
			writeErrorCheck(sb, 2, 4)
		case *NamedType:
			if _, ok := chb.RawSchema.Scalars[t.NativeType.SchemaName]; ok {
				chb.writeOperationExecution(sb, fn.OriginName, argumentParamStr, "result")
			} else {
				chb.writeOperationResult(sb, fn.OriginName, OperationProcedure, argumentParamStr, isNullable)
				sb.WriteString("\n    result, err := utils.EvalNestedColumnObject(selection, rawResult)\n")
				writeErrorCheck(sb, 2, 4)
			}
		}
		sb.WriteString("    return schema.NewProcedureResult(result).Encode(), nil\n")
	}

	_, _ = sb.WriteString(`
  default:
    return nil, utils.ErrHandlerNotfound
  }
}
`)
	chb.writeOperationNameEnums(sb, procedureEnumsName, procedureKeys)
}

func (chb connectorHandlerBuilder) writeOperationValidation(sb *strings.Builder, fn *OperationInfo, selector string, resultType Type) {
	_, _ = sb.WriteString("\n  case \"")
	_, _ = sb.WriteString(fn.Name)
	_, _ = sb.WriteString("\":\n")

	switch t := resultType.(type) {
	case *ArrayType:
		sb.WriteString("\n    selection, err := ")
		sb.WriteString(selector)
		sb.WriteString(`.AsArray()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be array", map[string]any{
        "cause": err.Error(),
      })
    }`)
	case *NamedType:
		if _, ok := chb.RawSchema.Scalars[t.NativeType.SchemaName]; ok {
			sb.WriteString("\n      if len(")
			sb.WriteString(selector)
			sb.WriteString(`) > 0 {
          return nil, schema.UnprocessableContentError("cannot evaluate selection fields for scalar", nil)
        }`)
			return
		}

		sb.WriteString("\n    selection, err := ")
		sb.WriteString(selector)
		sb.WriteString(`.AsObject()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
        "cause": err.Error(),
      })
    }`)
	}
}

func (chb connectorHandlerBuilder) writeOperationExecution(sb *strings.Builder, operationName string, argumentParamStr string, resultVarName string) {
	sb.WriteString("\n    ")
	sb.WriteString(resultVarName)
	sb.WriteString(", err := ")
	sb.WriteString(operationName)
	sb.WriteString("(ctx, state")
	sb.WriteString(argumentParamStr)
	sb.WriteString(")\n")
	writeErrorCheck(sb, 2, 4)
}

func (chb connectorHandlerBuilder) writeOperationResult(sb *strings.Builder, operationName string, operationKind OperationKind, argumentParamStr string, isNullable bool) {
	chb.writeOperationExecution(sb, operationName, argumentParamStr, "rawResult")
	if isNullable {
		sb.WriteString("\n    if rawResult == nil {\n")
		if operationKind == OperationProcedure {
			sb.WriteString("      return schema.NewProcedureResult(nil).Encode(), nil")
		} else {
			sb.WriteString("      return nil, nil")
		}
		sb.WriteString("\n    }")
	}
	sb.WriteString(`
    connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
      "raw_result": rawResult,
    })`)
}
