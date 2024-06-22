package utils

import (
	"testing"

	"github.com/hasura/ndc-sdk-go/schema"
)

type mockAuthor struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Address   struct {
		Street string `json:"street"`
	} `json:"address"`
}

func (ma mockAuthor) ToMap() (map[string]any, error) {
	result := map[string]any{}
	result["id"] = ma.ID
	result["first_name"] = ma.FirstName
	result["last_name"] = ma.LastName
	result["address"] = map[string]any{
		"street": ma.Address.Street,
	}
	return result, nil
}

type mockArticle struct {
	ID      int          `json:"id"`
	Authors []mockAuthor `json:"authors"`
}

func (ma mockArticle) ToMap() (map[string]any, error) {
	result := map[string]any{}
	result["id"] = ma.ID
	authors := make([]map[string]any, len(ma.Authors))
	for i, author := range ma.Authors {
		au, err := author.ToMap()
		if err != nil {
			return nil, err
		}
		authors[i] = au
	}
	result["authors"] = authors
	return result, nil
}

type mockArticleLazy mockArticle

func (ma mockArticleLazy) ToMap() (map[string]any, error) {
	result := map[string]any{}
	result["id"] = ma.ID
	result["authors"] = ma.Authors
	return result, nil
}

var benchmarkEvalNestedColumnFixture = mockArticle{
	ID: 1,
	Authors: []mockAuthor{
		{
			ID:        "1",
			FirstName: "Luke",
			LastName:  ToPtr("Skywalker"),
		},
	},
}

var benchmarkEvalSingleFieldSelection = schema.NewNestedObject(map[string]schema.FieldEncoder{
	"id": schema.NewColumnField("id", nil),
}).Encode()

var benchmarkEvalAllFieldsSelection = schema.NewNestedObject(map[string]schema.FieldEncoder{
	"id": schema.NewColumnField("id", nil),
	"authors": schema.NewColumnField("authors", schema.NewNestedArray(schema.NewNestedObject(map[string]schema.FieldEncoder{
		"id":         schema.NewColumnField("id", nil),
		"first_name": schema.NewColumnField("first_name", nil),
		"last_name":  schema.NewColumnField("last_name", nil),
		"address": schema.NewColumnField("address", schema.NewNestedObject(map[string]schema.FieldEncoder{
			"street": schema.NewColumnField("street", nil),
		})),
	}))),
}).Encode()

// BenchmarkEvalSingleFieldToMap-32    	 1390357	       797.6 ns/op	    1512 B/op	      15 allocs/op
func BenchmarkEvalSingleFieldToMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := EvalNestedColumnFields(benchmarkEvalSingleFieldSelection, benchmarkEvalNestedColumnFixture)
		if err != nil {
			panic(err)
		}
	}
}

// BenchmarkEvalSingleFieldToMapLazy-32    	 2420421	       491.2 ns/op	     800 B/op	       8 allocs/op
func BenchmarkEvalSingleFieldToMapLazy(b *testing.B) {
	value := mockArticleLazy(benchmarkEvalNestedColumnFixture)
	for i := 0; i < b.N; i++ {
		_, err := EvalNestedColumnFields(benchmarkEvalSingleFieldSelection, value)
		if err != nil {
			panic(err)
		}
	}
}

// BenchmarkEvalSingleFieldAny-32    	  706143	      1596 ns/op	    1776 B/op	      26 allocs/op
func BenchmarkEvalSingleFieldAny(b *testing.B) {
	type Object mockArticle
	var object = Object(benchmarkEvalNestedColumnFixture)
	for i := 0; i < b.N; i++ {
		_, err := EvalNestedColumnFields(benchmarkEvalSingleFieldSelection, object)
		if err != nil {
			panic(err)
		}
	}
}

// BenchmarkEvalAllFieldToMap-32    	  478923	      2319 ns/op	    2708 B/op	      38 allocs/op
func BenchmarkEvalAllFieldToMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := EvalNestedColumnFields(benchmarkEvalAllFieldsSelection, benchmarkEvalNestedColumnFixture)
		if err != nil {
			panic(err)
		}
	}
}

// BenchmarkEvalAllFieldToMapLazy-32    	  461304	      2345 ns/op	    2764 B/op	      38 allocs/op
func BenchmarkEvalAllFieldToMapLazy(b *testing.B) {
	value := mockArticleLazy(benchmarkEvalNestedColumnFixture)
	for i := 0; i < b.N; i++ {
		_, err := EvalNestedColumnFields(benchmarkEvalAllFieldsSelection, value)
		if err != nil {
			panic(err)
		}
	}
}

// BenchmarkEvalAllFieldAny-32    	  368212	      3279 ns/op	    2972 B/op	      49 allocs/op
func BenchmarkEvalAllFieldAny(b *testing.B) {
	type Object mockArticle
	var object = Object(benchmarkEvalNestedColumnFixture)
	for i := 0; i < b.N; i++ {
		_, err := EvalNestedColumnFields(benchmarkEvalAllFieldsSelection, object)
		if err != nil {
			panic(err)
		}
	}
}
