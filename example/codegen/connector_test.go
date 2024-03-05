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
	"github.com/hasura/ndc-sdk-go/schema"
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
						"CustomScalar": {
							"type": "literal",
							"value": "a comment"
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
								"value": "hello"
						},
						"IntPtr": {
								"type": "literal",
								"value": 1
						},
						"Int8Ptr": {
								"type": "literal",
								"value": 2
						},
						"Int16Ptr": {
								"type": "literal",
								"value": 3
						},
						"Int32Ptr": {
							"type": "literal",
							"value": 4
						},
						"Int64Ptr": {
							"type": "literal",
							"value": 5
						},
						"UintPtr": {
							"type": "literal",
							"value": 6
						},
						"Uint8Ptr": {
							"type": "literal",
							"value": 7
						},
						"Uint16Ptr": {
							"type": "literal",
							"value": 8
						},
						"Uint32Ptr": {
							"type": "literal",
							"value": 9
						},
						"Uint64Ptr": {
							"type": "literal",
							"value": 10
						},
						"Float32Ptr": {
							"type": "literal",
							"value": 1.1
						},
						"Float64Ptr": {
							"type": "literal",
							"value": 2.2
						},
						"TimePtr": {
							"type": "literal",
							"value": "2024-03-05T07:00:56Z"
						},
						"DurationPtr": {
							"type": "literal",
							"value": "10s"
						},
						"CustomScalarPtr": {
							"type": "literal",
							"value": "a comment"
						},
						"Object": {
							"type": "literal",
							"value": null
						},
						"ObjectPtr": {
							"type": "literal",
							"value": null
						},
						"ArrayObject": {
							"type": "literal",
							"value": null
						},
						"NamedObject": {
							"type": "literal",
							"value": null
						},
						"NamedObjectPtr": {
							"type": "literal",
							"value": null
						},
						"NamedArray": {
							"type": "literal",
							"value": null
						}
				},
				"query": {
						"fields": {
								"UUID": {
										"type": "column",
										"column": "UUID"
								},
								"Bool": {
										"type": "column",
										"column": "Bool"
								},
								"String": {
										"type": "column",
										"column": "String"
								},
								"Int": {
										"type": "column",
										"column": "Int"
								},
								"Int8": {
										"type": "column",
										"column": "Int8"
								},
								"Int16": {
										"type": "column",
										"column": "Int16"
								},
								"Int32": {
										"type": "column",
										"column": "Int32"
								},
								"Int64": {
										"type": "column",
										"column": "Int64"
								},
								"Uint": {
										"type": "column",
										"column": "Uint"
								},
								"Uint8": {
										"type": "column",
										"column": "Uint8"
								},
								"Uint16": {
										"type": "column",
										"column": "Uint16"
								},
								"Uint32": {
										"type": "column",
										"column": "Uint32"
								},
								"Uint64": {
										"type": "column",
										"column": "Uint64"
								},
								"Float32": {
										"type": "column",
										"column": "Float32"
								},
								"Float64": {
										"type": "column",
										"column": "Float64"
								},
								"Time": {
										"type": "column",
										"column": "Time"
								},
								"Duration": {
										"type": "column",
										"column": "Duration"
								},
								"CustomScalar": {
									"type": "column",
									"column": "CustomScalar"
								},
								"UUIDPtr": {
										"type": "column",
										"column": "UUIDPtr"
								},
								"BoolPtr": {
										"type": "column",
										"column": "BoolPtr"
								},
								"StringPtr": {
										"type": "column",
										"column": "StringPtr"
								},
								"IntPtr": {
										"type": "column",
										"column": "IntPtr"
								},
								"Int8Ptr": {
										"type": "column",
										"column": "Int8Ptr"
								},
								"Int16Ptr": {
										"type": "column",
										"column": "Int16Ptr"
								},
								"Int32Ptr": {
										"type": "column",
										"column": "Int32Ptr"
								},
								"Int64Ptr": {
										"type": "column",
										"column": "Int64Ptr"
								},
								"UintPtr": {
										"type": "column",
										"column": "UintPtr"
								},
								"Uint8Ptr": {
										"type": "column",
										"column": "Uint8Ptr"
								},
								"Uint16Ptr": {
										"type": "column",
										"column": "Uint16Ptr"
								},
								"Uint32Ptr": {
										"type": "column",
										"column": "Uint32Ptr"
								},
								"Uint64Ptr": {
										"type": "column",
										"column": "Uint64Ptr"
								},
								"Float32Ptr": {
										"type": "column",
										"column": "Float32Ptr"
								},
								"Float64Ptr": {
										"type": "column",
										"column": "Float64Ptr"
								},
								"TimePtr": {
										"type": "column",
										"column": "TimePtr"
								},
								"DurationPtr": {
										"type": "column",
										"column": "DurationPtr"
								},
								"CustomScalarPtr": {
									"type": "column",
									"column": "CustomScalarPtr"
								},
								"Object": {
									"type": "column",
									"column": "Object"
								},
								"ObjectPtr": {
									"type": "column",
									"column": "ObjectPtr"
								},
								"ArrayObject": {
									"type": "column",
									"column": "ArrayObject"
								},
								"NamedObject": {
									"type": "column",
									"column": "NamedObject"
								},
								"NamedObjectPtr": {
									"type": "column",
									"column": "NamedObjectPtr"
								},
								"NamedArray": {
									"type": "column",
									"column": "NamedArray"
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
				Time:            time.Date(2023, 3, 5, 7, 0, 56, 0, time.UTC),
				Duration:        10 * time.Second,
				CustomScalar:    commentText,
				UUIDPtr:         schema.ToPtr(uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5b")),
				BoolPtr:         schema.ToPtr(true),
				IntPtr:          schema.ToPtr(11),
				Int8Ptr:         schema.ToPtr(int8(12)),
				Int16Ptr:        schema.ToPtr(int16(13)),
				Int32Ptr:        schema.ToPtr(int32(14)),
				Int64Ptr:        schema.ToPtr(int64(15)),
				UintPtr:         schema.ToPtr(uint(16)),
				Uint8Ptr:        schema.ToPtr(uint8(17)),
				Uint16Ptr:       schema.ToPtr(uint16(18)),
				Uint32Ptr:       schema.ToPtr(uint32(19)),
				Uint64Ptr:       schema.ToPtr(uint64(20)),
				Float32Ptr:      schema.ToPtr(float32(3.3)),
				Float64Ptr:      schema.ToPtr(float64(4.4)),
				TimePtr:         schema.ToPtr(time.Date(2023, 3, 5, 7, 0, 0, 0, time.UTC)),
				DurationPtr:     schema.ToPtr(time.Minute),
				CustomScalarPtr: &commentTextPtr,
				Object: struct {
					ID        uuid.UUID `json:"id"`
					CreatedAt time.Time `json:"created_at"`
				}{
					ID:        uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5c"),
					CreatedAt: time.Date(2023, 3, 5, 6, 0, 0, 0, time.UTC),
				},
				ObjectPtr: &struct {
					Long int
					Lat  int
				}{
					Long: 1,
					Lat:  2,
				},
				ArrayObject: []struct {
					Content string "json:\"content\""
				}{
					{
						Content: "a content",
					},
				},
				NamedObject: functions.Author{
					ID:        "1",
					Duration:  10 * time.Minute,
					CreatedAt: time.Date(2023, 3, 5, 5, 0, 0, 0, time.UTC),
				},
				NamedObjectPtr: &functions.Author{
					ID:        "2",
					Duration:  11 * time.Minute,
					CreatedAt: time.Date(2023, 3, 5, 4, 0, 0, 0, time.UTC),
				},
				NamedArray: []functions.Author{
					{
						ID:        "3",
						Duration:  12 * time.Minute,
						CreatedAt: time.Date(2023, 3, 5, 3, 0, 0, 0, time.UTC),
					},
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
