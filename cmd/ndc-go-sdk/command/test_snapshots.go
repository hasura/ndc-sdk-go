package command

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk/command/internal"
	"github.com/hasura/ndc-sdk-go/schema"
)

// GenTestSnapshotArguments represents arguments for test snapshot generation
type GenTestSnapshotArguments struct {
	Schema   string `help:"NDC schema file path. Use either endpoint or schema path"`
	Endpoint string `help:"The endpoint of the connector. Use either endpoint or schema path"`
	Dir      string `help:"The directory of test snapshots."`
	Depth    uint   `help:"The selection depth of nested fields in result types." default:"10"`
	Seed     *int64 `help:"Using a fixed seed will produce the same output on every run."`
}

// genTestSnapshotsCommand
type genTestSnapshotsCommand struct {
	args   *GenTestSnapshotArguments
	schema schema.SchemaResponse
	random *rand.Rand
}

// GenTestSnapshots generates test snapshots from NDC schema
func GenTestSnapshots(args *GenTestSnapshotArguments) error {
	seed := time.Now().UnixNano()
	if args.Seed != nil {
		seed = *args.Seed
	}
	random := rand.New(rand.NewSource(seed))
	cmd := genTestSnapshotsCommand{
		args:   args,
		random: random,
	}

	if err := cmd.fetchSchema(); err != nil {
		return err
	}

	for _, fn := range cmd.schema.Functions {
		if err := cmd.genFunction(&fn); err != nil {
			return err
		}
	}

	for _, proc := range cmd.schema.Procedures {
		if err := cmd.genProcedure(&proc); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *genTestSnapshotsCommand) fetchSchema() error {
	if cmd.args.Schema != "" {
		rawBytes, err := os.ReadFile(cmd.args.Schema)
		if err != nil {
			return fmt.Errorf("failed to read schema from %s: %s", cmd.args.Schema, err)
		}
		if err := json.Unmarshal(rawBytes, &cmd.schema); err != nil {
			return fmt.Errorf("failed to decode schema json from %s: %s", cmd.args.Schema, err)
		}
		return nil
	}

	if cmd.args.Endpoint != "" {
		resp, err := http.Get(cmd.args.Schema)
		if err != nil {
			return fmt.Errorf("failed to fetch schema from %s: %s", cmd.args.Endpoint, err)
		}

		if resp.StatusCode != http.StatusOK {
			var respBytes []byte
			if resp.Body != nil {
				respBytes, _ = io.ReadAll(resp.Body)
			}
			if len(respBytes) == 0 {
				respBytes = []byte(http.StatusText(resp.StatusCode))
			}
			return fmt.Errorf("failed to fetch schema from %s: %s", cmd.args.Endpoint, string(respBytes))
		}
		if resp.Body == nil {
			return fmt.Errorf("received empty response from %s", cmd.args.Endpoint)
		}

		if err := json.NewDecoder(resp.Body).Decode(&cmd.schema); err != nil {
			return fmt.Errorf("failed to decode schema json from %s: %s", cmd.args.Schema, err)
		}
		return nil
	}

	return fmt.Errorf("required either endpoint or file path to the schema")
}

func (cmd *genTestSnapshotsCommand) genFunction(fn *schema.FunctionInfo) error {

	args, err := cmd.genQueryArguments(fn.Arguments)
	if err != nil {
		return fmt.Errorf("failed to generate arguments for %s function: %s", fn.Name, err)
	}
	fields, value, err := cmd.genNestFieldAndValue(fn.ResultType)
	if err != nil {
		return fmt.Errorf("failed to generate result for %s function: %s", fn.Name, err)
	}

	queryReq := schema.QueryRequest{
		Collection: fn.Name,
		Query: schema.Query{
			Fields: schema.QueryFields{
				"__value": schema.NewColumnField("__value", fields).Encode(),
			},
		},
		Arguments:               args,
		CollectionRelationships: schema.QueryRequestCollectionRelationships{},
	}

	queryResp := schema.QueryResponse{
		{
			Rows: []map[string]any{
				{
					"__value": value,
				},
			},
		},
	}

	snapshotDir := path.Join(cmd.args.Dir, "query", queryReq.Collection)
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return err
	}

	if err := internal.WritePrettyFileJSON(path.Join(snapshotDir, "request.json"), queryReq); err != nil {
		return err
	}

	return internal.WritePrettyFileJSON(path.Join(snapshotDir, "expected.json"), queryResp)
}

func (cmd *genTestSnapshotsCommand) genQueryArguments(arguments schema.FunctionInfoArguments) (schema.QueryRequestArguments, error) {
	result := schema.QueryRequestArguments{}
	for key, arg := range arguments {
		_, value, err := cmd.genNestFieldAndValue(arg.Type)
		if err != nil {
			return nil, err
		}
		result[key] = schema.Argument{
			Type:  schema.ArgumentTypeLiteral,
			Value: value,
		}
	}
	return result, nil
}

func (cmd *genTestSnapshotsCommand) genProcedure(proc *schema.ProcedureInfo) error {
	args, err := cmd.genOperationArguments(proc.Arguments)
	if err != nil {
		return fmt.Errorf("failed to generate arguments for %s procedure: %s", proc.Name, err)
	}

	fields, value, err := cmd.genNestFieldAndValue(proc.ResultType)
	if err != nil {
		return fmt.Errorf("failed to generate result for %s procedure: %s", proc.Name, err)
	}
	var rawFields schema.NestedField
	if fields != nil {
		rawFields = fields.Encode()
	}
	mutationReq := schema.MutationRequest{
		Operations: []schema.MutationOperation{
			{
				Type:      schema.MutationOperationProcedure,
				Name:      proc.Name,
				Arguments: args,
				Fields:    rawFields,
			},
		},
		CollectionRelationships: make(schema.MutationRequestCollectionRelationships),
	}

	mutationResp := schema.MutationResponse{
		OperationResults: []schema.MutationOperationResults{
			schema.NewProcedureResult(value).Encode(),
		},
	}

	snapshotDir := path.Join(cmd.args.Dir, "mutation", proc.Name)
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return err
	}

	if err := internal.WritePrettyFileJSON(path.Join(snapshotDir, "request.json"), mutationReq); err != nil {
		return err
	}

	return internal.WritePrettyFileJSON(path.Join(snapshotDir, "expected.json"), mutationResp)
}

func (cmd *genTestSnapshotsCommand) genOperationArguments(arguments schema.ProcedureInfoArguments) ([]byte, error) {
	result := map[string]any{}
	for key, arg := range arguments {
		_, value, err := cmd.genNestFieldAndValue(arg.Type)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}

	return json.Marshal(result)
}

func (cmd *genTestSnapshotsCommand) genNestFieldAndValue(rawType schema.Type) (schema.NestedFieldEncoder, any, error) {
	nestedField, value, _, err := cmd.genNestFieldAndValueInternal(rawType, 0)
	return nestedField, value, err
}

func (cmd *genTestSnapshotsCommand) genNestFieldAndValueInternal(rawType schema.Type, currentDepth uint) (schema.NestedFieldEncoder, any, bool, error) {
	resultType, err := rawType.InterfaceT()

	switch ty := resultType.(type) {
	case *schema.NullableType:
		return cmd.genNestFieldAndValueInternal(ty.UnderlyingType, currentDepth)
	case *schema.ArrayType:
		if currentDepth >= cmd.args.Depth {
			return nil, nil, false, nil
		}
		innerType, data, isScalar, err := cmd.genNestFieldAndValueInternal(ty.ElementType, currentDepth+1)
		if err != nil {
			return nil, nil, false, err
		}
		if isScalar {
			return nil, []any{data}, isScalar, nil
		}
		return schema.NewNestedArray(innerType), []any{data}, isScalar, nil
	case *schema.NamedType:
		if currentDepth >= cmd.args.Depth {
			return nil, nil, false, nil
		}
		if scalar, ok := cmd.schema.ScalarTypes[ty.Name]; ok {
			return nil, internal.GenRandomScalarValue(cmd.random, ty.Name, &scalar), true, nil
		}
		objectType, ok := cmd.schema.ObjectTypes[ty.Name]
		if !ok {
			return nil, nil, false, fmt.Errorf("the named type <%s> does not exist", ty.Name)
		}

		fields := make(map[string]schema.FieldEncoder)
		values := make(map[string]any)
		for key, field := range objectType.Fields {
			innerType, value, _, err := cmd.genNestFieldAndValueInternal(field.Type, currentDepth+1)
			if err != nil {
				return nil, nil, false, err
			}
			fields[key] = schema.NewColumnField(key, innerType)
			values[key] = value
		}
		return schema.NewNestedObject(fields), values, false, nil
	default:
		return nil, nil, false, err
	}
}
