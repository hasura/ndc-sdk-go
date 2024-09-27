package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
)

func TestEvalNestedFields(t *testing.T) {
	testCases := []struct {
		Name      string
		Input     any
		Selection schema.NestedField
		Expected  any
	}{
		{
			Name: "no_field",
			Input: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   "1",
				Name: "John",
			},
			Selection: schema.NewNestedObject(map[string]schema.FieldEncoder{}).Encode(),
			Expected: map[string]any{
				"id":   "1",
				"name": "John",
			},
		},
		{
			Name: "column_field",
			Input: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   "1",
				Name: "John",
			},
			Selection: schema.NewNestedObject(map[string]schema.FieldEncoder{
				"id": schema.NewColumnField("id", nil),
			}).Encode(),
			Expected: map[string]any{
				"id": "1",
			},
		},
		{
			Name: "nested_fields",
			Input: struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Articles []struct {
					ID        int       `json:"id"`
					Name      string    `json:"name"`
					CreatedAt time.Time `json:"created_at"`
				} `json:"articles"`
			}{
				ID:   "1",
				Name: "John",
				Articles: []struct {
					ID        int       `json:"id"`
					Name      string    `json:"name"`
					CreatedAt time.Time `json:"created_at"`
				}{
					{
						ID:        1,
						Name:      "Article 1",
						CreatedAt: time.Date(2020, 01, 01, 00, 0, 0, 0, time.UTC),
					},
				},
			},
			Selection: schema.NewNestedObject(map[string]schema.FieldEncoder{
				"id": schema.NewColumnField("id", nil),
				"articles": schema.NewColumnField("articles", schema.NewNestedArray(schema.NewNestedObject(map[string]schema.FieldEncoder{
					"id":         schema.NewColumnField("id", nil),
					"created_at": schema.NewColumnField("created_at", nil),
				}))),
			}).Encode(),
			Expected: map[string]any{
				"id": "1",
				"articles": []map[string]any{
					{
						"id":         1,
						"created_at": time.Date(2020, 01, 01, 00, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			Name: "rename_fields",
			Input: struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Articles []struct {
					Name string
				} `json:"articles"`
			}{
				ID:   "1",
				Name: "John",
				Articles: []struct {
					Name string
				}{
					{
						Name: "Article 1",
					},
				},
			},
			Selection: schema.NewNestedObject(map[string]schema.FieldEncoder{
				"id":   schema.NewColumnField("id", nil),
				"Name": schema.NewColumnField("name", nil),
				"articles": schema.NewColumnField("articles", schema.NewNestedArray(schema.NewNestedObject(map[string]schema.FieldEncoder{
					"name": schema.NewColumnField("Name", nil),
				}))),
			}).Encode(),
			Expected: map[string]any{
				"id":   "1",
				"Name": "John",
				"articles": []map[string]any{
					{
						"name": "Article 1",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := EvalNestedColumnFields(tc.Selection, tc.Input)
			if err != nil {
				t.Errorf("failed to evaluate nested column fields: %s", err)
				t.FailNow()
			}
			if !internal.DeepEqual(tc.Expected, result) {
				t.Errorf("evaluated result does not equal\nexpected: %+v, got: %+v", tc.Expected, result)
				t.FailNow()
			}
		})
	}
}

func TestMergeSchemas(t *testing.T) {
	testCases := []struct {
		Name     string
		Inputs   []*schema.SchemaResponse
		Expected schema.SchemaResponse
		Errors   []error
	}{
		{
			Name: "simple",
			Inputs: []*schema.SchemaResponse{
				nil,
				{
					Collections: []schema.CollectionInfo{
						{
							Name:      "Foo",
							Arguments: schema.CollectionInfoArguments{},
						},
					},
					ObjectTypes: schema.SchemaResponseObjectTypes{
						"GetArticlesResult": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"Name": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
						"GetTypesArguments": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"ArrayBigInt": schema.ObjectField{
									Type: schema.NewArrayType(schema.NewNamedType("BigInt")).Encode(),
								},
							},
						},
						"GetTypesArgumentsArrayObject": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"content": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
						"GetTypesArgumentsObjectPtr": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"Lat": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
								"Long": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
							},
						},
						"HelloResult": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"error": schema.ObjectField{
									Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
								},
								"foo": schema.ObjectField{
									Type: schema.NewNamedType("Foo").Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("UUID").Encode(),
								},
								"num": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
								"text": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
					},
					Functions: []schema.FunctionInfo{
						{
							Name:        "getBool",
							Description: ToPtr("return an scalar boolean"),
							ResultType:  schema.NewNamedType("Boolean").Encode(),
							Arguments:   map[string]schema.ArgumentInfo{},
						},
						{
							Name:        "hello",
							Description: ToPtr("sends a hello message"),
							ResultType:  schema.NewNullableType(schema.NewNamedType("HelloResult")).Encode(),
							Arguments:   map[string]schema.ArgumentInfo{},
						},
						{
							Name:        "getArticles",
							Description: ToPtr("GetArticles"),
							ResultType:  schema.NewArrayType(schema.NewNamedType("GetArticlesResult")).Encode(),
							Arguments: map[string]schema.ArgumentInfo{
								"Limit": {
									Type: schema.NewNamedType("Float64").Encode(),
								},
							},
						},
					},
					Procedures: []schema.ProcedureInfo{
						{
							Name:        "create_article",
							Description: ToPtr("CreateArticle"),
							ResultType:  schema.NewNullableType(schema.NewNamedType("CreateArticleResult")).Encode(),
							Arguments: map[string]schema.ArgumentInfo{
								"author": {
									Type: schema.NewNamedType("CreateArticleArgumentsAuthor").Encode(),
								},
							},
						},
						{
							Name:        "createAuthor",
							Description: ToPtr("creates an author"),
							ResultType:  schema.NewNullableType(schema.NewNamedType("CreateAuthorResult")).Encode(),
							Arguments: map[string]schema.ArgumentInfo{
								"name": {
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
						{
							Name:        "createAuthors",
							Description: ToPtr("creates a list of authors"),
							ResultType:  schema.NewArrayType(schema.NewNamedType("CreateAuthorResult")).Encode(),
							Arguments: map[string]schema.ArgumentInfo{
								"names": {
									Type: schema.NewArrayType(schema.NewNamedType("String")).Encode(),
								},
							},
						},
					},
					ScalarTypes: schema.SchemaResponseScalarTypes{
						"BigInt": schema.ScalarType{
							AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
							ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
							Representation:      schema.NewTypeRepresentationBigInteger().Encode(),
						},
						"JSON": schema.ScalarType{
							AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
							ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
							Representation:      schema.NewTypeRepresentationJSON().Encode(),
						},
						"TimestampTZ": schema.ScalarType{
							AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
							ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
							Representation:      schema.NewTypeRepresentationTimestampTZ().Encode(),
						},
						"UUID": schema.ScalarType{
							AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
							ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
							Representation:      schema.NewTypeRepresentationUUID().Encode(),
						},
					},
				},
				{
					Collections: []schema.CollectionInfo{
						{
							Name:      "Foo",
							Arguments: schema.CollectionInfoArguments{},
						},
					},
					ObjectTypes: schema.SchemaResponseObjectTypes{
						"Author": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"created_at": schema.ObjectField{
									Type: schema.NewNamedType("TimestampTZ").Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
						"CreateArticleArgumentsAuthor": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"created_at": schema.ObjectField{
									Type: schema.NewNamedType("TimestampTZ").Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("UUID").Encode(),
								},
							},
						},
						"CreateArticleResult": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"authors": schema.ObjectField{
									Type: schema.NewArrayType(schema.NewNamedType("Author")).Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
							},
						},
						"CreateAuthorResult": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"created_at": schema.ObjectField{
									Type: schema.NewNamedType("TimestampTZ").Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
								"name": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
						"GetTypesArgumentsArrayObjectPtr": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"content": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
						"GetTypesArgumentsObject": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"created_at": schema.ObjectField{
									Type: schema.NewNamedType("TimestampTZ").Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("UUID").Encode(),
								},
							},
						},
						"GetTypesArgumentsObjectPtr": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"Lat": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
								"Long": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
							},
						},
						"HelloResult": schema.ObjectType{
							Fields: schema.ObjectTypeFields{
								"error": schema.ObjectField{
									Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
								},
								"foo": schema.ObjectField{
									Type: schema.NewNamedType("Foo").Encode(),
								},
								"id": schema.ObjectField{
									Type: schema.NewNamedType("UUID").Encode(),
								},
								"num": schema.ObjectField{
									Type: schema.NewNamedType("Int32").Encode(),
								},
								"text": schema.ObjectField{
									Type: schema.NewNamedType("String").Encode(),
								},
							},
						},
					},
					Functions: []schema.FunctionInfo{
						{
							Name:        "getBool",
							Description: ToPtr("return an scalar boolean"),
							ResultType:  schema.NewNamedType("Boolean").Encode(),
							Arguments:   map[string]schema.ArgumentInfo{},
						},
						{
							Name:       "getTypes",
							ResultType: schema.NewNullableType(schema.NewNamedType("GetTypesArguments")).Encode(),
							Arguments: map[string]schema.ArgumentInfo{
								"ArrayBigInt": {
									Type: schema.NewArrayType(schema.NewNamedType("BigInt")).Encode(),
								},
							},
						},
					},
					Procedures: []schema.ProcedureInfo{
						{
							Name:        "create_article",
							Description: ToPtr("CreateArticle"),
							ResultType:  schema.NewNullableType(schema.NewNamedType("CreateArticleResult")).Encode(),
							Arguments: map[string]schema.ArgumentInfo{
								"author": {
									Type: schema.NewNamedType("CreateArticleArgumentsAuthor").Encode(),
								},
							},
						},
						{
							Name:        "increase",
							Description: ToPtr("Increase"),
							ResultType:  schema.NewNamedType("Int32").Encode(),
							Arguments:   map[string]schema.ArgumentInfo{},
						},
					},
					ScalarTypes: schema.SchemaResponseScalarTypes{
						"BigInt": schema.ScalarType{
							AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
							ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
							Representation:      schema.NewTypeRepresentationBigInteger().Encode(),
						},
						"Boolean": schema.ScalarType{
							AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
							ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
							Representation:      schema.NewTypeRepresentationBoolean().Encode(),
						},
					},
				},
			},
			Expected: schema.SchemaResponse{
				Collections: []schema.CollectionInfo{
					{
						Name:      "Foo",
						Arguments: schema.CollectionInfoArguments{},
					},
				},
				ObjectTypes: schema.SchemaResponseObjectTypes{
					"Author": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"created_at": schema.ObjectField{
								Type: schema.NewNamedType("TimestampTZ").Encode(),
							},
							"id": schema.ObjectField{
								Type: schema.NewNamedType("String").Encode(),
							},
						},
					},
					"CreateArticleArgumentsAuthor": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"created_at": schema.ObjectField{
								Type: schema.NewNamedType("TimestampTZ").Encode(),
							},
							"id": schema.ObjectField{
								Type: schema.NewNamedType("UUID").Encode(),
							},
						},
					},
					"CreateArticleResult": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"authors": schema.ObjectField{
								Type: schema.NewArrayType(schema.NewNamedType("Author")).Encode(),
							},
							"id": schema.ObjectField{
								Type: schema.NewNamedType("Int32").Encode(),
							},
						},
					},
					"CreateAuthorResult": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"created_at": schema.ObjectField{
								Type: schema.NewNamedType("TimestampTZ").Encode(),
							},
							"id": schema.ObjectField{
								Type: schema.NewNamedType("Int32").Encode(),
							},
							"name": schema.ObjectField{
								Type: schema.NewNamedType("String").Encode(),
							},
						},
					},
					"GetArticlesResult": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"Name": schema.ObjectField{
								Type: schema.NewNamedType("String").Encode(),
							},
							"id": schema.ObjectField{
								Type: schema.NewNamedType("String").Encode(),
							},
						},
					},
					"GetTypesArguments": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"ArrayBigInt": schema.ObjectField{
								Type: schema.NewArrayType(schema.NewNamedType("BigInt")).Encode(),
							},
						},
					},
					"GetTypesArgumentsArrayObject": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"content": schema.ObjectField{
								Type: schema.NewNamedType("String").Encode(),
							},
						},
					},
					"GetTypesArgumentsArrayObjectPtr": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"content": schema.ObjectField{
								Type: schema.NewNamedType("String").Encode(),
							},
						},
					},
					"GetTypesArgumentsObject": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"created_at": schema.ObjectField{
								Type: schema.NewNamedType("TimestampTZ").Encode(),
							},
							"id": schema.ObjectField{
								Type: schema.NewNamedType("UUID").Encode(),
							},
						},
					},
					"GetTypesArgumentsObjectPtr": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"Lat": schema.ObjectField{
								Type: schema.NewNamedType("Int32").Encode(),
							},
							"Long": schema.ObjectField{
								Type: schema.NewNamedType("Int32").Encode(),
							},
						},
					},
					"HelloResult": schema.ObjectType{
						Fields: schema.ObjectTypeFields{
							"error": schema.ObjectField{
								Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
							},
							"foo": schema.ObjectField{
								Type: schema.NewNamedType("Foo").Encode(),
							},
							"id": schema.ObjectField{
								Type: schema.NewNamedType("UUID").Encode(),
							},
							"num": schema.ObjectField{
								Type: schema.NewNamedType("Int32").Encode(),
							},
							"text": schema.ObjectField{
								Type: schema.NewNamedType("String").Encode(),
							},
						},
					},
				},
				Functions: []schema.FunctionInfo{
					{
						Name:        "getArticles",
						Description: ToPtr("GetArticles"),
						ResultType:  schema.NewArrayType(schema.NewNamedType("GetArticlesResult")).Encode(),
						Arguments: map[string]schema.ArgumentInfo{
							"Limit": {
								Type: schema.NewNamedType("Float64").Encode(),
							},
						},
					},
					{
						Name:        "getBool",
						Description: ToPtr("return an scalar boolean"),
						ResultType:  schema.NewNamedType("Boolean").Encode(),
						Arguments:   map[string]schema.ArgumentInfo{},
					},
					{
						Name:       "getTypes",
						ResultType: schema.NewNullableType(schema.NewNamedType("GetTypesArguments")).Encode(),
						Arguments: map[string]schema.ArgumentInfo{
							"ArrayBigInt": {
								Type: schema.NewArrayType(schema.NewNamedType("BigInt")).Encode(),
							},
						},
					},
					{
						Name:        "hello",
						Description: ToPtr("sends a hello message"),
						ResultType:  schema.NewNullableType(schema.NewNamedType("HelloResult")).Encode(),
						Arguments:   map[string]schema.ArgumentInfo{},
					},
				},
				Procedures: []schema.ProcedureInfo{
					{
						Name:        "createAuthor",
						Description: ToPtr("creates an author"),
						ResultType:  schema.NewNullableType(schema.NewNamedType("CreateAuthorResult")).Encode(),
						Arguments: map[string]schema.ArgumentInfo{
							"name": {
								Type: schema.NewNamedType("String").Encode(),
							},
						},
					},
					{
						Name:        "createAuthors",
						Description: ToPtr("creates a list of authors"),
						ResultType:  schema.NewArrayType(schema.NewNamedType("CreateAuthorResult")).Encode(),
						Arguments: map[string]schema.ArgumentInfo{
							"names": {
								Type: schema.NewArrayType(schema.NewNamedType("String")).Encode(),
							},
						},
					},
					{
						Name:        "create_article",
						Description: ToPtr("CreateArticle"),
						ResultType:  schema.NewNullableType(schema.NewNamedType("CreateArticleResult")).Encode(),
						Arguments: map[string]schema.ArgumentInfo{
							"author": {
								Type: schema.NewNamedType("CreateArticleArgumentsAuthor").Encode(),
							},
						},
					},
					{
						Name:        "increase",
						Description: ToPtr("Increase"),
						ResultType:  schema.NewNamedType("Int32").Encode(),
						Arguments:   map[string]schema.ArgumentInfo{},
					},
				},
				ScalarTypes: schema.SchemaResponseScalarTypes{
					"BigInt": schema.ScalarType{
						AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
						ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
						Representation:      schema.NewTypeRepresentationBigInteger().Encode(),
					},
					"Boolean": schema.ScalarType{
						AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
						ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
						Representation:      schema.NewTypeRepresentationBoolean().Encode(),
					},
					"JSON": schema.ScalarType{
						AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
						ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
						Representation:      schema.NewTypeRepresentationJSON().Encode(),
					},
					"TimestampTZ": schema.ScalarType{
						AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
						ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
						Representation:      schema.NewTypeRepresentationTimestampTZ().Encode(),
					},
					"UUID": schema.ScalarType{
						AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
						ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
						Representation:      schema.NewTypeRepresentationUUID().Encode(),
					},
				},
			},
			Errors: []error{
				errors.New("collection type Foo exists"),
				errors.New("function type getBool exists"),
				errors.New("procedure type create_article exists"),
				errors.New("object type HelloResult exists"),
				errors.New("object type GetTypesArgumentsObjectPtr exists"),
				errors.New("scalar type BigInt exists"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, errs := MergeSchemas(tc.Inputs...)
			assertDeepEqual(t, tc.Errors, errs)
			assertDeepEqual(t, tc.Expected.Collections, result.Collections)
			assertDeepEqual(t, tc.Expected.Functions, result.Functions)
			assertDeepEqual(t, tc.Expected.ObjectTypes, result.ObjectTypes)
			assertDeepEqual(t, tc.Expected.Procedures, result.Procedures)
			assertDeepEqual(t, tc.Expected.ScalarTypes, result.ScalarTypes)
		})
	}
}
