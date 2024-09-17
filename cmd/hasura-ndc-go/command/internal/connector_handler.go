package internal

import (
	"fmt"
	"strings"
)

const (
	functionEnumsName = "enumValues_FunctionName"
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
	bs.imports[chb.RawSchema.TypesPackagePath] = ""

	_, _ = bs.builder.WriteString(`
// DataConnectorHandler implements the data connector handler 
type DataConnectorHandler struct{}
`)
	chb.renderQuery(bs.builder)
	chb.renderMutation(bs.builder)

	bs.builder.WriteString(`		
func connector_addSpanEvent(span trace.Span, logger *slog.Logger, name string, data map[string]any, options ...trace.EventOption) {
	logger.Debug(name, slog.Any("data", data))
	attrs := utils.DebugJSONAttributes(data, utils.IsDebug(logger))
	span.AddEvent(name, append(options, trace.WithAttributes(attrs...))...)
}`)
}

func (chb connectorHandlerBuilder) renderOperationNameEnums(sb *strings.Builder, name string, values []string) {
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

func (chb connectorHandlerBuilder) renderQuery(sb *strings.Builder) {
	if len(chb.Functions) == 0 {
		return
	}

	_, _ = sb.WriteString(`
func (dch DataConnectorHandler) Query(ctx context.Context, state *types.State, request *schema.QueryRequest, rawArgs map[string]any) (*schema.RowSet, error) {
	if !slices.Contains(`)
	_, _ = sb.WriteString(functionEnumsName)
	_, _ = sb.WriteString(`, request.Collection) {
		return nil, utils.ErrHandlerNotfound
	}
	queryFields, err := utils.EvalFunctionSelectionFieldValue(request)
	if err != nil {
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}

	result, err := dch.execQuery(ctx, state, request, queryFields, rawArgs)
	if err != nil {
		return nil, schema.UnprocessableContentError(err.Error(), nil)
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
	
func (dch DataConnectorHandler) execQuery(ctx context.Context, state *types.State, request *schema.QueryRequest, queryFields schema.NestedField, rawArgs map[string]any) (any, error) {
	span := trace.SpanFromContext(ctx)
	logger := connector.GetLogger(ctx)
	switch request.Collection {`)

	functionKeys := make([]string, len(chb.Functions))
	for i, fn := range chb.Functions {
		functionKeys[i] = fn.Name
		_, _ = sb.WriteString("\n  case \"")
		_, _ = sb.WriteString(fn.Name)
		_, _ = sb.WriteString("\":\n")

		if fn.ResultType.IsScalar {
			sb.WriteString(`
		if len(queryFields) > 0 {
			return nil, schema.UnprocessableContentError("cannot evaluate selection fields for scalar", nil)
		}`)
		} else if fn.ResultType.IsArray() {
			sb.WriteString(`
		selection, err := queryFields.AsArray()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be array", map[string]any{
				"cause": err.Error(),
			})
		}`)
		} else {
			sb.WriteString(`
		selection, err := queryFields.AsObject()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}`)
		}

		var argumentParamStr string
		if fn.ArgumentsType != nil {
			argName := fn.ArgumentsType.Name
			if fn.ArgumentsType.PackagePath != "" && fn.ArgumentsType.PackagePath != chb.Builder.packagePath {
				chb.Builder.imports[fn.ArgumentsType.PackagePath] = ""
				argName = fmt.Sprintf("%s.%s", fn.ArgumentsType.PackageName, fn.ArgumentsType.Name)
			}

			argumentStr := fmt.Sprintf(`
		var args %s
		if err = args.FromValue(rawArgs); err != nil {
			return nil, schema.UnprocessableContentError("failed to resolve arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		
		connector_addSpanEvent(span, logger, "execute_function", map[string]any{
			"arguments": args,
		})`, argName)
			sb.WriteString(argumentStr)
			argumentParamStr = ", &args"
		}

		if fn.ResultType.IsScalar {
			sb.WriteString(fmt.Sprintf("\n    return %s(ctx, state%s)\n", fn.OriginName, argumentParamStr))
			continue
		}

		sb.WriteString(fmt.Sprintf("\n    rawResult, err := %s(ctx, state%s)", fn.OriginName, argumentParamStr))
		chb.genGeneralOperationResult(sb, fn.ResultType)

		sb.WriteString(`
		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})`)
		if fn.ResultType.IsArray() {
			sb.WriteString("\n    result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)")
		} else {
			sb.WriteString("\n    result, err := utils.EvalNestedColumnObject(selection, rawResult)")
		}
		sb.WriteString(textBlockErrorCheck2)
		sb.WriteString("    return result, nil\n")
	}

	_, _ = sb.WriteString(`
	default:
		return nil, utils.ErrHandlerNotfound
	}
}
`)
	chb.renderOperationNameEnums(sb, functionEnumsName, functionKeys)
}

func (chb connectorHandlerBuilder) renderMutation(sb *strings.Builder) {
	if len(chb.Procedures) == 0 {
		return
	}
	chb.Builder.imports["encoding/json"] = ""

	_, _ = sb.WriteString(`
func (dch DataConnectorHandler) Mutate(ctx context.Context, state *types.State, operation *schema.MutationOperation) (schema.MutationOperationResults, error) {
	span := trace.SpanFromContext(ctx)	
	logger := connector.GetLogger(ctx)
	connector_addSpanEvent(span, logger, "validate_request", map[string]any{
		"operations_name": operation.Name,
	})
	
	switch operation.Name {`)

	for _, fn := range chb.Procedures {
		_, _ = sb.WriteString("\n  case \"")
		_, _ = sb.WriteString(fn.Name)
		_, _ = sb.WriteString("\":\n")

		if fn.ResultType.IsScalar {
			sb.WriteString(`
    if len(operation.Fields) > 0 {
      return nil, schema.UnprocessableContentError("cannot evaluate selection fields for scalar", nil)
    }`)
		} else if fn.ResultType.IsArray() {
			sb.WriteString(`
    selection, err := operation.Fields.AsArray()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be array", map[string]any{
        "cause": err.Error(),
      })
    }`)
		} else {
			sb.WriteString(`
    selection, err := operation.Fields.AsObject()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
        "cause": err.Error(),
      })
    }`)
		}

		var argumentParamStr string
		if fn.ArgumentsType != nil {
			argName := fn.ArgumentsType.Name
			if fn.ArgumentsType.PackagePath != "" && fn.ArgumentsType.PackagePath != chb.Builder.packagePath {
				chb.Builder.imports[fn.ArgumentsType.PackagePath] = ""
				argName = fmt.Sprintf("%s.%s", fn.ArgumentsType.PackageName, fn.ArgumentsType.Name)
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
		if fn.ResultType.IsScalar {
			sb.WriteString(fmt.Sprintf(`
    var err error
    result, err := %s(ctx, state%s)`, fn.OriginName, argumentParamStr))
		} else {
			sb.WriteString(fmt.Sprintf("\n    rawResult, err := %s(ctx, state%s)\n", fn.OriginName, argumentParamStr))
			chb.genGeneralOperationResult(sb, fn.ResultType)

			sb.WriteString(`    connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})`)
			if fn.ResultType.IsArray() {
				sb.WriteString("\n    result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)\n")
			} else {
				sb.WriteString("\n    result, err := utils.EvalNestedColumnObject(selection, rawResult)\n")
			}
		}

		sb.WriteString(textBlockErrorCheck2)
		sb.WriteString("    return schema.NewProcedureResult(result).Encode(), nil\n")
	}

	_, _ = sb.WriteString(`
	default:
		return nil, utils.ErrHandlerNotfound
	}
}
`)
}

func (chb connectorHandlerBuilder) genGeneralOperationResult(sb *strings.Builder, resultType *TypeInfo) {
	sb.WriteString(textBlockErrorCheck2)
	if resultType.IsNullable() {
		sb.WriteString(`
    if rawResult == nil {
      return nil, nil
    }
`)
	} else {
		sb.WriteString(`
    if rawResult == nil {
      return nil, schema.UnprocessableContentError("expected not null result", nil)
    }
`)
	}
}
