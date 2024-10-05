package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-codegen-example/types/arguments"
	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/ndctest"
	"github.com/hasura/ndc-sdk-go/scalar"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"gotest.tools/v3/assert"
)

func createTestServer(t *testing.T, options ...connector.ServeOption) *connector.Server[types.Configuration, types.State] {
	// reset global envs
	loadGlobalEnvOnce = sync.Once{}
	_globalEnvironments = globalEnvironments{}

	server, err := connector.NewServer[types.Configuration, types.State](&Connector{}, &connector.ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, append(options, connector.WithoutRecovery())...)

	assert.NilError(t, err)

	return server
}

func TestConnector(t *testing.T) {
	ndctest.TestConnector(t, &Connector{}, ndctest.TestConnectorOptions{
		Configuration: "{}",
		InlineConfig:  true,
		TestDataDir:   "testdata",
	})
}

func TestQueryGetTypes(t *testing.T) {
	commentText := types.CommentText{}
	assert.NilError(t, commentText.FromValue("a comment"))
	commentTextPtr := types.CommentText{}
	assert.NilError(t, commentTextPtr.FromValue("a comment pointer"))

	testCases := []struct {
		name     string
		body     string
		status   int
		response arguments.GetTypesArguments
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
					"URL": {
						"type": "literal",
						"value": "https://example.com"
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
							"tags": [],
							"created_at": "2024-03-05T05:00:00Z"
						}
					},
					"NamedObjectPtr": {
						"type": "literal",
						"value": {
							"id": "2",
							"duration": 11,
							"tags": [],
							"created_at": "2024-03-05T04:00:00Z"
						}
					},
					"NamedArray": {
						"type": "literal",
						"value": [
							{
								"id": "3",
								"duration": 12,
								"tags": [],
								"created_at": "2024-03-05T03:00:00Z"
							}
						]
					},
					"NamedArrayPtr": {
						"type": "literal",
						"value": [
							{
								"created_at": "2024-03-05T02:00:00Z",
								"tags": [],
								"id": "bPgG5cs38N"
							}
						]
					},
					"ArrayUUID": {
						"type": "literal",
						"value": [
							"b085b0b9-007c-440e-9661-0d8f2de98a5a",
							"b085b0b9-007c-440e-9661-0d8f2de98a5b"
						]
					},
					"ArrayUUIDPtr": {
						"type": "literal",
						"value": [
							"b085b0b9-007c-440e-9661-0d8f2de98a5c",
							"b085b0b9-007c-440e-9661-0d8f2de98a5d"
						]
					},
					"Map": {
						"type": "literal",
						"value": {
							"foo": "bar"
						}
					},
					"MapPtr": {
						"type": "literal",
						"value": {
							"foo": "bar_ptr"
						}
					},
					"ArrayMap": {
						"type": "literal",
						"value": [{
							"foo": "bar"
						}]
					},
					"ArrayMapPtr": {
						"type": "literal",
						"value": [{
							"foo": "bar_ptr"
						}]
					},
					"JSON": {
						"type": "literal",
						"value": {"message":"json"}
					},
					"JSONPtr": {
						"type": "literal",
						"value": {"message":"json_ptr"}
					},
					"RawJSON": {
						"type": "literal",
						"value": {"message":"raw_json"}
					},
					"RawJSONPtr": {
						"type": "literal",
						"value": {"message":"raw_json_ptr"}
					},
					"Bytes": {
						"type": "literal",
						"value": "aGVsbG8gd29ybGQ="
					},
					"BytesPtr": {
						"type": "literal",
						"value": "aGVsbG8gcG9pbnRlcg=="
					},
					"ArrayBool": {
						"type": "literal",
						"value": [true]
					},
					"ArrayString": {
						"type": "literal",
						"value": ["foo"]
					},
					"ArrayInt": {
						"type": "literal",
						"value": [1]
					},
					"ArrayInt8": {
						"type": "literal",
						"value": [2]
					},
					"ArrayInt16": {
						"type": "literal",
						"value": [3]
					},
					"ArrayInt32": {
						"type": "literal",
						"value": [4]
					},
					"ArrayInt64": {
						"type": "literal",
						"value": [5]
					},
					"ArrayUint": {
						"type": "literal",
						"value": [6]
					},
					"ArrayUint8": {
						"type": "literal",
						"value": [7]
					},
					"ArrayUint16": {
						"type": "literal",
						"value": [8]
					},
					"ArrayUint32": {
						"type": "literal",
						"value": [9]
					},
					"ArrayUint64": {
						"type": "literal",
						"value": [10]
					},
					"ArrayFloat32": {
						"type": "literal",
						"value": [11.1]
					},
					"ArrayFloat64": {
						"type": "literal",
						"value": [12.2]
					},
					"ArrayBoolPtr": {
						"type": "literal",
						"value": [true]
					},
					"ArrayStringPtr": {
						"type": "literal",
						"value": ["foo"]
					},
					"ArrayIntPtr": {
						"type": "literal",
						"value": [1]
					},
					"ArrayInt8Ptr": {
						"type": "literal",
						"value": [2]
					},
					"ArrayInt16Ptr": {
						"type": "literal",
						"value": [3]
					},
					"ArrayInt32Ptr": {
						"type": "literal",
						"value": [4]
					},
					"ArrayInt64Ptr": {
						"type": "literal",
						"value": [5]
					},
					"ArrayUintPtr": {
						"type": "literal",
						"value": [6]
					},
					"ArrayUint8Ptr": {
						"type": "literal",
						"value": [7]
					},
					"ArrayUint16Ptr": {
						"type": "literal",
						"value": [8]
					},
					"ArrayUint32Ptr": {
						"type": "literal",
						"value": [9]
					},
					"ArrayUint64Ptr": {
						"type": "literal",
						"value": [10]
					},
					"ArrayFloat32Ptr": {
						"type": "literal",
						"value": [11.1]
					},
					"ArrayFloat64Ptr": {
						"type": "literal",
						"value": [12.2]
					},
					"ArrayBigInt": {
						"type": "literal",
						"value": ["13"]
					},
					"ArrayBigIntPtr": {
						"type": "literal",
						"value": ["14"]
					},
					"ArrayJSON": {
						"type": "literal",
						"value": [{"foo": "bar"}]
					},
					"ArrayJSONPtr": {
						"type": "literal",
						"value": [{"foo": "baz"}]
					},
					"ArrayRawJSON": {
						"type": "literal",
						"value": [{"message":"raw_json"}]
					},
					"ArrayRawJSONPtr": {
						"type": "literal",
						"value": [{"message":"raw_json_ptr"}]
					},

					"PtrArrayBool": {
						"type": "literal",
						"value": [true]
					},
					"PtrArrayString": {
						"type": "literal",
						"value": ["foo"]
					},
					"PtrArrayInt": {
						"type": "literal",
						"value": [1]
					},
					"PtrArrayInt8": {
						"type": "literal",
						"value": [2]
					},
					"PtrArrayInt16": {
						"type": "literal",
						"value": [3]
					},
					"PtrArrayInt32": {
						"type": "literal",
						"value": [4]
					},
					"PtrArrayInt64": {
						"type": "literal",
						"value": [5]
					},
					"PtrArrayUint": {
						"type": "literal",
						"value": [6]
					},
					"PtrArrayUint8": {
						"type": "literal",
						"value": [7]
					},
					"PtrArrayUint16": {
						"type": "literal",
						"value": [8]
					},
					"PtrArrayUint32": {
						"type": "literal",
						"value": [9]
					},
					"PtrArrayUint64": {
						"type": "literal",
						"value": [10]
					},
					"PtrArrayFloat32": {
						"type": "literal",
						"value": [11.1]
					},
					"PtrArrayFloat64": {
						"type": "literal",
						"value": [12.2]
					},
					"PtrArrayBoolPtr": {
						"type": "literal",
						"value": [true]
					},
					"PtrArrayStringPtr": {
						"type": "literal",
						"value": ["foo"]
					},
					"PtrArrayIntPtr": {
						"type": "literal",
						"value": [1]
					},
					"PtrArrayInt8Ptr": {
						"type": "literal",
						"value": [2]
					},
					"PtrArrayInt16Ptr": {
						"type": "literal",
						"value": [3]
					},
					"PtrArrayInt32Ptr": {
						"type": "literal",
						"value": [4]
					},
					"PtrArrayInt64Ptr": {
						"type": "literal",
						"value": [5]
					},
					"PtrArrayUintPtr": {
						"type": "literal",
						"value": [6]
					},
					"PtrArrayUint8Ptr": {
						"type": "literal",
						"value": [7]
					},
					"PtrArrayUint16Ptr": {
						"type": "literal",
						"value": [8]
					},
					"PtrArrayUint32Ptr": {
						"type": "literal",
						"value": [9]
					},
					"PtrArrayUint64Ptr": {
						"type": "literal",
						"value": [10]
					},
					"PtrArrayFloat32Ptr": {
						"type": "literal",
						"value": [11.1]
					},
					"PtrArrayFloat64Ptr": {
						"type": "literal",
						"value": [12.2]
					},
					"PtrArrayBigInt": {
						"type": "literal",
						"value": ["13"]
					},
					"PtrArrayBigIntPtr": {
						"type": "literal",
						"value": ["14"]
					},
					"PtrArrayJSON": {
						"type": "literal",
						"value": [{"foo": "bar"}]
					},
					"PtrArrayJSONPtr": {
						"type": "literal",
						"value": [{"foo": "baz"}]
					},
					"PtrArrayRawJSON": {
						"type": "literal",
						"value": [{"message":"raw_json"}]
					},
					"PtrArrayRawJSONPtr": {
						"type": "literal",
						"value": [{"message":"raw_json_ptr"}]
					},
					"ArrayTime": {
						"type": "literal",
						"value": ["2024-03-05T07:00:56Z"]
					},
					"ArrayTimePtr": {
						"type": "literal",
						"value": ["2024-03-05T07:00:56Z"]
					},
					"PtrArrayTime": {
						"type": "literal",
						"value": ["2024-03-05T07:00:56Z"]
					},
					"PtrArrayTimePtr": {
						"type": "literal",
						"value": ["2024-03-05T07:00:56Z"]
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
									"URL": {
										"column": "URL",
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
									"ArrayUUID": {
										"column": "ArrayUUID",
										"type": "column"
									},
									"ArrayUUIDPtr": {
										"column": "ArrayUUIDPtr",
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
									},
									"Map": {
										"column": "Map",
										"type": "column"
									},
									"MapPtr": {
										"column": "MapPtr",
										"type": "column"
									},
									"ArrayMap": {
										"column": "ArrayMap",
										"type": "column"
									},
									"ArrayMapPtr": {
										"column": "ArrayMapPtr",
										"type": "column"
									},
									"JSON": {
										"column": "JSON",
										"type": "column"
									},
									"JSONPtr": {
										"column": "JSONPtr",
										"type": "column"
									},
									"RawJSON": {
										"column": "RawJSON",
										"type": "column"
									},
									"RawJSONPtr": {
										"column": "RawJSONPtr",
										"type": "column"
									},
									"Bytes": {
										"column": "Bytes",
										"type": "column"
									},
									"BytesPtr": {
										"column": "BytesPtr",
										"type": "column"
									},
									"ArrayBool": {
										"column": "ArrayBool",
										"type": "column"
									},
									"ArrayString": {
										"column": "ArrayString",
										"type": "column"
									},
									"ArrayInt": {
										"column": "ArrayInt",
										"type": "column"
									},
									"ArrayInt8": {
										"column": "ArrayInt8",
										"type": "column"
									},
									"ArrayInt16": {
										"column": "ArrayInt16",
										"type": "column"
									},
									"ArrayInt32": {
										"column": "ArrayInt32",
										"type": "column"
									},
									"ArrayInt64": {
										"column": "ArrayInt64",
										"type": "column"
									},
									"ArrayUint": {
										"column": "ArrayUint",
										"type": "column"
									},
									"ArrayUint8": {
										"column": "ArrayUint8",
										"type": "column"
									},
									"ArrayUint16": {
										"column": "ArrayUint16",
										"type": "column"
									},
									"ArrayUint32": {
										"column": "ArrayUint32",
										"type": "column"
									},
									"ArrayUint64": {
										"column": "ArrayUint64",
										"type": "column"
									},
									"ArrayFloat32": {
										"column": "ArrayFloat32",
										"type": "column"
									},
									"ArrayFloat64": {
										"column": "ArrayFloat64",
										"type": "column"
									},
									"ArrayBoolPtr": {
										"column": "ArrayBoolPtr",
										"type": "column"
									},
									"ArrayStringPtr": {
										"column": "ArrayStringPtr",
										"type": "column"
									},
									"ArrayIntPtr": {
										"column": "ArrayIntPtr",
										"type": "column"
									},
									"ArrayInt8Ptr": {
										"column": "ArrayInt8Ptr",
										"type": "column"
									},
									"ArrayInt16Ptr": {
										"column": "ArrayInt16Ptr",
										"type": "column"
									},
									"ArrayInt32Ptr": {
										"column": "ArrayInt32Ptr",
										"type": "column"
									},
									"ArrayInt64Ptr": {
										"column": "ArrayInt64Ptr",
										"type": "column"
									},
									"ArrayUintPtr": {
										"column": "ArrayUintPtr",
										"type": "column"
									},
									"ArrayUint8Ptr": {
										"column": "ArrayUint8Ptr",
										"type": "column"
									},
									"ArrayUint16Ptr": {
										"column": "ArrayUint16Ptr",
										"type": "column"
									},
									"ArrayUint32Ptr": {
										"column": "ArrayUint32Ptr",
										"type": "column"
									},
									"ArrayUint64Ptr": {
										"column": "ArrayUint64Ptr",
										"type": "column"
									},
									"ArrayFloat32Ptr": {
										"column": "ArrayFloat32Ptr",
										"type": "column"
									},
									"ArrayFloat64Ptr": {
										"column": "ArrayFloat64Ptr",
										"type": "column"
									},
									"ArrayBigInt": {
										"column": "ArrayBigInt",
										"type": "column"
									},
									"ArrayBigIntPtr": {
										"column": "ArrayBigIntPtr",
										"type": "column"
									},
									"ArrayJSON": {
										"column": "ArrayJSON",
										"type": "column"
									},
									"ArrayJSONPtr": {
										"column": "ArrayJSONPtr",
										"type": "column"
									},
									"ArrayRawJSON": {
										"column": "ArrayRawJSON",
										"type": "column"
									},
									"ArrayRawJSONPtr": {
										"column": "ArrayRawJSONPtr",
										"type": "column"
									},
									"ArrayTime": {
										"column": "ArrayTime",
										"type": "column"
									},
									"ArrayTimePtr": {
										"column": "ArrayTimePtr",
										"type": "column"
									},

									"PtrArrayBool": {
										"column": "PtrArrayBool",
										"type": "column"
									},
									"PtrArrayString": {
										"column": "PtrArrayString",
										"type": "column"
									},
									"PtrArrayInt": {
										"column": "PtrArrayInt",
										"type": "column"
									},
									"PtrArrayInt8": {
										"column": "PtrArrayInt8",
										"type": "column"
									},
									"PtrArrayInt16": {
										"column": "PtrArrayInt16",
										"type": "column"
									},
									"PtrArrayInt32": {
										"column": "PtrArrayInt32",
										"type": "column"
									},
									"PtrArrayInt64": {
										"column": "PtrArrayInt64",
										"type": "column"
									},
									"PtrArrayUint": {
										"column": "PtrArrayUint",
										"type": "column"
									},
									"PtrArrayUint8": {
										"column": "PtrArrayUint8",
										"type": "column"
									},
									"PtrArrayUint16": {
										"column": "PtrArrayUint16",
										"type": "column"
									},
									"PtrArrayUint32": {
										"column": "PtrArrayUint32",
										"type": "column"
									},
									"PtrArrayUint64": {
										"column": "PtrArrayUint64",
										"type": "column"
									},
									"PtrArrayFloat32": {
										"column": "PtrArrayFloat32",
										"type": "column"
									},
									"PtrArrayFloat64": {
										"column": "PtrArrayFloat64",
										"type": "column"
									},
									"PtrArrayBoolPtr": {
										"column": "PtrArrayBoolPtr",
										"type": "column"
									},
									"PtrArrayStringPtr": {
										"column": "PtrArrayStringPtr",
										"type": "column"
									},
									"PtrArrayIntPtr": {
										"column": "PtrArrayIntPtr",
										"type": "column"
									},
									"PtrArrayInt8Ptr": {
										"column": "PtrArrayInt8Ptr",
										"type": "column"
									},
									"PtrArrayInt16Ptr": {
										"column": "PtrArrayInt16Ptr",
										"type": "column"
									},
									"PtrArrayInt32Ptr": {
										"column": "PtrArrayInt32Ptr",
										"type": "column"
									},
									"PtrArrayInt64Ptr": {
										"column": "PtrArrayInt64Ptr",
										"type": "column"
									},
									"PtrArrayUintPtr": {
										"column": "PtrArrayUintPtr",
										"type": "column"
									},
									"PtrArrayUint8Ptr": {
										"column": "PtrArrayUint8Ptr",
										"type": "column"
									},
									"PtrArrayUint16Ptr": {
										"column": "PtrArrayUint16Ptr",
										"type": "column"
									},
									"PtrArrayUint32Ptr": {
										"column": "PtrArrayUint32Ptr",
										"type": "column"
									},
									"PtrArrayUint64Ptr": {
										"column": "PtrArrayUint64Ptr",
										"type": "column"
									},
									"PtrArrayFloat32Ptr": {
										"column": "PtrArrayFloat32Ptr",
										"type": "column"
									},
									"PtrArrayFloat64Ptr": {
										"column": "PtrArrayFloat64Ptr",
										"type": "column"
									},
									"PtrArrayBigInt": {
										"column": "PtrArrayBigInt",
										"type": "column"
									},
									"PtrArrayBigIntPtr": {
										"column": "PtrArrayBigIntPtr",
										"type": "column"
									},
									"PtrArrayJSON": {
										"column": "PtrArrayJSON",
										"type": "column"
									},
									"PtrArrayJSONPtr": {
										"column": "PtrArrayJSONPtr",
										"type": "column"
									},
									"PtrArrayRawJSON": {
										"column": "PtrArrayRawJSON",
										"type": "column"
									},
									"PtrArrayRawJSONPtr": {
										"column": "PtrArrayRawJSONPtr",
										"type": "column"
									},
									"PtrArrayTime": {
										"column": "PtrArrayTime",
										"type": "column"
									},
									"PtrArrayTimePtr": {
										"column": "PtrArrayTimePtr",
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
			response: arguments.GetTypesArguments{
				UUID:         uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5a"),
				Bool:         true,
				String:       "hello",
				Int:          1,
				Int8:         2,
				Int16:        3,
				Int32:        4,
				Int64:        5,
				Uint:         6,
				Uint8:        7,
				Uint16:       8,
				Uint32:       9,
				Uint64:       10,
				Float32:      1.1,
				Float64:      2.2,
				Time:         time.Date(2024, 3, 5, 7, 0, 56, 0, time.UTC),
				Text:         "text",
				CustomScalar: commentText,
				Enum:         types.SomeEnumFoo,
				BigInt:       10000,
				Date:         *scalar.NewDate(2024, 04, 02),
				URL: func() scalar.URL {
					r, _ := scalar.NewURL("https://example.com")
					return *r
				}(),

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
				TextPtr:         utils.ToPtr(types.Text("text pointer")),
				CustomScalarPtr: &commentTextPtr,
				EnumPtr:         utils.ToPtr(types.SomeEnumBar),
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
				NamedObject: types.Author{
					ID:        "1",
					CreatedAt: time.Date(2024, 3, 5, 5, 0, 0, 0, time.UTC),
				},
				NamedObjectPtr: &types.Author{
					ID:        "2",
					CreatedAt: time.Date(2024, 3, 5, 4, 0, 0, 0, time.UTC),
				},
				NamedArray: []types.Author{
					{
						ID:        "3",
						CreatedAt: time.Date(2024, 3, 5, 3, 0, 0, 0, time.UTC),
					},
				},
				NamedArrayPtr: &[]types.Author{
					{
						ID:        "bPgG5cs38N",
						CreatedAt: time.Date(2024, 3, 5, 2, 0, 0, 0, time.UTC),
					},
				},
				ArrayUUID: []uuid.UUID{
					uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5a"),
					uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5b"),
				},
				ArrayUUIDPtr: []*uuid.UUID{
					utils.ToPtr(uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5c")),
					utils.ToPtr(uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5d")),
				},
				Map: map[string]any{
					"foo": "bar",
				},
				MapPtr: &map[string]any{
					"foo": "bar_ptr",
				},
				ArrayMap: []map[string]any{
					{"foo": "bar"},
				},
				ArrayMapPtr: &[]map[string]any{
					{"foo": "bar_ptr"},
				},
				JSON:            map[string]any{"message": "json"},
				JSONPtr:         utils.ToPtr(any(map[string]any{"message": "json_ptr"})),
				RawJSON:         json.RawMessage(`{"message":"raw_json"}`),
				RawJSONPtr:      utils.ToPtr(json.RawMessage(`{"message":"raw_json_ptr"}`)),
				Bytes:           *scalar.NewBytes([]byte("hello world")),
				BytesPtr:        scalar.NewBytes([]byte("hello pointer")),
				ArrayBool:       []bool{true},
				ArrayString:     []string{"foo"},
				ArrayInt:        []int{1},
				ArrayInt8:       []int8{2},
				ArrayInt16:      []int16{3},
				ArrayInt32:      []int32{4},
				ArrayInt64:      []int64{5},
				ArrayUint:       []uint{6},
				ArrayUint8:      []uint8{7},
				ArrayUint16:     []uint16{8},
				ArrayUint32:     []uint32{9},
				ArrayUint64:     []uint64{10},
				ArrayFloat32:    []float32{11.1},
				ArrayFloat64:    []float64{12.2},
				ArrayBoolPtr:    []*bool{utils.ToPtr(true)},
				ArrayStringPtr:  []*string{utils.ToPtr("foo")},
				ArrayIntPtr:     []*int{utils.ToPtr(1)},
				ArrayInt8Ptr:    []*int8{utils.ToPtr[int8](2)},
				ArrayInt16Ptr:   []*int16{utils.ToPtr[int16](3)},
				ArrayInt32Ptr:   []*int32{utils.ToPtr[int32](4)},
				ArrayInt64Ptr:   []*int64{utils.ToPtr[int64](5)},
				ArrayUintPtr:    []*uint{utils.ToPtr[uint](6)},
				ArrayUint8Ptr:   []*uint8{utils.ToPtr[uint8](7)},
				ArrayUint16Ptr:  []*uint16{utils.ToPtr[uint16](8)},
				ArrayUint32Ptr:  []*uint32{utils.ToPtr[uint32](9)},
				ArrayUint64Ptr:  []*uint64{utils.ToPtr[uint64](10)},
				ArrayFloat32Ptr: []*float32{utils.ToPtr[float32](11.1)},
				ArrayFloat64Ptr: []*float64{utils.ToPtr[float64](12.2)},
				ArrayBigInt:     []scalar.BigInt{scalar.BigInt(13)},
				ArrayBigIntPtr:  []*scalar.BigInt{utils.ToPtr(scalar.BigInt(14))},
				ArrayJSON:       []any{map[string]any{"foo": "bar"}},
				ArrayJSONPtr:    []*any{utils.ToPtr(any(map[string]any{"foo": "baz"}))},
				ArrayRawJSON:    []json.RawMessage{json.RawMessage(`{"message":"raw_json"}`)},
				ArrayRawJSONPtr: []*json.RawMessage{utils.ToPtr(json.RawMessage(`{"message":"raw_json_ptr"}`))},
				ArrayTime:       []time.Time{time.Date(2024, 3, 5, 7, 0, 56, 0, time.UTC)},
				ArrayTimePtr:    []*time.Time{utils.ToPtr(time.Date(2024, 3, 5, 7, 0, 56, 0, time.UTC))},

				PtrArrayBool:       &[]bool{true},
				PtrArrayString:     &[]string{"foo"},
				PtrArrayInt:        &[]int{1},
				PtrArrayInt8:       &[]int8{2},
				PtrArrayInt16:      &[]int16{3},
				PtrArrayInt32:      &[]int32{4},
				PtrArrayInt64:      &[]int64{5},
				PtrArrayUint:       &[]uint{6},
				PtrArrayUint8:      &[]uint8{7},
				PtrArrayUint16:     &[]uint16{8},
				PtrArrayUint32:     &[]uint32{9},
				PtrArrayUint64:     &[]uint64{10},
				PtrArrayFloat32:    &[]float32{11.1},
				PtrArrayFloat64:    &[]float64{12.2},
				PtrArrayBoolPtr:    &[]*bool{utils.ToPtr(true)},
				PtrArrayStringPtr:  &[]*string{utils.ToPtr("foo")},
				PtrArrayIntPtr:     &[]*int{utils.ToPtr(1)},
				PtrArrayInt8Ptr:    &[]*int8{utils.ToPtr[int8](2)},
				PtrArrayInt16Ptr:   &[]*int16{utils.ToPtr[int16](3)},
				PtrArrayInt32Ptr:   &[]*int32{utils.ToPtr[int32](4)},
				PtrArrayInt64Ptr:   &[]*int64{utils.ToPtr[int64](5)},
				PtrArrayUintPtr:    &[]*uint{utils.ToPtr[uint](6)},
				PtrArrayUint8Ptr:   &[]*uint8{utils.ToPtr[uint8](7)},
				PtrArrayUint16Ptr:  &[]*uint16{utils.ToPtr[uint16](8)},
				PtrArrayUint32Ptr:  &[]*uint32{utils.ToPtr[uint32](9)},
				PtrArrayUint64Ptr:  &[]*uint64{utils.ToPtr[uint64](10)},
				PtrArrayFloat32Ptr: &[]*float32{utils.ToPtr[float32](11.1)},
				PtrArrayFloat64Ptr: &[]*float64{utils.ToPtr[float64](12.2)},
				PtrArrayBigInt:     &[]scalar.BigInt{scalar.BigInt(13)},
				PtrArrayBigIntPtr:  &[]*scalar.BigInt{utils.ToPtr(scalar.BigInt(14))},
				PtrArrayJSON:       &[]any{map[string]any{"foo": "bar"}},
				PtrArrayJSONPtr:    &[]*any{utils.ToPtr(any(map[string]any{"foo": "baz"}))},
				PtrArrayRawJSON:    &[]json.RawMessage{json.RawMessage(`{"message":"raw_json"}`)},
				PtrArrayRawJSONPtr: &[]*json.RawMessage{utils.ToPtr(json.RawMessage(`{"message":"raw_json_ptr"}`))},
				PtrArrayTime:       &[]time.Time{time.Date(2024, 3, 5, 7, 0, 56, 0, time.UTC)},
				PtrArrayTimePtr:    &[]*time.Time{utils.ToPtr(time.Date(2024, 3, 5, 7, 0, 56, 0, time.UTC))},
			},
		},
	}

	testServer := createTestServer(t).BuildTestServer()
	defer testServer.Close()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/query", testServer.URL), "application/json", bytes.NewReader([]byte(tc.body)))
			assert.NilError(t, err, "failed to request query")
			assert.Equal(t, tc.status, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if tc.errorMsg != "" {
				assert.NilError(t, err)
				assert.Check(t, strings.Contains(string(respBody), tc.errorMsg))
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("expected successful response, got error: %s", string(respBody))
			} else {
				log.Print("response: ", string(respBody))
				var results []struct {
					Rows []struct {
						Value arguments.GetTypesArguments `json:"__value"`
					} `json:"rows,omitempty" mapstructure:"rows,omitempty"`
				}
				assert.NilError(t, json.Unmarshal(respBody, &results), "failed to decode response")
				assert.Equal(t, 1, len(results))
				assert.Equal(t, 1, len(results[0].Rows))
				assert.DeepEqual(t, tc.response, results[0].Rows[0].Value, cmp.Exporter(func(t reflect.Type) bool { return true }))
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
									"text": { "type": "column", "column": "text", "fields": null },
									"error": { "type": "column", "column": "error", "fields": null }
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
				"text": "world",
				"error": {}
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
	}

	testServer := createTestServer(t).BuildTestServer()
	defer testServer.Close()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/query", testServer.URL), "application/json", bytes.NewReader([]byte(tc.body)))
			assert.NilError(t, err, "failed to request query")
			defer resp.Body.Close()
			assert.Equal(t, tc.status, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if tc.errorMsg != "" {
				assert.NilError(t, err)
				assert.Check(t, strings.Contains(string(respBody), tc.errorMsg))
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("expected successful response, got error: %s", string(respBody))
			} else {
				log.Print("response: ", string(respBody))
				var expected any
				assert.NilError(t, json.Unmarshal([]byte(tc.response), &expected))
				var results []struct {
					Rows []struct {
						Value any `json:"__value"`
					} `json:"rows,omitempty" mapstructure:"rows,omitempty"`
				}
				assert.NilError(t, json.Unmarshal(respBody, &results), "failed to decode response")
				assert.Equal(t, 1, len(results))
				assert.Equal(t, 1, len(results[0].Rows))
				assert.DeepEqual(t, expected, results[0].Rows[0].Value)
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
	defer testServer.Close()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/mutation", testServer.URL), "application/json", bytes.NewReader([]byte(tc.body)))
			assert.NilError(t, err, "failed to request mutation")
			defer resp.Body.Close()

			assert.Equal(t, tc.status, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			if tc.errorMsg != "" {
				assert.NilError(t, err)
				assert.Check(t, strings.Contains(string(respBody), tc.errorMsg))
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("expected successful response, got error: %s", string(respBody))
			} else {
				log.Print("response: ", string(respBody))
				var expected any
				assert.NilError(t, json.Unmarshal([]byte(tc.response), &expected))
				var results struct {
					OperationResults []struct {
						Result any `json:"result"`
					} `json:"operation_results"`
				}
				assert.NilError(t, json.Unmarshal(respBody, &results), "failed to decode response")
				assert.Equal(t, 1, len(results.OperationResults))
				assert.DeepEqual(t, expected, results.OperationResults[0].Result)
			}
		})
	}
}

func TestQuerySync(t *testing.T) {
	latency := testQueryLatency(t)
	assert.Assert(t, latency >= (3*time.Second))
}

func TestQueryAsync(t *testing.T) {
	t.Setenv("QUERY_CONCURRENCY_LIMIT", "2")
	latency := testQueryLatency(t)
	assert.Assert(t, latency < (3*time.Second))
}

func TestQueryAsyncByName(t *testing.T) {
	t.Setenv("QUERY_CONCURRENCY", "getArticles=2")
	latency := testQueryLatency(t)
	assert.Assert(t, latency < (3*time.Second))
}

func testQueryLatency(t *testing.T) time.Duration {
	testServer := createTestServer(t).BuildTestServer()
	defer testServer.Close()
	rawRequestBody, err := os.ReadFile("testdata/query/getArticles/request.json")
	assert.NilError(t, err)
	rawExpectedBody, err := os.ReadFile("testdata/query/getArticles/expected.json")
	assert.NilError(t, err)

	var expected, respBody schema.QueryResponse
	assert.NilError(t, json.Unmarshal(rawExpectedBody, &expected))

	start := time.Now()
	resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/query", testServer.URL), "application/json", bytes.NewReader(rawRequestBody))
	latency := time.Since(start)
	assert.NilError(t, err, "failed to request query")
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.NilError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	assert.DeepEqual(t, expected, respBody)

	return latency
}

func TestMutationSync(t *testing.T) {
	latency := testMutationLatency(t)
	assert.Assert(t, latency >= (1500*time.Millisecond))
}

func TestMutationAsync2(t *testing.T) {
	t.Setenv("MUTATION_CONCURRENCY_LIMIT", "2")
	latency := testMutationLatency(t)
	assert.Assert(t, latency < (2000*time.Millisecond))
}

func TestMutationAsync3(t *testing.T) {
	t.Setenv("MUTATION_CONCURRENCY_LIMIT", "3")
	latency := testMutationLatency(t)
	assert.Assert(t, latency < (1500*time.Millisecond))
}

func testMutationLatency(t *testing.T) time.Duration {
	testServer := createTestServer(t).BuildTestServer()
	defer testServer.Close()
	rawRequestBody, err := os.ReadFile("testdata/mutation/create_article/request.json")
	assert.NilError(t, err)
	rawExpectedBody, err := os.ReadFile("testdata/mutation/create_article/expected.json")
	assert.NilError(t, err)

	var expected, respBody schema.MutationResponse
	assert.NilError(t, json.Unmarshal(rawExpectedBody, &expected))

	start := time.Now()
	resp, err := http.DefaultClient.Post(fmt.Sprintf("%s/mutation", testServer.URL), "application/json", bytes.NewReader(rawRequestBody))
	latency := time.Since(start)
	assert.NilError(t, err, "failed to request mutation")
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.NilError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	assert.DeepEqual(t, expected, respBody)

	return latency
}
