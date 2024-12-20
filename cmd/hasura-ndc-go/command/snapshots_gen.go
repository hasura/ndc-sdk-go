package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/command/internal"
	"github.com/hasura/ndc-sdk-go/schema"
)

// GenTestSnapshotArguments represents arguments for test snapshot generation.
type GenTestSnapshotArguments struct {
	Schema        string                     `help:"NDC schema file path. Use either endpoint or schema path"`
	Endpoint      string                     `help:"The endpoint of the connector. Use either endpoint or schema path"`
	Dir           string                     `help:"The directory of test snapshots."`
	Depth         uint                       `help:"The selection depth of nested fields in result types." default:"10"`
	Query         []string                   `help:"Specify individual queries to be generated. Separated by commas, or 'all' for all queries"`
	Mutation      []string                   `help:"Specify individual mutations to be generated. Separated by commas, or 'all' for all mutations"`
	FetchResponse bool                       `help:"Fetch snapshot responses from the connector server"`
	Strategy      internal.WriteFileStrategy `help:"Decide the strategy to do when the snapshot file exists. Accept: none, override, update" enum:"none,override,update" default:"none"`
}

// genTestSnapshotsCommand.
type genTestSnapshotsCommand struct {
	args     *GenTestSnapshotArguments
	schema   schema.SchemaResponse
	random   *rand.Rand
	endpoint string
}

// GenTestSnapshots generates test snapshots from NDC schema.
func GenTestSnapshots(args *GenTestSnapshotArguments) error {
	if args.FetchResponse && args.Endpoint == "" {
		return errors.New("require --endpoint argument if the --server-response argument is enabled")
	}

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))
	cmd := genTestSnapshotsCommand{
		args:     args,
		random:   random,
		endpoint: strings.TrimRight(args.Endpoint, "/"),
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
			return fmt.Errorf("failed to read schema from %s: %w", cmd.args.Schema, err)
		}
		if err := json.Unmarshal(rawBytes, &cmd.schema); err != nil {
			return fmt.Errorf("failed to decode schema json from %s: %w", cmd.args.Schema, err)
		}
		return nil
	}

	if cmd.args.Endpoint != "" {
		resp, err := http.Get(cmd.endpoint + "/schema")
		if err != nil {
			return fmt.Errorf("failed to fetch schema from %s: %w", cmd.args.Endpoint, err)
		}
		defer resp.Body.Close()

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
			return fmt.Errorf("failed to decode schema json from %s: %w", cmd.args.Schema, err)
		}
		return nil
	}

	return errors.New("required either endpoint or file path to the schema")
}

func (cmd *genTestSnapshotsCommand) prepareFilePaths(snapshotDir string) (string, string, error) {
	if err := os.MkdirAll(snapshotDir, 0o755); err != nil {
		return "", "", err
	}

	requestFilePath := filepath.Join(snapshotDir, "request.json")
	expectedFilePath := filepath.Join(snapshotDir, "expected.json")

	return requestFilePath, expectedFilePath, nil
}

func (cmd *genTestSnapshotsCommand) genFunction(fn *schema.FunctionInfo) error {
	if !cmd.hasQuery(fn.Name) {
		return nil
	}

	snapshotDir := filepath.Join(cmd.args.Dir, "query", fn.Name)
	requestFilePath, expectedFilePath, err := cmd.prepareFilePaths(snapshotDir)
	if err != nil {
		return err
	}

	currentRequest := readSnapshotFile[schema.QueryRequest](requestFilePath)
	currentResponse := readSnapshotFile[schema.QueryResponse](expectedFilePath)

	if currentRequest != nil && currentResponse != nil && cmd.args.Strategy == internal.WriteFileStrategyNone {
		return nil
	}

	fields, value, err := cmd.genNestFieldAndValue(fn.ResultType)
	if err != nil {
		return fmt.Errorf("failed to generate result for %s function: %w", fn.Name, err)
	}

	queryReq := currentRequest
	if currentRequest == nil || cmd.args.Strategy == internal.WriteFileStrategyOverride {
		args, err := cmd.genQueryArguments(fn.Arguments)
		if err != nil {
			return fmt.Errorf("failed to generate arguments for %s function: %w", fn.Name, err)
		}

		queryReq = &schema.QueryRequest{
			Collection: fn.Name,
			Query: schema.Query{
				Fields: schema.QueryFields{
					"__value": schema.NewColumnField("__value", fields).Encode(),
				},
			},
			Arguments:               args,
			CollectionRelationships: schema.QueryRequestCollectionRelationships{},
		}
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

	requestBytes, err := json.MarshalIndent(queryReq, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: failed to encode the request snapshot: %w", fn.Name, err)
	}

	if cmd.args.FetchResponse {
		httpRequest, err := http.NewRequest(http.MethodPost, cmd.endpoint+"/query", bytes.NewBuffer(requestBytes))
		if err != nil {
			return fmt.Errorf("%s: failed to create http request: %w", fn.Name, err)
		}

		httpRequest.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(httpRequest)
		if err != nil {
			slog.Error("failed to execute http request: "+err.Error(), slog.String("function", fn.Name))
		}

		defer resp.Body.Close()
		commonAttrs := []any{slog.String("function", fn.Name), slog.Int("status", resp.StatusCode)}
		if resp.StatusCode != http.StatusOK {
			errorBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("failed to read http response: "+err.Error(), commonAttrs...)
			} else {
				slog.Error("received non-200 error: "+string(errorBytes), commonAttrs...)
			}
		} else {
			var response schema.QueryResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				slog.Error("failed to read http response: "+err.Error(), commonAttrs...)
			} else {
				queryResp = response
			}
		}
	}

	if err := os.WriteFile(requestFilePath, requestBytes, 0o644); err != nil {
		return fmt.Errorf("%s: failed write request file: %w", fn.Name, err)
	}

	return internal.WritePrettyFileJSON(expectedFilePath, queryResp)
}

func (cmd *genTestSnapshotsCommand) genQueryArguments(arguments schema.FunctionInfoArguments) (schema.QueryRequestArguments, error) {
	result := schema.QueryRequestArguments{}
	for key, arg := range arguments {
		_, value, err := cmd.genNestFieldAndValue(arg.Type)
		if err != nil {
			return nil, err
		}
		result[key] = schema.NewArgumentLiteral(value).Encode()
	}
	return result, nil
}

func (cmd *genTestSnapshotsCommand) genProcedure(proc *schema.ProcedureInfo) error {
	if !cmd.hasMutation(proc.Name) {
		return nil
	}

	snapshotDir := filepath.Join(cmd.args.Dir, "mutation", proc.Name)
	requestFilePath, expectedFilePath, err := cmd.prepareFilePaths(snapshotDir)
	if err != nil {
		return err
	}

	currentRequest := readSnapshotFile[schema.MutationRequest](requestFilePath)
	currentResponse := readSnapshotFile[schema.MutationResponse](expectedFilePath)

	if currentRequest != nil && currentResponse != nil && cmd.args.Strategy == internal.WriteFileStrategyNone {
		return nil
	}

	mutationReq := currentRequest
	fields, value, err := cmd.genNestFieldAndValue(proc.ResultType)
	if err != nil {
		return fmt.Errorf("failed to generate result for %s procedure: %w", proc.Name, err)
	}

	if currentRequest == nil || cmd.args.Strategy == internal.WriteFileStrategyOverride {
		args, err := cmd.genOperationArguments(proc.Arguments)
		if err != nil {
			return fmt.Errorf("failed to generate arguments for %s procedure: %w", proc.Name, err)
		}

		var rawFields schema.NestedField
		if fields != nil {
			rawFields = fields.Encode()
		}
		mutationReq = &schema.MutationRequest{
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
	}

	mutationResp := schema.MutationResponse{
		OperationResults: []schema.MutationOperationResults{
			schema.NewProcedureResult(value).Encode(),
		},
	}

	requestBytes, err := json.MarshalIndent(mutationReq, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: failed to encode the request snapshot: %w", proc.Name, err)
	}

	if cmd.args.FetchResponse {
		httpRequest, err := http.NewRequest(http.MethodPost, cmd.endpoint+"/mutation", bytes.NewBuffer(requestBytes))
		if err != nil {
			return fmt.Errorf("%s: failed to create http request: %w", proc.Name, err)
		}

		httpRequest.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(httpRequest)
		if err != nil {
			slog.Error("failed to execute http request: "+err.Error(), slog.String("function", proc.Name))
		}

		defer resp.Body.Close()
		commonAttrs := []any{slog.String("procedure", proc.Name), slog.Int("status", resp.StatusCode)}
		if resp.StatusCode != http.StatusOK {
			errorBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("failed to read http response: "+err.Error(), commonAttrs...)
			} else {
				slog.Error("received non-200 error: "+string(errorBytes), commonAttrs...)
			}
		} else {
			var response schema.MutationResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				slog.Error("failed to read http response: "+err.Error(), commonAttrs...)
			} else {
				mutationResp = response
			}
		}
	}

	if err := os.WriteFile(requestFilePath, requestBytes, 0o644); err != nil {
		return fmt.Errorf("%s: failed write request file: %w", proc.Name, err)
	}

	return internal.WritePrettyFileJSON(expectedFilePath, mutationResp)
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
		if innerType == nil || isScalar {
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
		if len(fields) == 0 {
			return nil, values, false, nil
		}
		return schema.NewNestedObject(fields), values, false, nil
	default:
		return nil, nil, false, err
	}
}

func (cmd genTestSnapshotsCommand) hasQuery(name string) bool {
	if (len(cmd.args.Query) == 0 && len(cmd.args.Mutation) == 0) || slices.Contains(cmd.args.Query, "all") {
		return true
	}
	return slices.Contains(cmd.args.Query, name)
}

func (cmd genTestSnapshotsCommand) hasMutation(name string) bool {
	if (len(cmd.args.Query) == 0 && len(cmd.args.Mutation) == 0) || slices.Contains(cmd.args.Mutation, "all") {
		return true
	}
	return slices.Contains(cmd.args.Mutation, name)
}

func readSnapshotFile[T any](snapshotPath string) *T {
	rawBytes, err := os.ReadFile(snapshotPath)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error(err.Error(), slog.String("path", snapshotPath))
		}

		return nil
	}

	var result T
	if err := json.Unmarshal(rawBytes, &result); err != nil {
		slog.Error(err.Error(), slog.String("path", snapshotPath))

		return nil
	}

	return &result
}
