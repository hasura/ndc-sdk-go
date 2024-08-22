// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package schema

import "encoding/json"
import "fmt"
import "reflect"

// The definition of an aggregation function on a scalar type
type AggregateFunctionDefinition struct {
	// The scalar or object type of the result of this function
	ResultType Type `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AggregateFunctionDefinition) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["result_type"]; raw != nil && !ok {
		return fmt.Errorf("field result_type in AggregateFunctionDefinition: required")
	}
	type Plain AggregateFunctionDefinition
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AggregateFunctionDefinition(plain)
	return nil
}

type ArgumentInfo struct {
	// Argument description
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the type of this argument
	Type Type `json:"type" yaml:"type" mapstructure:"type"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArgumentInfo) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ArgumentInfo: required")
	}
	type Plain ArgumentInfo
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ArgumentInfo(plain)
	return nil
}

// Describes the features of the specification which a data connector implements.
type Capabilities struct {
	// Mutation corresponds to the JSON schema field "mutation".
	Mutation MutationCapabilities `json:"mutation" yaml:"mutation" mapstructure:"mutation"`

	// Query corresponds to the JSON schema field "query".
	Query QueryCapabilities `json:"query" yaml:"query" mapstructure:"query"`

	// Relationships corresponds to the JSON schema field "relationships".
	Relationships interface{} `json:"relationships,omitempty" yaml:"relationships,omitempty" mapstructure:"relationships,omitempty"`
}

type CapabilitiesResponse struct {
	// Capabilities corresponds to the JSON schema field "capabilities".
	Capabilities Capabilities `json:"capabilities" yaml:"capabilities" mapstructure:"capabilities"`

	// Version corresponds to the JSON schema field "version".
	Version string `json:"version" yaml:"version" mapstructure:"version"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CapabilitiesResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["capabilities"]; raw != nil && !ok {
		return fmt.Errorf("field capabilities in CapabilitiesResponse: required")
	}
	if _, ok := raw["version"]; raw != nil && !ok {
		return fmt.Errorf("field version in CapabilitiesResponse: required")
	}
	type Plain CapabilitiesResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CapabilitiesResponse(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Capabilities) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["mutation"]; raw != nil && !ok {
		return fmt.Errorf("field mutation in Capabilities: required")
	}
	if _, ok := raw["query"]; raw != nil && !ok {
		return fmt.Errorf("field query in Capabilities: required")
	}
	type Plain Capabilities
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Capabilities(plain)
	return nil
}

type CollectionInfo struct {
	// Any arguments that this collection requires
	Arguments CollectionInfoArguments `json:"arguments" yaml:"arguments" mapstructure:"arguments"`

	// Description of the collection
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Any foreign key constraints enforced on this collection
	ForeignKeys CollectionInfoForeignKeys `json:"foreign_keys" yaml:"foreign_keys" mapstructure:"foreign_keys"`

	// The name of the collection
	//
	// Note: these names are abstract - there is no requirement that this name
	// correspond to the name of an actual collection in the database.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The name of the collection's object type
	Type string `json:"type" yaml:"type" mapstructure:"type"`

	// Any uniqueness constraints enforced on this collection
	UniquenessConstraints CollectionInfoUniquenessConstraints `json:"uniqueness_constraints" yaml:"uniqueness_constraints" mapstructure:"uniqueness_constraints"`
}

// Any arguments that this collection requires
type CollectionInfoArguments map[string]ArgumentInfo

// Any foreign key constraints enforced on this collection
type CollectionInfoForeignKeys map[string]ForeignKeyConstraint

// Any uniqueness constraints enforced on this collection
type CollectionInfoUniquenessConstraints map[string]UniquenessConstraint

// UnmarshalJSON implements json.Unmarshaler.
func (j *CollectionInfo) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["arguments"]; raw != nil && !ok {
		return fmt.Errorf("field arguments in CollectionInfo: required")
	}
	if _, ok := raw["foreign_keys"]; raw != nil && !ok {
		return fmt.Errorf("field foreign_keys in CollectionInfo: required")
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in CollectionInfo: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in CollectionInfo: required")
	}
	if _, ok := raw["uniqueness_constraints"]; raw != nil && !ok {
		return fmt.Errorf("field uniqueness_constraints in CollectionInfo: required")
	}
	type Plain CollectionInfo
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CollectionInfo(plain)
	return nil
}

// The definition of a comparison operator on a scalar type

type ErrorResponse struct {
	// Any additional structured information about the error
	Details interface{} `json:"details" yaml:"details" mapstructure:"details"`

	// A human-readable summary of the error
	Message string `json:"message" yaml:"message" mapstructure:"message"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ErrorResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["details"]; raw != nil && !ok {
		return fmt.Errorf("field details in ErrorResponse: required")
	}
	if _, ok := raw["message"]; raw != nil && !ok {
		return fmt.Errorf("field message in ErrorResponse: required")
	}
	type Plain ErrorResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ErrorResponse(plain)
	return nil
}

type ExistsCapabilities struct {
	// Does the connector support ExistsInCollection::NestedCollection
	NestedCollections interface{} `json:"nested_collections,omitempty" yaml:"nested_collections,omitempty" mapstructure:"nested_collections,omitempty"`
}

type ExplainResponse struct {
	// A list of human-readable key-value pairs describing a query execution plan. For
	// example, a connector for a relational database might return the generated SQL
	// and/or the output of the `EXPLAIN` command. An API-based connector might encode
	// a list of statically-known API calls which would be made.
	Details ExplainResponseDetails `json:"details" yaml:"details" mapstructure:"details"`
}

// A list of human-readable key-value pairs describing a query execution plan. For
// example, a connector for a relational database might return the generated SQL
// and/or the output of the `EXPLAIN` command. An API-based connector might encode
// a list of statically-known API calls which would be made.
type ExplainResponseDetails map[string]string

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExplainResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["details"]; raw != nil && !ok {
		return fmt.Errorf("field details in ExplainResponse: required")
	}
	type Plain ExplainResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ExplainResponse(plain)
	return nil
}

type ForeignKeyConstraint struct {
	// The columns on which you want want to define the foreign key.
	ColumnMapping ForeignKeyConstraintColumnMapping `json:"column_mapping" yaml:"column_mapping" mapstructure:"column_mapping"`

	// The name of a collection
	ForeignCollection string `json:"foreign_collection" yaml:"foreign_collection" mapstructure:"foreign_collection"`
}

// The columns on which you want want to define the foreign key.
type ForeignKeyConstraintColumnMapping map[string]string

// UnmarshalJSON implements json.Unmarshaler.
func (j *ForeignKeyConstraint) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["column_mapping"]; raw != nil && !ok {
		return fmt.Errorf("field column_mapping in ForeignKeyConstraint: required")
	}
	if _, ok := raw["foreign_collection"]; raw != nil && !ok {
		return fmt.Errorf("field foreign_collection in ForeignKeyConstraint: required")
	}
	type Plain ForeignKeyConstraint
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ForeignKeyConstraint(plain)
	return nil
}

type FunctionInfo struct {
	// Any arguments that this collection requires
	Arguments FunctionInfoArguments `json:"arguments" yaml:"arguments" mapstructure:"arguments"`

	// Description of the function
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the function
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The name of the function's result type
	ResultType Type `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// Any arguments that this collection requires
type FunctionInfoArguments map[string]ArgumentInfo

// UnmarshalJSON implements json.Unmarshaler.
func (j *FunctionInfo) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["arguments"]; raw != nil && !ok {
		return fmt.Errorf("field arguments in FunctionInfo: required")
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in FunctionInfo: required")
	}
	if _, ok := raw["result_type"]; raw != nil && !ok {
		return fmt.Errorf("field result_type in FunctionInfo: required")
	}
	type Plain FunctionInfo
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = FunctionInfo(plain)
	return nil
}

// A unit value to indicate a particular leaf capability is supported. This is an
// empty struct to allow for future sub-capabilities.
type LeafCapability map[string]interface{}

type MutationCapabilities struct {
	// Does the connector support explaining mutations
	Explain interface{} `json:"explain,omitempty" yaml:"explain,omitempty" mapstructure:"explain,omitempty"`

	// Does the connector support executing multiple mutations in a transaction.
	Transactional interface{} `json:"transactional,omitempty" yaml:"transactional,omitempty" mapstructure:"transactional,omitempty"`
}

type MutationRequest struct {
	// The relationships between collections involved in the entire mutation request
	CollectionRelationships MutationRequestCollectionRelationships `json:"collection_relationships" yaml:"collection_relationships" mapstructure:"collection_relationships"`

	// The mutation operations to perform
	Operations []MutationOperation `json:"operations" yaml:"operations" mapstructure:"operations"`
}

// The relationships between collections involved in the entire mutation request
type MutationRequestCollectionRelationships map[string]Relationship

// UnmarshalJSON implements json.Unmarshaler.
func (j *MutationRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["collection_relationships"]; raw != nil && !ok {
		return fmt.Errorf("field collection_relationships in MutationRequest: required")
	}
	if _, ok := raw["operations"]; raw != nil && !ok {
		return fmt.Errorf("field operations in MutationRequest: required")
	}
	type Plain MutationRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = MutationRequest(plain)
	return nil
}

type MutationResponse struct {
	// The results of each mutation operation, in the same order as they were received
	OperationResults []MutationOperationResults `json:"operation_results" yaml:"operation_results" mapstructure:"operation_results"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MutationResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["operation_results"]; raw != nil && !ok {
		return fmt.Errorf("field operation_results in MutationResponse: required")
	}
	type Plain MutationResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = MutationResponse(plain)
	return nil
}

type NestedFieldCapabilities struct {
	// Does the connector support aggregating values within nested fields
	Aggregates interface{} `json:"aggregates,omitempty" yaml:"aggregates,omitempty" mapstructure:"aggregates,omitempty"`

	// Does the connector support filtering by values of nested fields
	FilterBy interface{} `json:"filter_by,omitempty" yaml:"filter_by,omitempty" mapstructure:"filter_by,omitempty"`

	// Does the connector support ordering by values of nested fields
	OrderBy interface{} `json:"order_by,omitempty" yaml:"order_by,omitempty" mapstructure:"order_by,omitempty"`
}

// The definition of an object field
type ObjectField struct {
	// The arguments available to the field - Matches implementation from
	// CollectionInfo
	Arguments ObjectFieldArguments `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments,omitempty"`

	// Description of this field
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The type of this field
	Type Type `json:"type" yaml:"type" mapstructure:"type"`
}

// The arguments available to the field - Matches implementation from
// CollectionInfo
type ObjectFieldArguments map[string]ArgumentInfo

// UnmarshalJSON implements json.Unmarshaler.
func (j *ObjectField) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ObjectField: required")
	}
	type Plain ObjectField
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ObjectField(plain)
	return nil
}

// The definition of an object type
type ObjectType struct {
	// Description of this type
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Fields defined on this object type
	Fields ObjectTypeFields `json:"fields" yaml:"fields" mapstructure:"fields"`
}

// Fields defined on this object type
type ObjectTypeFields map[string]ObjectField

// UnmarshalJSON implements json.Unmarshaler.
func (j *ObjectType) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["fields"]; raw != nil && !ok {
		return fmt.Errorf("field fields in ObjectType: required")
	}
	type Plain ObjectType
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ObjectType(plain)
	return nil
}

type OrderBy struct {
	// The elements to order by, in priority order
	Elements []OrderByElement `json:"elements" yaml:"elements" mapstructure:"elements"`
}

type OrderByElement struct {
	// OrderDirection corresponds to the JSON schema field "order_direction".
	OrderDirection OrderDirection `json:"order_direction" yaml:"order_direction" mapstructure:"order_direction"`

	// Target corresponds to the JSON schema field "target".
	Target OrderByTarget `json:"target" yaml:"target" mapstructure:"target"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OrderByElement) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["order_direction"]; raw != nil && !ok {
		return fmt.Errorf("field order_direction in OrderByElement: required")
	}
	if _, ok := raw["target"]; raw != nil && !ok {
		return fmt.Errorf("field target in OrderByElement: required")
	}
	type Plain OrderByElement
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OrderByElement(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OrderBy) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["elements"]; raw != nil && !ok {
		return fmt.Errorf("field elements in OrderBy: required")
	}
	type Plain OrderBy
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = OrderBy(plain)
	return nil
}

type OrderDirection string

const OrderDirectionAsc OrderDirection = "asc"
const OrderDirectionDesc OrderDirection = "desc"

var enumValues_OrderDirection = []interface{}{
	"asc",
	"desc",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OrderDirection) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_OrderDirection {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_OrderDirection, v)
	}
	*j = OrderDirection(v)
	return nil
}

type PathElement struct {
	// Values to be provided to any collection arguments
	Arguments PathElementArguments `json:"arguments" yaml:"arguments" mapstructure:"arguments"`

	// A predicate expression to apply to the target collection
	Predicate Expression `json:"predicate,omitempty" yaml:"predicate,omitempty" mapstructure:"predicate,omitempty"`

	// The name of the relationship to follow
	Relationship string `json:"relationship" yaml:"relationship" mapstructure:"relationship"`
}

// Values to be provided to any collection arguments
type PathElementArguments map[string]RelationshipArgument

// UnmarshalJSON implements json.Unmarshaler.
func (j *PathElement) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["arguments"]; raw != nil && !ok {
		return fmt.Errorf("field arguments in PathElement: required")
	}
	if _, ok := raw["relationship"]; raw != nil && !ok {
		return fmt.Errorf("field relationship in PathElement: required")
	}
	type Plain PathElement
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PathElement(plain)
	return nil
}

type ProcedureInfo struct {
	// Any arguments that this collection requires
	Arguments ProcedureInfoArguments `json:"arguments" yaml:"arguments" mapstructure:"arguments"`

	// Column description
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the procedure
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The name of the result type
	ResultType Type `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// Any arguments that this collection requires
type ProcedureInfoArguments map[string]ArgumentInfo

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProcedureInfo) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["arguments"]; raw != nil && !ok {
		return fmt.Errorf("field arguments in ProcedureInfo: required")
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in ProcedureInfo: required")
	}
	if _, ok := raw["result_type"]; raw != nil && !ok {
		return fmt.Errorf("field result_type in ProcedureInfo: required")
	}
	type Plain ProcedureInfo
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ProcedureInfo(plain)
	return nil
}

type Query struct {
	// Aggregate fields of the query
	Aggregates QueryAggregates `json:"aggregates,omitempty" yaml:"aggregates,omitempty" mapstructure:"aggregates,omitempty"`

	// Fields of the query
	Fields QueryFields `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields,omitempty"`

	// Optionally limit to N results
	Limit *int `json:"limit,omitempty" yaml:"limit,omitempty" mapstructure:"limit,omitempty"`

	// Optionally offset from the Nth result
	Offset *int `json:"offset,omitempty" yaml:"offset,omitempty" mapstructure:"offset,omitempty"`

	// OrderBy corresponds to the JSON schema field "order_by".
	OrderBy *OrderBy `json:"order_by,omitempty" yaml:"order_by,omitempty" mapstructure:"order_by,omitempty"`

	// Predicate corresponds to the JSON schema field "predicate".
	Predicate Expression `json:"predicate,omitempty" yaml:"predicate,omitempty" mapstructure:"predicate,omitempty"`
}

// Aggregate fields of the query
type QueryAggregates map[string]Aggregate

type QueryCapabilities struct {
	// Does the connector support aggregate queries
	Aggregates interface{} `json:"aggregates,omitempty" yaml:"aggregates,omitempty" mapstructure:"aggregates,omitempty"`

	// Does the connector support EXISTS predicates
	Exists ExistsCapabilities `json:"exists,omitempty" yaml:"exists,omitempty" mapstructure:"exists,omitempty"`

	// Does the connector support explaining queries
	Explain interface{} `json:"explain,omitempty" yaml:"explain,omitempty" mapstructure:"explain,omitempty"`

	// Does the connector support nested fields
	NestedFields NestedFieldCapabilities `json:"nested_fields,omitempty" yaml:"nested_fields,omitempty" mapstructure:"nested_fields,omitempty"`

	// Does the connector support queries which use variables
	Variables interface{} `json:"variables,omitempty" yaml:"variables,omitempty" mapstructure:"variables,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *QueryCapabilities) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain QueryCapabilities
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["exists"]; !ok || v == nil {
		plain.Exists = ExistsCapabilities{}
	}
	if v, ok := raw["nested_fields"]; !ok || v == nil {
		plain.NestedFields = NestedFieldCapabilities{}
	}
	*j = QueryCapabilities(plain)
	return nil
}

// Fields of the query
type QueryFields map[string]Field

// This is the request body of the query POST endpoint
type QueryRequest struct {
	// Values to be provided to any collection arguments
	Arguments QueryRequestArguments `json:"arguments" yaml:"arguments" mapstructure:"arguments"`

	// The name of a collection
	Collection string `json:"collection" yaml:"collection" mapstructure:"collection"`

	// Any relationships between collections involved in the query request
	CollectionRelationships QueryRequestCollectionRelationships `json:"collection_relationships" yaml:"collection_relationships" mapstructure:"collection_relationships"`

	// The query syntax tree
	Query Query `json:"query" yaml:"query" mapstructure:"query"`

	// One set of named variables for each rowset to fetch. Each variable set should
	// be subtituted in turn, and a fresh set of rows returned.
	Variables []QueryRequestVariablesElem `json:"variables,omitempty" yaml:"variables,omitempty" mapstructure:"variables,omitempty"`
}

// Values to be provided to any collection arguments
type QueryRequestArguments map[string]Argument

// Any relationships between collections involved in the query request
type QueryRequestCollectionRelationships map[string]Relationship

type QueryRequestVariablesElem map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *QueryRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["arguments"]; raw != nil && !ok {
		return fmt.Errorf("field arguments in QueryRequest: required")
	}
	if _, ok := raw["collection"]; raw != nil && !ok {
		return fmt.Errorf("field collection in QueryRequest: required")
	}
	if _, ok := raw["collection_relationships"]; raw != nil && !ok {
		return fmt.Errorf("field collection_relationships in QueryRequest: required")
	}
	if _, ok := raw["query"]; raw != nil && !ok {
		return fmt.Errorf("field query in QueryRequest: required")
	}
	type Plain QueryRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = QueryRequest(plain)
	return nil
}

// Query responses may return multiple RowSets when using queries with variables.
// Else, there should always be exactly one RowSet
type QueryResponse []RowSet

type Relationship struct {
	// Values to be provided to any collection arguments
	Arguments RelationshipArguments `json:"arguments" yaml:"arguments" mapstructure:"arguments"`

	// A mapping between columns on the source collection to columns on the target
	// collection
	ColumnMapping RelationshipColumnMapping `json:"column_mapping" yaml:"column_mapping" mapstructure:"column_mapping"`

	// RelationshipType corresponds to the JSON schema field "relationship_type".
	RelationshipType RelationshipType `json:"relationship_type" yaml:"relationship_type" mapstructure:"relationship_type"`

	// The name of a collection
	TargetCollection string `json:"target_collection" yaml:"target_collection" mapstructure:"target_collection"`
}

// Values to be provided to any collection arguments
type RelationshipArguments map[string]RelationshipArgument

type RelationshipCapabilities struct {
	// Does the connector support ordering by an aggregated array relationship?
	OrderByAggregate interface{} `json:"order_by_aggregate,omitempty" yaml:"order_by_aggregate,omitempty" mapstructure:"order_by_aggregate,omitempty"`

	// Does the connector support comparisons that involve related collections (ie.
	// joins)?
	RelationComparisons interface{} `json:"relation_comparisons,omitempty" yaml:"relation_comparisons,omitempty" mapstructure:"relation_comparisons,omitempty"`
}

// A mapping between columns on the source collection to columns on the target
// collection
type RelationshipColumnMapping map[string]string

type RelationshipType string

const RelationshipTypeArray RelationshipType = "array"
const RelationshipTypeObject RelationshipType = "object"

var enumValues_RelationshipType = []interface{}{
	"object",
	"array",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationshipType) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_RelationshipType {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_RelationshipType, v)
	}
	*j = RelationshipType(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Relationship) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["arguments"]; raw != nil && !ok {
		return fmt.Errorf("field arguments in Relationship: required")
	}
	if _, ok := raw["column_mapping"]; raw != nil && !ok {
		return fmt.Errorf("field column_mapping in Relationship: required")
	}
	if _, ok := raw["relationship_type"]; raw != nil && !ok {
		return fmt.Errorf("field relationship_type in Relationship: required")
	}
	if _, ok := raw["target_collection"]; raw != nil && !ok {
		return fmt.Errorf("field target_collection in Relationship: required")
	}
	type Plain Relationship
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Relationship(plain)
	return nil
}

type RowFieldValue interface{}

type RowSet struct {
	// The results of the aggregates returned by the query
	Aggregates RowSetAggregates `json:"aggregates,omitempty" yaml:"aggregates,omitempty" mapstructure:"aggregates,omitempty"`

	// The rows returned by the query, corresponding to the query's fields
	Rows []map[string]any `json:"rows,omitempty" yaml:"rows,omitempty" mapstructure:"rows,omitempty"`
}

// The results of the aggregates returned by the query
type RowSetAggregates map[string]interface{}

// The definition of a scalar type, i.e. types that can be used as the types of
// columns.
type ScalarType struct {
	// A map from aggregate function names to their definitions. Result type names
	// must be defined scalar types declared in ScalarTypesCapabilities.
	AggregateFunctions ScalarTypeAggregateFunctions `json:"aggregate_functions" yaml:"aggregate_functions" mapstructure:"aggregate_functions"`

	// A map from comparison operator names to their definitions. Argument type names
	// must be defined scalar types declared in ScalarTypesCapabilities.
	ComparisonOperators map[string]ComparisonOperatorDefinition `json:"comparison_operators" yaml:"comparison_operators" mapstructure:"comparison_operators"`

	// A description of valid values for this scalar type. Defaults to
	// `TypeRepresentation::JSON` if omitted
	Representation TypeRepresentation `json:"representation,omitempty" yaml:"representation,omitempty" mapstructure:"representation,omitempty"`
}

// A map from aggregate function names to their definitions. Result type names must
// be defined scalar types declared in ScalarTypesCapabilities.
type ScalarTypeAggregateFunctions map[string]AggregateFunctionDefinition

// A map from comparison operator names to their definitions. Argument type names
// must be defined scalar types declared in ScalarTypesCapabilities.

// UnmarshalJSON implements json.Unmarshaler.
func (j *ScalarType) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["aggregate_functions"]; raw != nil && !ok {
		return fmt.Errorf("field aggregate_functions in ScalarType: required")
	}
	if _, ok := raw["comparison_operators"]; raw != nil && !ok {
		return fmt.Errorf("field comparison_operators in ScalarType: required")
	}
	type Plain ScalarType
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ScalarType(plain)
	return nil
}

type SchemaGeneratedJson struct {
	// CapabilitiesResponse corresponds to the JSON schema field
	// "capabilities_response".
	CapabilitiesResponse CapabilitiesResponse `json:"capabilities_response" yaml:"capabilities_response" mapstructure:"capabilities_response"`

	// ErrorResponse corresponds to the JSON schema field "error_response".
	ErrorResponse ErrorResponse `json:"error_response" yaml:"error_response" mapstructure:"error_response"`

	// ExplainResponse corresponds to the JSON schema field "explain_response".
	ExplainResponse ExplainResponse `json:"explain_response" yaml:"explain_response" mapstructure:"explain_response"`

	// MutationRequest corresponds to the JSON schema field "mutation_request".
	MutationRequest MutationRequest `json:"mutation_request" yaml:"mutation_request" mapstructure:"mutation_request"`

	// MutationResponse corresponds to the JSON schema field "mutation_response".
	MutationResponse MutationResponse `json:"mutation_response" yaml:"mutation_response" mapstructure:"mutation_response"`

	// QueryRequest corresponds to the JSON schema field "query_request".
	QueryRequest QueryRequest `json:"query_request" yaml:"query_request" mapstructure:"query_request"`

	// QueryResponse corresponds to the JSON schema field "query_response".
	QueryResponse QueryResponse `json:"query_response" yaml:"query_response" mapstructure:"query_response"`

	// SchemaResponse corresponds to the JSON schema field "schema_response".
	SchemaResponse SchemaResponse `json:"schema_response" yaml:"schema_response" mapstructure:"schema_response"`

	// ValidateResponse corresponds to the JSON schema field "validate_response".
	ValidateResponse ValidateResponse `json:"validate_response" yaml:"validate_response" mapstructure:"validate_response"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SchemaGeneratedJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["capabilities_response"]; raw != nil && !ok {
		return fmt.Errorf("field capabilities_response in SchemaGeneratedJson: required")
	}
	if _, ok := raw["error_response"]; raw != nil && !ok {
		return fmt.Errorf("field error_response in SchemaGeneratedJson: required")
	}
	if _, ok := raw["explain_response"]; raw != nil && !ok {
		return fmt.Errorf("field explain_response in SchemaGeneratedJson: required")
	}
	if _, ok := raw["mutation_request"]; raw != nil && !ok {
		return fmt.Errorf("field mutation_request in SchemaGeneratedJson: required")
	}
	if _, ok := raw["mutation_response"]; raw != nil && !ok {
		return fmt.Errorf("field mutation_response in SchemaGeneratedJson: required")
	}
	if _, ok := raw["query_request"]; raw != nil && !ok {
		return fmt.Errorf("field query_request in SchemaGeneratedJson: required")
	}
	if _, ok := raw["query_response"]; raw != nil && !ok {
		return fmt.Errorf("field query_response in SchemaGeneratedJson: required")
	}
	if _, ok := raw["schema_response"]; raw != nil && !ok {
		return fmt.Errorf("field schema_response in SchemaGeneratedJson: required")
	}
	if _, ok := raw["validate_response"]; raw != nil && !ok {
		return fmt.Errorf("field validate_response in SchemaGeneratedJson: required")
	}
	type Plain SchemaGeneratedJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SchemaGeneratedJson(plain)
	return nil
}

type SchemaResponse struct {
	// Collections which are available for queries
	Collections []CollectionInfo `json:"collections" yaml:"collections" mapstructure:"collections"`

	// Functions (i.e. collections which return a single column and row)
	Functions []FunctionInfo `json:"functions" yaml:"functions" mapstructure:"functions"`

	// A list of object types which can be used as the types of arguments, or return
	// types of procedures. Names should not overlap with scalar type names.
	ObjectTypes SchemaResponseObjectTypes `json:"object_types" yaml:"object_types" mapstructure:"object_types"`

	// Procedures which are available for execution as part of mutations
	Procedures []ProcedureInfo `json:"procedures" yaml:"procedures" mapstructure:"procedures"`

	// A list of scalar types which will be used as the types of collection columns
	ScalarTypes SchemaResponseScalarTypes `json:"scalar_types" yaml:"scalar_types" mapstructure:"scalar_types"`
}

// A list of object types which can be used as the types of arguments, or return
// types of procedures. Names should not overlap with scalar type names.
type SchemaResponseObjectTypes map[string]ObjectType

// A list of scalar types which will be used as the types of collection columns
type SchemaResponseScalarTypes map[string]ScalarType

// UnmarshalJSON implements json.Unmarshaler.
func (j *SchemaResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["collections"]; raw != nil && !ok {
		return fmt.Errorf("field collections in SchemaResponse: required")
	}
	if _, ok := raw["functions"]; raw != nil && !ok {
		return fmt.Errorf("field functions in SchemaResponse: required")
	}
	if _, ok := raw["object_types"]; raw != nil && !ok {
		return fmt.Errorf("field object_types in SchemaResponse: required")
	}
	if _, ok := raw["procedures"]; raw != nil && !ok {
		return fmt.Errorf("field procedures in SchemaResponse: required")
	}
	if _, ok := raw["scalar_types"]; raw != nil && !ok {
		return fmt.Errorf("field scalar_types in SchemaResponse: required")
	}
	type Plain SchemaResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SchemaResponse(plain)
	return nil
}

// Types track the valid representations of values as JSON

// Representations of scalar types

type UnaryComparisonOperator string

const UnaryComparisonOperatorIsNull UnaryComparisonOperator = "is_null"

var enumValues_UnaryComparisonOperator = []interface{}{
	"is_null",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *UnaryComparisonOperator) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_UnaryComparisonOperator {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_UnaryComparisonOperator, v)
	}
	*j = UnaryComparisonOperator(v)
	return nil
}

type UniquenessConstraint struct {
	// A list of columns which this constraint requires to be unique
	UniqueColumns []string `json:"unique_columns" yaml:"unique_columns" mapstructure:"unique_columns"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *UniquenessConstraint) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["unique_columns"]; raw != nil && !ok {
		return fmt.Errorf("field unique_columns in UniquenessConstraint: required")
	}
	type Plain UniquenessConstraint
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = UniquenessConstraint(plain)
	return nil
}

type ValidateResponse struct {
	// Capabilities corresponds to the JSON schema field "capabilities".
	Capabilities CapabilitiesResponse `json:"capabilities" yaml:"capabilities" mapstructure:"capabilities"`

	// ResolvedConfiguration corresponds to the JSON schema field
	// "resolved_configuration".
	ResolvedConfiguration string `json:"resolved_configuration" yaml:"resolved_configuration" mapstructure:"resolved_configuration"`

	// Schema corresponds to the JSON schema field "schema".
	Schema SchemaResponse `json:"schema" yaml:"schema" mapstructure:"schema"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ValidateResponse) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["capabilities"]; raw != nil && !ok {
		return fmt.Errorf("field capabilities in ValidateResponse: required")
	}
	if _, ok := raw["resolved_configuration"]; raw != nil && !ok {
		return fmt.Errorf("field resolved_configuration in ValidateResponse: required")
	}
	if _, ok := raw["schema"]; raw != nil && !ok {
		return fmt.Errorf("field schema in ValidateResponse: required")
	}
	type Plain ValidateResponse
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ValidateResponse(plain)
	return nil
}
