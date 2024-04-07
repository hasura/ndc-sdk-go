package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/functions"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/scalar"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/stretchr/testify/assert"
)

func createTestServer(t *testing.T) *connector.Server[types.Configuration, types.State] {
	server, err := connector.NewServer[types.Configuration, types.State](&Connector{}, &connector.ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, connector.WithoutRecovery())

	assert.NoError(t, err)

	return server
}

func TestQueryGetTypes(t *testing.T) {
	commentText := functions.CommentText{}
	assert.NoError(t, commentText.FromValue("a comment"))
	commentTextPtr := functions.CommentText{}
	assert.NoError(t, commentTextPtr.FromValue("a comment pointer"))

	testCases := []struct {
		name     string
		body     string
		status   int
		response functions.GetTypesArguments
		errorMsg string
	}{
		{
			name:   "all_fields",
			status: http.StatusOK,
			body: `{
				"collection": "getTypes",
				"arguments": {
					"UUID": {
						"type": "literal",
						"value": "b085b0b9-007c-440e-9661-0d8f2de98a5a"
					},
					"Bool": {
						"type": "literal",
						"value": true
					},
					"String": {
						"type": "literal",
						"value": "hello"
					},
					"Int": {
						"type": "literal",
						"value": 1
					},
					"Int8": {
						"type": "literal",
						"value": 2
					},
					"Int16": {
						"type": "literal",
						"value": 3
					},
					"Int32": {
						"type": "literal",
						"value": 4
					},
					"Int64": {
						"type": "literal",
						"value": 5
					},
					"Uint": {
						"type": "literal",
						"value": 6
					},
					"Uint8": {
						"type": "literal",
						"value": 7
					},
					"Uint16": {
						"type": "literal",
						"value": 8
					},
					"Uint32": {
						"type": "literal",
						"value": 9
					},
					"Uint64": {
						"type": "literal",
						"value": 10
					},
					"Float32": {
						"type": "literal",
						"value": 1.1
					},
					"Float64": {
						"type": "literal",
						"value": 2.2
					},
					"Time": {
						"type": "literal",
						"value": "2024-03-05T07:00:56Z"
					},
					"Duration": {
						"type": "literal",
						"value": "10s"
					},
					"Text": {
						"type": "literal",
						"value": "text"
					},
					"CustomScalar": {
						"type": "literal",
						"value": "a comment"
					},
					"Enum": {
						"type": "literal",
						"value": "foo"
					},
					"BigInt": {
						"type": "literal",
						"value": "10000"
					},
					"Date": {
						"type": "literal",
						"value": "2024-04-02"
					},
					"UUIDPtr": {
						"type": "literal",
						"value": "b085b0b9-007c-440e-9661-0d8f2de98a5b"
					},
			
					"BoolPtr": {
						"type": "literal",
						"value": true
					},
					"StringPtr": {
						"type": "literal",
						"value": "world"
					},
					"IntPtr": {
						"type": "literal",
						"value": 11
					},
					"Int8Ptr": {
						"type": "literal",
						"value": 12
					},
					"Int16Ptr": {
						"type": "literal",
						"value": 13
					},
					"Int32Ptr": {
						"type": "literal",
						"value": 14
					},
					"Int64Ptr": {
						"type": "literal",
						"value": 15
					},
					"UintPtr": {
						"type": "literal",
						"value": 16
					},
					"Uint8Ptr": {
						"type": "literal",
						"value": 17
					},
					"Uint16Ptr": {
						"type": "literal",
						"value": 18
					},
					"Uint32Ptr": {
						"type": "literal",
						"value": 19
					},
					"Uint64Ptr": {
						"type": "literal",
						"value": 20
					},
					"Float32Ptr": {
						"type": "literal",
						"value": 3.3
					},
					"Float64Ptr": {
						"type": "literal",
						"value": 4.4
					},
					"TimePtr": {
						"type": "literal",
						"value": "2024-03-05T07:00:00Z"
					},
					"DurationPtr": {
						"type": "literal",
						"value": "1m"
					},
					"TextPtr": {
						"type": "literal",
						"value": "text pointer"
					},
					"CustomScalarPtr": {
						"type": "literal",
						"value": "a comment pointer"
					},
					"EnumPtr": {
						"type": "literal",
						"value": "bar"
					},
					"BigIntPtr": {
						"type": "literal",
						"value": "20000"
					},
					"DatePtr": {
						"type": "literal",
						"value": "2024-04-03"
					},
					"Object": {
						"type": "literal",
						"value": {
							"id": "b085b0b9-007c-440e-9661-0d8f2de98a5c",
							"created_at": "2024-03-05T06:00:00Z"
						}
					},
					"ObjectPtr": {
						"type": "literal",
						"value": {
							"Long": 1,
							"Lat": 2
						}
					},
					"ArrayObject": {
						"type": "literal",
						"value": [
							{
								"content": "a content"
							}
						]
					},
					"ArrayObjectPtr": {
						"type": "literal",
						"value": [{ "content": "a content pointer" }]
					},
					"NamedObject": {
						"type": "literal",
						"value": {
							"id": "1",
							"duration": 10,
							"created_at": "2024-03-05T05:00:00Z"
						}
					},
					"NamedObjectPtr": {
						"type": "literal",
						"value": {
							"id": "2",
							"duration": 11,
							"created_at": "2024-03-05T04:00:00Z"
						}
					},
					"NamedArray": {
						"type": "literal",
						"value": [
							{
								"id": "3",
								"duration": 12,
								"created_at": "2024-03-05T03:00:00Z"
							}
						]
					},
					"NamedArrayPtr": {
						"type": "literal",
						"value": [
							{
								"created_at": "2024-03-05T02:00:00Z",
								"id": "bPgG5cs38N"
							}
						]
					},
					"UUIDArray": {
						"type": "literal",
						"value": [
							"b085b0b9-007c-440e-9661-0d8f2de98a5a",
							"b085b0b9-007c-440e-9661-0d8f2de98a5b"
						]
					}
				},
				"query": {
					"fields": {
						"__value": {
							"column": "__value",
							"fields": {
								"fields": {
									"ArrayObject": {
										"column": "ArrayObject",
										"fields": {
											"fields": {
												"fields": {
													"content": {
														"column": "content",
														"type": "column"
													}
												},
												"type": "object"
											},
											"type": "array"
										},
										"type": "column"
									},
									"ArrayObjectPtr": {
										"column": "ArrayObjectPtr",
										"fields": {
											"fields": {
												"fields": {
													"content": {
														"column": "content",
														"type": "column"
													}
												},
												"type": "object"
											},
											"type": "array"
										},
										"type": "column"
									},
									"Bool": {
										"column": "Bool",
										"type": "column"
									},
									"BoolPtr": {
										"column": "BoolPtr",
										"type": "column"
									},
									"CustomScalar": {
										"column": "CustomScalar",
										"type": "column"
									},
									"CustomScalarPtr": {
										"column": "CustomScalarPtr",
										"type": "column"
									},
									"Enum": {
										"column": "Enum",
										"type": "column"
									},
									"EnumPtr": {
										"column": "EnumPtr",
										"type": "column"
									},
									"Float32": {
										"column": "Float32",
										"type": "column"
									},
									"Float32Ptr": {
										"column": "Float32Ptr",
										"type": "column"
									},
									"Float64": {
										"column": "Float64",
										"type": "column"
									},
									"Float64Ptr": {
										"column": "Float64Ptr",
										"type": "column"
									},
									"Int": {
										"column": "Int",
										"type": "column"
									},
									"Int16": {
										"column": "Int16",
										"type": "column"
									},
									"Int16Ptr": {
										"column": "Int16Ptr",
										"type": "column"
									},
									"Int32": {
										"column": "Int32",
										"type": "column"
									},
									"Int32Ptr": {
										"column": "Int32Ptr",
										"type": "column"
									},
									"Int64": {
										"column": "Int64",
										"type": "column"
									},
									"Int64Ptr": {
										"column": "Int64Ptr",
										"type": "column"
									},
									"Int8": {
										"column": "Int8",
										"type": "column"
									},
									"Int8Ptr": {
										"column": "Int8Ptr",
										"type": "column"
									},
									"IntPtr": {
										"column": "IntPtr",
										"type": "column"
									},
									"BigInt": {
										"column": "BigInt",
										"type": "column"
									},
									"BigIntPtr": {
										"column": "BigIntPtr",
										"type": "column"
									},
									"Date": {
										"column": "Date",
										"type": "column"
									},
									"DatePtr": {
										"column": "DatePtr",
										"type": "column"
									},
									"NamedArray": {
										"column": "NamedArray",
										"fields": {
											"fields": {
												"fields": {
													"created_at": {
														"column": "created_at",
														"type": "column"
													},
													"id": {
														"column": "id",
														"type": "column"
													}
												},
												"type": "object"
											},
											"type": "array"
										},
										"type": "column"
									},
									"NamedArrayPtr": {
										"column": "NamedArrayPtr",
										"fields": {
											"fields": {
												"fields": {
													"created_at": {
														"column": "created_at",
														"type": "column"
													},
													"id": {
														"column": "id",
														"type": "column"
													}
												},
												"type": "object"
											},
											"type": "array"
										},
										"type": "column"
									},
									"NamedObject": {
										"column": "NamedObject",
										"fields": {
											"fields": {
												"created_at": {
													"column": "created_at",
													"type": "column"
												},
												"id": {
													"column": "id",
													"type": "column"
												}
											},
											"type": "object"
										},
										"type": "column"
									},
									"NamedObjectPtr": {
										"column": "NamedObjectPtr",
										"fields": {
											"fields": {
												"created_at": {
													"column": "created_at",
													"type": "column"
												},
												"id": {
													"column": "id",
													"type": "column"
												}
											},
											"type": "object"
										},
										"type": "column"
									},
									"Object": {
										"column": "Object",
										"fields": {
											"fields": {
												"created_at": {
													"column": "created_at",
													"type": "column"
												},
												"id": {
													"column": "id",
													"type": "column"
												}
											},
											"type": "object"
										},
										"type": "column"
									},
									"ObjectPtr": {
										"column": "ObjectPtr",
										"fields": {
											"fields": {
												"Lat": {
													"column": "Lat",
													"type": "column"
												},
												"Long": {
													"column": "Long",
													"type": "column"
												}
											},
											"type": "object"
										},
										"type": "column"
									},
									"String": {
										"column": "String",
										"type": "column"
									},
									"StringPtr": {
										"column": "StringPtr",
										"type": "column"
									},
									"Text": {
										"column": "Text",
										"type": "column"
									},
									"TextPtr": {
										"column": "TextPtr",
										"type": "column"
									},
									"Time": {
										"column": "Time",
										"type": "column"
									},
									"TimePtr": {
										"column": "TimePtr",
										"type": "column"
									},
									"UUID": {
										"column": "UUID",
										"type": "column"
									},
									"UUIDArray": {
										"column": "UUIDArray",
										"type": "column"
									},
									"UUIDPtr": {
										"column": "UUIDPtr",
										"type": "column"
									},
									"Uint": {
										"column": "Uint",
										"type": "column"
									},
									"Uint16": {
										"column": "Uint16",
										"type": "column"
									},
									"Uint16Ptr": {
										"column": "Uint16Ptr",
										"type": "column"
									},
									"Uint32": {
										"column": "Uint32",
										"type": "column"
									},
									"Uint32Ptr": {
										"column": "Uint32Ptr",
										"type": "column"
									},
									"Uint64": {
										"column": "Uint64",
										"type": "column"
									},
									"Uint64Ptr": {
										"column": "Uint64Ptr",
										"type": "column"
									},
									"Uint8": {
										"column": "Uint8",
										"type": "column"
									},
									"Uint8Ptr": {
										"column": "Uint8Ptr",
										"type": "column"
									},
									"UintPtr": {
										"column": "UintPtr",
										"type": "column"
									}
								},
								"type": "object"
							},
							"type": "column"
						}
					}
				},
				"collection_relationships": {}
			}`,
			response: functions.GetTypesArguments{
				UUID:            uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5a"),
				Bool:            true,
				String:          "hello",
				Int:             1,
				Int8:            2,
				Int16:           3,
				Int32:           4,
				Int64:           5,
				Uint:            6,
				Uint8:           7,
				Uint16:          8,
				Uint32:          9,
				Uint64:          10,
				Float32:         1.1,
				Float64:         2.2,
				Time:            time.Date(2024, 3, 5, 7, 0, 56, 0, time.UTC),
				Text:            "text",
				CustomScalar:    commentText,
				Enum:            functions.SomeEnumFoo,
				BigInt:          10000,
				Date:            *scalar.NewDate(2024, 04, 02),
				UUIDPtr:         utils.ToPtr(uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5b")),
				BoolPtr:         utils.ToPtr(true),
				StringPtr:       utils.ToPtr("world"),
				IntPtr:          utils.ToPtr(11),
				Int8Ptr:         utils.ToPtr(int8(12)),
				Int16Ptr:        utils.ToPtr(int16(13)),
				Int32Ptr:        utils.ToPtr(int32(14)),
				Int64Ptr:        utils.ToPtr(int64(15)),
				UintPtr:         utils.ToPtr(uint(16)),
				Uint8Ptr:        utils.ToPtr(uint8(17)),
				Uint16Ptr:       utils.ToPtr(uint16(18)),
				Uint32Ptr:       utils.ToPtr(uint32(19)),
				Uint64Ptr:       utils.ToPtr(uint64(20)),
				Float32Ptr:      utils.ToPtr(float32(3.3)),
				Float64Ptr:      utils.ToPtr(float64(4.4)),
				TimePtr:         utils.ToPtr(time.Date(2024, 3, 5, 7, 0, 0, 0, time.UTC)),
				TextPtr:         utils.ToPtr(functions.Text("text pointer")),
				CustomScalarPtr: &commentTextPtr,
				EnumPtr:         utils.ToPtr(functions.SomeEnumBar),
				BigIntPtr:       utils.ToPtr(scalar.BigInt(20000)),
				DatePtr:         scalar.NewDate(2024, 04, 03),
				Object: struct {
					ID        uuid.UUID `json:"id"`
					CreatedAt time.Time `json:"created_at"`
				}{
					ID:        uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5c"),
					CreatedAt: time.Date(2024, 3, 5, 6, 0, 0, 0, time.UTC),
				},
				ObjectPtr: &struct {
					Long int
					Lat  int
				}{
					Long: 1,
					Lat:  2,
				},
				ArrayObject: []struct {
					Content string `json:"content"`
				}{
					{
						Content: "a content",
					},
				},
				ArrayObjectPtr: &[]struct {
					Content string `json:"content"`
				}{
					{
						Content: "a content pointer",
					},
				},
				NamedObject: functions.Author{
					ID:        "1",
					CreatedAt: time.Date(2024, 3, 5, 5, 0, 0, 0, time.UTC),
				},
				NamedObjectPtr: &functions.Author{
					ID:        "2",
					CreatedAt: time.Date(2024, 3, 5, 4, 0, 0, 0, time.UTC),
				},
				NamedArray: []functions.Author{
					{
						ID:        "3",
						CreatedAt: time.Date(2024, 3, 5, 3, 0, 0, 0, time.UTC),
					},
				},
				NamedArrayPtr: &[]functions.Author{
					{
						ID:        "bPgG5cs38N",
						CreatedAt: time.Date(2024, 3, 5, 2, 0, 0, 0, time.UTC),
					},
				},
				UUIDArray: []uuid.UUID{
					uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5a"),
					uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5b"),
				},
			},
		},
	}

	testServer := createTestServer(t).BuildTestServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/query", testServer.URL), "application/json", bytes.NewReader([]byte(tc.body)))
			assert.NoError(t, err, "failed to request query")
			assert.Equal(t, tc.status, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if tc.errorMsg != "" {
				assert.NoError(t, err)
				assert.Contains(t, string(respBody), tc.errorMsg)
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("expected successful response, got error: %s", string(respBody))
			} else {
				log.Print("response: ", string(respBody))
				var results []struct {
					Rows []struct {
						Value functions.GetTypesArguments `json:"__value"`
					} `json:"rows,omitempty" mapstructure:"rows,omitempty"`
				}
				assert.NoError(t, json.Unmarshal(respBody, &results), "failed to decode response")
				assert.Equal(t, 1, len(results))
				assert.Equal(t, 1, len(results[0].Rows))
				assert.Equal(t, tc.response, results[0].Rows[0].Value)
			}
		})
	}
}

func TestQueries(t *testing.T) {
	testCases := []struct {
		name     string
		body     string
		status   int
		response string
		errorMsg string
	}{
		{
			name:   "hello_success",
			status: http.StatusOK,
			body: `{
				"collection": "hello",
				"query": {
					"fields": {
						"__value": {
							"type": "column",
							"column": "__value",
							"fields": {
								"type": "object",
								"fields": {
									"num": { "type": "column", "column": "num", "fields": null },
									"text": { "type": "column", "column": "text", "fields": null }
								}
							}
						}
					}
				},
				"arguments": {},
				"collection_relationships": {}
			}`,
			response: `{
				"num": 1,
				"text": "world"
			}`,
		},
		{
			name:   "hello_failure_array",
			status: http.StatusUnprocessableEntity,
			body: `{
				"collection": "hello",
				"query": {
					"fields": {
						"__value": {
							"type": "column",
							"column": "__value",
							"fields": {
								"type": "array",
								"fields": {
									"type": "object",
									"fields": {
										"num": { "type": "column", "column": "num", "fields": null },
										"text": { "type": "column", "column": "text", "fields": null }
									}
								}
							}
						}
					}
				},
				"arguments": {},
				"collection_relationships": {}
			}`,
			errorMsg: "the selection field type must be object",
		},
		{
			name:   "getBool",
			status: http.StatusOK,
			body: `{
				"collection": "getBool",
				"arguments": {},
				"query": {
					"fields": {
						"__value": {
							"type": "column",
							"column": "__value",
							"fields": null
						}
					}
				},
				"collection_relationships": {}
			}`,
			response: `true`,
		},
		{
			name:   "getArticles_failure_object",
			status: http.StatusUnprocessableEntity,
			body: `{
				"collection": "getArticles",
				"query": {
					"fields": {
						"__value": {
							"type": "column",
							"column": "__value",
							"fields": {
								"type": "object",
								"fields": {
									"num": { "type": "column", "column": "num", "fields": null },
									"text": { "type": "column", "column": "text", "fields": null }
								}
							}
						}
					}
				},
				"arguments": {},
				"collection_relationships": {}
			}`,
			errorMsg: "the selection field type must be array",
		},
		{
			name:   "getArticles_success",
			status: http.StatusOK,
			body: `{
				"collection": "getArticles",
				"query": {
					"fields": {
						"__value": {
							"type": "column",
							"column": "__value",
							"fields": {
								"type": "array",
								"fields": {
									"type": "object",
									"fields": {
										"id": { "type": "column", "column": "id" }
									}
								}
							}
						}
					}
				},
				"arguments": {
					"Limit": {
						"type": "literal",
						"value": 1
					}
				},
				"collection_relationships": {}
			}`,
			response: `[{
				"id": "1"
			}]`,
		},
	}

	testServer := createTestServer(t).BuildTestServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/query", testServer.URL), "application/json", bytes.NewReader([]byte(tc.body)))
			assert.NoError(t, err, "failed to request query")
			assert.Equal(t, tc.status, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if tc.errorMsg != "" {
				assert.NoError(t, err)
				assert.Contains(t, string(respBody), tc.errorMsg)
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("expected successful response, got error: %s", string(respBody))
			} else {
				log.Print("response: ", string(respBody))
				var expected any
				assert.NoError(t, json.Unmarshal([]byte(tc.response), &expected))
				var results []struct {
					Rows []struct {
						Value any `json:"__value"`
					} `json:"rows,omitempty" mapstructure:"rows,omitempty"`
				}
				assert.NoError(t, json.Unmarshal(respBody, &results), "failed to decode response")
				assert.Equal(t, 1, len(results))
				assert.Equal(t, 1, len(results[0].Rows))
				assert.Equal(t, expected, results[0].Rows[0].Value)
			}
		})
	}
}

func TestProcedures(t *testing.T) {
	testCases := []struct {
		name     string
		body     string
		status   int
		response string
		errorMsg string
	}{
		{
			name:   "create_article_success",
			status: http.StatusOK,
			body: `{
				"operations": [
					{
						"type": "procedure",
						"name": "create_article",
						"arguments": {},
						"fields": {
							"type": "object",
							"fields": {
								"id": {
										"type": "column",
										"column": "id"
								}
							}
						}
					}
				],
				"collection_relationships": {}
			}`,
			response: `{
				"id": 1
			}`,
		},
		{
			name:   "create_article_array_400",
			status: http.StatusUnprocessableEntity,
			body: `{
				"operations": [
					{
						"type": "procedure",
						"name": "create_article",
						"arguments": {},
						"fields": {
							"type": "array",
							"fields": {
								"type": "object",
								"fields": {
									"id": {
										"type": "column",
										"column": "id"
									}
								}
							}
						}
					}
				],
				"collection_relationships": {}
			}`,
			errorMsg: "the selection field type must be object",
		},
		{
			name:   "createAuthors_success",
			status: http.StatusOK,
			body: `{
				"operations": [
					{
						"type": "procedure",
						"name": "createAuthors",
						"arguments": {},
						"fields": {
							"type": "array",
							"fields": {
								"type": "object",
								"fields": {
									"id": {
										"type": "column",
										"column": "id"
									}
								}
							}
						}
					}
				],
				"collection_relationships": {}
			}`,
			response: `[{
				"id": 1
			}]`,
		},
		{
			name:   "createAuthors_object_422",
			status: http.StatusUnprocessableEntity,
			body: `{
				"operations": [
					{
						"type": "procedure",
						"name": "createAuthors",
						"arguments": {},
						"fields": {
							"type": "object",
							"fields": {
								"id": {
										"type": "column",
										"column": "id"
								}
							}
						}
					}
				],
				"collection_relationships": {}
			}`,
			errorMsg: "the selection field type must be array",
		},
		{
			name:   "increase_success",
			status: http.StatusOK,
			body: `{
				"operations": [
					{
						"type": "procedure",
						"name": "increase",
						"arguments": {},
						"fields": null
					}
				],
				"collection_relationships": {}
			}`,
			response: `1`,
		},
		{
			name:   "increase_422",
			status: http.StatusUnprocessableEntity,
			body: `{
				"operations": [
					{
						"type": "procedure",
						"name": "increase",
						"arguments": {},
						"fields": {
							"type": "object",
							"fields": {
								"id": {
										"type": "column",
										"column": "id"
								}
							}
						}
					}
				],
				"collection_relationships": {}
			}`,
			errorMsg: "cannot evaluate selection fields for scalar",
		},
	}

	testServer := createTestServer(t).BuildTestServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/mutation", testServer.URL), "application/json", bytes.NewReader([]byte(tc.body)))
			assert.NoError(t, err, "failed to request mutation")
			assert.Equal(t, tc.status, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if tc.errorMsg != "" {
				assert.NoError(t, err)
				assert.Contains(t, string(respBody), tc.errorMsg)
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("expected successful response, got error: %s", string(respBody))
			} else {
				log.Print("response: ", string(respBody))
				var expected any
				assert.NoError(t, json.Unmarshal([]byte(tc.response), &expected))
				var results struct {
					OperationResults []struct {
						Result any `json:"result"`
					} `json:"operation_results"`
				}
				assert.NoError(t, json.Unmarshal(respBody, &results), "failed to decode response")
				assert.Equal(t, 1, len(results.OperationResults))
				assert.Equal(t, expected, results.OperationResults[0].Result)
			}
		})
	}
}
