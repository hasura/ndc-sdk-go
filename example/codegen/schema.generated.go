package main

import (
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

// GetConnectorSchema gets the generated connector schema
func GetConnectorSchema() schema.SchemaResponse {
	return schema.SchemaResponse{
		Collections: []schema.CollectionInfo{},
		ObjectTypes: schema.SchemaResponseObjectTypes{
			"Author": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Type: schema.NewNamedType("String").Encode(),
					},
					"created_at": schema.ObjectField{
						Type: schema.NewNamedType("TimestampTZ").Encode(),
					},
					"tags": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("String")).Encode(),
					},
					"author": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Author")).Encode(),
					},
				},
			},
			"CreateArticleArgumentsAuthor": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Type: schema.NewNamedType("UUID").Encode(),
					},
					"created_at": schema.ObjectField{
						Type: schema.NewNamedType("TimestampTZ").Encode(),
					},
				},
			},
			"CreateArticleResult": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"authors": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Author")).Encode(),
					},
				},
			},
			"CreateAuthorResult": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"name": schema.ObjectField{
						Type: schema.NewNamedType("String").Encode(),
					},
					"created_at": schema.ObjectField{
						Type: schema.NewNamedType("TimestampTZ").Encode(),
					},
				},
			},
			"GetArticlesResult": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Type: schema.NewNamedType("String").Encode(),
					},
					"Name": schema.ObjectField{
						Type: schema.NewNamedType("String").Encode(),
					},
				},
			},
			"GetTypesArguments": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"ArrayBool": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Boolean")).Encode(),
					},
					"PtrArrayInt8Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8")))).Encode(),
					},
					"PtrArrayInt32Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"UUID": schema.ObjectField{
						Type: schema.NewNamedType("UUID").Encode(),
					},
					"UintPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"ArrayInt64": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int64")).Encode(),
					},
					"ArrayBigInt": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("BigInt")).Encode(),
					},
					"PtrArrayUint8": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int8"))).Encode(),
					},
					"Int32": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"Uint32Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"Object": schema.ObjectField{
						Type: schema.NewNamedType("GetTypesArgumentsObject").Encode(),
					},
					"ArrayMapPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("JSON"))).Encode(),
					},
					"Int": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"BigInt": schema.ObjectField{
						Type: schema.NewNamedType("BigInt").Encode(),
					},
					"JSON": schema.ObjectField{
						Type: schema.NewNamedType("JSON").Encode(),
					},
					"ArrayUint": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"ArrayJSONPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("JSON"))).Encode(),
					},
					"PtrArrayUint16": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int16"))).Encode(),
					},
					"PtrArrayUint32Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"ObjectPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("GetTypesArgumentsObjectPtr")).Encode(),
					},
					"Float64Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Float64")).Encode(),
					},
					"ArrayInt32": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"ArrayUUIDPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("UUID"))).Encode(),
					},
					"PtrArrayInt32": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"ArrayMap": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("JSON")).Encode(),
					},
					"PtrArrayTime": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("TimestampTZ"))).Encode(),
					},
					"ArrayObject": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("GetTypesArgumentsArrayObject")).Encode(),
					},
					"Int8Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int8")).Encode(),
					},
					"ArrayString": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("String")).Encode(),
					},
					"ArrayFloat64": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Float64")).Encode(),
					},
					"ArrayTimePtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("TimestampTZ"))).Encode(),
					},
					"PtrArrayUUIDPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("UUID")))).Encode(),
					},
					"String": schema.ObjectField{
						Type: schema.NewNamedType("String").Encode(),
					},
					"Int16": schema.ObjectField{
						Type: schema.NewNamedType("Int16").Encode(),
					},
					"PtrArrayFloat64": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Float64"))).Encode(),
					},
					"Bytes": schema.ObjectField{
						Type: schema.NewNamedType("Bytes").Encode(),
					},
					"ArrayInt8Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8"))).Encode(),
					},
					"ArrayInt64Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64"))).Encode(),
					},
					"PtrArrayStringPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("String")))).Encode(),
					},
					"PtrArrayFloat32Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float32")))).Encode(),
					},
					"NamedObjectPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Author")).Encode(),
					},
					"PtrArrayString": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("String"))).Encode(),
					},
					"PtrArrayUint8Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8")))).Encode(),
					},
					"StringPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("String")).Encode(),
					},
					"Uint64Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int64")).Encode(),
					},
					"ArrayInt": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"ArrayInt32Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"ArrayUint32Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayJSONPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("JSON")))).Encode(),
					},
					"BoolPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Boolean")).Encode(),
					},
					"Int16Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int16")).Encode(),
					},
					"TextPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("String")).Encode(),
					},
					"ArrayInt16": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int16")).Encode(),
					},
					"PtrArrayUint16Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16")))).Encode(),
					},
					"CustomScalarPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("CommentString")).Encode(),
					},
					"PtrArrayUint": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayInt64Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64")))).Encode(),
					},
					"ArrayObjectPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("GetTypesArgumentsArrayObjectPtr"))).Encode(),
					},
					"Uint": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"Time": schema.ObjectField{
						Type: schema.NewNamedType("TimestampTZ").Encode(),
					},
					"ArrayUint8Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8"))).Encode(),
					},
					"ArrayJSON": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("JSON")).Encode(),
					},
					"PtrArrayUintPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"Float32Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Float32")).Encode(),
					},
					"ArrayUint64Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64"))).Encode(),
					},
					"PtrArrayBoolPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Boolean")))).Encode(),
					},
					"IntPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"PtrArrayRawJSONPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("RawJSON")))).Encode(),
					},
					"PtrArrayUint32": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayUint64Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64")))).Encode(),
					},
					"PtrArrayBigIntPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("BigInt")))).Encode(),
					},
					"Int8": schema.ObjectField{
						Type: schema.NewNamedType("Int8").Encode(),
					},
					"Uint8": schema.ObjectField{
						Type: schema.NewNamedType("Int8").Encode(),
					},
					"Date": schema.ObjectField{
						Type: schema.NewNamedType("Date").Encode(),
					},
					"DatePtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Date")).Encode(),
					},
					"PtrArrayInt16": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int16"))).Encode(),
					},
					"ArrayUUID": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("UUID")).Encode(),
					},
					"Map": schema.ObjectField{
						Type: schema.NewNamedType("JSON").Encode(),
					},
					"ArrayRawJSONPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("RawJSON"))).Encode(),
					},
					"PtrArrayRawJSON": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("RawJSON"))).Encode(),
					},
					"Float32": schema.ObjectField{
						Type: schema.NewNamedType("Float32").Encode(),
					},
					"ArrayInt8": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int8")).Encode(),
					},
					"ArrayBoolPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Boolean"))).Encode(),
					},
					"ArrayIntPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"ArrayRawJSON": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("RawJSON")).Encode(),
					},
					"ArrayUint8": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int8")).Encode(),
					},
					"ArrayBigIntPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("BigInt"))).Encode(),
					},
					"PtrArrayJSON": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("JSON"))).Encode(),
					},
					"PtrArrayTimePtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("TimestampTZ")))).Encode(),
					},
					"PtrArrayIntPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"EnumPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("SomeEnum")).Encode(),
					},
					"PtrArrayBool": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Boolean"))).Encode(),
					},
					"UUIDPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("UUID")).Encode(),
					},
					"Int64Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int64")).Encode(),
					},
					"Uint16Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int16")).Encode(),
					},
					"BigIntPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("BigInt")).Encode(),
					},
					"ArrayUintPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"TimePtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("TimestampTZ")).Encode(),
					},
					"PtrArrayInt16Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16")))).Encode(),
					},
					"Uint32": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"ArrayUint64": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int64")).Encode(),
					},
					"PtrArrayInt8": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int8"))).Encode(),
					},
					"PtrArrayFloat32": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Float32"))).Encode(),
					},
					"RawJSONPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("RawJSON")).Encode(),
					},
					"Uint16": schema.ObjectField{
						Type: schema.NewNamedType("Int16").Encode(),
					},
					"PtrArrayBigInt": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("BigInt"))).Encode(),
					},
					"NamedArrayPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Author"))).Encode(),
					},
					"MapPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
					},
					"PtrArrayFloat64Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float64")))).Encode(),
					},
					"RawJSON": schema.ObjectField{
						Type: schema.NewNamedType("RawJSON").Encode(),
					},
					"Float64": schema.ObjectField{
						Type: schema.NewNamedType("Float64").Encode(),
					},
					"CustomScalar": schema.ObjectField{
						Type: schema.NewNamedType("CommentString").Encode(),
					},
					"ArrayUint16": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int16")).Encode(),
					},
					"ArrayUint32": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"ArrayFloat32Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float32"))).Encode(),
					},
					"NamedObject": schema.ObjectField{
						Type: schema.NewNamedType("Author").Encode(),
					},
					"Bool": schema.ObjectField{
						Type: schema.NewNamedType("Boolean").Encode(),
					},
					"Uint8Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int8")).Encode(),
					},
					"ArrayStringPtr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("String"))).Encode(),
					},
					"ArrayInt16Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16"))).Encode(),
					},
					"ArrayTime": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("TimestampTZ")).Encode(),
					},
					"Text": schema.ObjectField{
						Type: schema.NewNamedType("String").Encode(),
					},
					"Enum": schema.ObjectField{
						Type: schema.NewNamedType("SomeEnum").Encode(),
					},
					"URL": schema.ObjectField{
						Type: schema.NewNamedType("URL").Encode(),
					},
					"BytesPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Bytes")).Encode(),
					},
					"Uint64": schema.ObjectField{
						Type: schema.NewNamedType("Int64").Encode(),
					},
					"Int32Ptr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"ArrayFloat32": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Float32")).Encode(),
					},
					"ArrayFloat64Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float64"))).Encode(),
					},
					"PtrArrayUint64": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int64"))).Encode(),
					},
					"Int64": schema.ObjectField{
						Type: schema.NewNamedType("Int64").Encode(),
					},
					"PtrArrayInt": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayUUID": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("UUID"))).Encode(),
					},
					"ArrayUint16Ptr": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16"))).Encode(),
					},
					"PtrArrayInt64": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int64"))).Encode(),
					},
					"NamedArray": schema.ObjectField{
						Type: schema.NewArrayType(schema.NewNamedType("Author")).Encode(),
					},
					"JSONPtr": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
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
					"id": schema.ObjectField{
						Type: schema.NewNamedType("UUID").Encode(),
					},
					"created_at": schema.ObjectField{
						Type: schema.NewNamedType("TimestampTZ").Encode(),
					},
				},
			},
			"GetTypesArgumentsObjectPtr": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"Long": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"Lat": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
				},
			},
			"HelloResult": schema.ObjectType{
				Fields: schema.ObjectTypeFields{
					"num": schema.ObjectField{
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"text": schema.ObjectField{
						Type: schema.NewNamedType("String").Encode(),
					},
					"foo": schema.ObjectField{
						Type: schema.NewNamedType("Foo").Encode(),
					},
					"error": schema.ObjectField{
						Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
					},
					"id": schema.ObjectField{
						Type: schema.NewNamedType("UUID").Encode(),
					},
				},
			},
		},
		Functions: []schema.FunctionInfo{
			{
				Name:        "getBool",
				Description: utils.ToPtr("return an scalar boolean"),
				ResultType:  schema.NewNamedType("Boolean").Encode(),
				Arguments:   map[string]schema.ArgumentInfo{},
			},
			{
				Name:       "getTypes",
				ResultType: schema.NewNullableType(schema.NewNamedType("GetTypesArguments")).Encode(),
				Arguments: map[string]schema.ArgumentInfo{
					"IntPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"PtrArrayInt32Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"PtrArrayBoolPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Boolean")))).Encode(),
					},
					"ArrayRawJSON": {
						Type: schema.NewArrayType(schema.NewNamedType("RawJSON")).Encode(),
					},
					"NamedArray": {
						Type: schema.NewArrayType(schema.NewNamedType("Author")).Encode(),
					},
					"PtrArrayInt32": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayInt64": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int64"))).Encode(),
					},
					"ArrayUint8Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8"))).Encode(),
					},
					"ArrayUint8": {
						Type: schema.NewArrayType(schema.NewNamedType("Int8")).Encode(),
					},
					"NamedObjectPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("Author")).Encode(),
					},
					"PtrArrayInt": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayInt64Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64")))).Encode(),
					},
					"PtrArrayTime": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("TimestampTZ"))).Encode(),
					},
					"Text": {
						Type: schema.NewNamedType("String").Encode(),
					},
					"Int64Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int64")).Encode(),
					},
					"PtrArrayUint32": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"Int32Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"StringPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("String")).Encode(),
					},
					"ArrayBool": {
						Type: schema.NewArrayType(schema.NewNamedType("Boolean")).Encode(),
					},
					"ArrayTimePtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("TimestampTZ"))).Encode(),
					},
					"ArrayObjectPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("GetTypesArgumentsArrayObjectPtr"))).Encode(),
					},
					"ArrayIntPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayJSON": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("JSON"))).Encode(),
					},
					"Int8Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int8")).Encode(),
					},
					"ArrayTime": {
						Type: schema.NewArrayType(schema.NewNamedType("TimestampTZ")).Encode(),
					},
					"Uint64": {
						Type: schema.NewNamedType("Int64").Encode(),
					},
					"BigInt": {
						Type: schema.NewNamedType("BigInt").Encode(),
					},
					"PtrArrayUint64Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64")))).Encode(),
					},
					"ArrayInt8Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8"))).Encode(),
					},
					"Uint32": {
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"ArrayStringPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("String"))).Encode(),
					},
					"ArrayUint16": {
						Type: schema.NewArrayType(schema.NewNamedType("Int16")).Encode(),
					},
					"ArrayInt16Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16"))).Encode(),
					},
					"ArrayString": {
						Type: schema.NewArrayType(schema.NewNamedType("String")).Encode(),
					},
					"URL": {
						Type: schema.NewNamedType("URL").Encode(),
					},
					"PtrArrayUint16Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16")))).Encode(),
					},
					"PtrArrayUint": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int32"))).Encode(),
					},
					"PtrArrayBigInt": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("BigInt"))).Encode(),
					},
					"Uint64Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int64")).Encode(),
					},
					"ArrayUint64Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64"))).Encode(),
					},
					"ArrayUUIDPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("UUID"))).Encode(),
					},
					"ArrayFloat32Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float32"))).Encode(),
					},
					"Int64": {
						Type: schema.NewNamedType("Int64").Encode(),
					},
					"ArrayUint32": {
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"ArrayUUID": {
						Type: schema.NewArrayType(schema.NewNamedType("UUID")).Encode(),
					},
					"PtrArrayInt8": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int8"))).Encode(),
					},
					"Bytes": {
						Type: schema.NewNamedType("Bytes").Encode(),
					},
					"ArrayInt": {
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"PtrArrayString": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("String"))).Encode(),
					},
					"ArrayMap": {
						Type: schema.NewArrayType(schema.NewNamedType("JSON")).Encode(),
					},
					"ArrayInt8": {
						Type: schema.NewArrayType(schema.NewNamedType("Int8")).Encode(),
					},
					"ArrayBoolPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Boolean"))).Encode(),
					},
					"Uint": {
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"Date": {
						Type: schema.NewNamedType("Date").Encode(),
					},
					"Float32": {
						Type: schema.NewNamedType("Float32").Encode(),
					},
					"Enum": {
						Type: schema.NewNamedType("SomeEnum").Encode(),
					},
					"PtrArrayFloat64": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Float64"))).Encode(),
					},
					"PtrArrayUUIDPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("UUID")))).Encode(),
					},
					"ArrayRawJSONPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("RawJSON"))).Encode(),
					},
					"ArrayObject": {
						Type: schema.NewArrayType(schema.NewNamedType("GetTypesArgumentsArrayObject")).Encode(),
					},
					"ArrayUint64": {
						Type: schema.NewArrayType(schema.NewNamedType("Int64")).Encode(),
					},
					"Int16": {
						Type: schema.NewNamedType("Int16").Encode(),
					},
					"PtrArrayUint32Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"CustomScalar": {
						Type: schema.NewNamedType("CommentString").Encode(),
					},
					"Uint16Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int16")).Encode(),
					},
					"ArrayBigIntPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("BigInt"))).Encode(),
					},
					"TimePtr": {
						Type: schema.NewNullableType(schema.NewNamedType("TimestampTZ")).Encode(),
					},
					"ArrayInt64Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int64"))).Encode(),
					},
					"ArrayUintPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"ArrayUint16Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16"))).Encode(),
					},
					"PtrArrayInt8Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8")))).Encode(),
					},
					"ArrayMapPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("JSON"))).Encode(),
					},
					"Uint16": {
						Type: schema.NewNamedType("Int16").Encode(),
					},
					"BigIntPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("BigInt")).Encode(),
					},
					"RawJSONPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("RawJSON")).Encode(),
					},
					"PtrArrayRawJSONPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("RawJSON")))).Encode(),
					},
					"ArrayJSONPtr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("JSON"))).Encode(),
					},
					"PtrArrayRawJSON": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("RawJSON"))).Encode(),
					},
					"Int": {
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"UintPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"DatePtr": {
						Type: schema.NewNullableType(schema.NewNamedType("Date")).Encode(),
					},
					"Object": {
						Type: schema.NewNamedType("GetTypesArgumentsObject").Encode(),
					},
					"PtrArrayUint8": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int8"))).Encode(),
					},
					"ArrayUint": {
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"ObjectPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("GetTypesArgumentsObjectPtr")).Encode(),
					},
					"PtrArrayStringPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("String")))).Encode(),
					},
					"MapPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
					},
					"PtrArrayBigIntPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("BigInt")))).Encode(),
					},
					"ArrayFloat64": {
						Type: schema.NewArrayType(schema.NewNamedType("Float64")).Encode(),
					},
					"Float64": {
						Type: schema.NewNamedType("Float64").Encode(),
					},
					"PtrArrayTimePtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("TimestampTZ")))).Encode(),
					},
					"PtrArrayUUID": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("UUID"))).Encode(),
					},
					"PtrArrayFloat64Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float64")))).Encode(),
					},
					"UUIDPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("UUID")).Encode(),
					},
					"BytesPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("Bytes")).Encode(),
					},
					"ArrayJSON": {
						Type: schema.NewArrayType(schema.NewNamedType("JSON")).Encode(),
					},
					"PtrArrayUint8Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int8")))).Encode(),
					},
					"PtrArrayFloat32Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float32")))).Encode(),
					},
					"ArrayInt32": {
						Type: schema.NewArrayType(schema.NewNamedType("Int32")).Encode(),
					},
					"PtrArrayInt16Ptr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int16")))).Encode(),
					},
					"ArrayUint32Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"Map": {
						Type: schema.NewNamedType("JSON").Encode(),
					},
					"PtrArrayBool": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Boolean"))).Encode(),
					},
					"Time": {
						Type: schema.NewNamedType("TimestampTZ").Encode(),
					},
					"ArrayBigInt": {
						Type: schema.NewArrayType(schema.NewNamedType("BigInt")).Encode(),
					},
					"PtrArrayIntPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"ArrayInt32Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32"))).Encode(),
					},
					"Float32Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Float32")).Encode(),
					},
					"String": {
						Type: schema.NewNamedType("String").Encode(),
					},
					"PtrArrayUint64": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int64"))).Encode(),
					},
					"JSON": {
						Type: schema.NewNamedType("JSON").Encode(),
					},
					"ArrayInt16": {
						Type: schema.NewArrayType(schema.NewNamedType("Int16")).Encode(),
					},
					"Uint32Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int32")).Encode(),
					},
					"Int32": {
						Type: schema.NewNamedType("Int32").Encode(),
					},
					"Int16Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int16")).Encode(),
					},
					"NamedObject": {
						Type: schema.NewNamedType("Author").Encode(),
					},
					"EnumPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("SomeEnum")).Encode(),
					},
					"JSONPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("JSON")).Encode(),
					},
					"Uint8Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Int8")).Encode(),
					},
					"Bool": {
						Type: schema.NewNamedType("Boolean").Encode(),
					},
					"BoolPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("Boolean")).Encode(),
					},
					"UUID": {
						Type: schema.NewNamedType("UUID").Encode(),
					},
					"NamedArrayPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Author"))).Encode(),
					},
					"Uint8": {
						Type: schema.NewNamedType("Int8").Encode(),
					},
					"RawJSON": {
						Type: schema.NewNamedType("RawJSON").Encode(),
					},
					"CustomScalarPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("CommentString")).Encode(),
					},
					"Float64Ptr": {
						Type: schema.NewNullableType(schema.NewNamedType("Float64")).Encode(),
					},
					"ArrayInt64": {
						Type: schema.NewArrayType(schema.NewNamedType("Int64")).Encode(),
					},
					"PtrArrayInt16": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int16"))).Encode(),
					},
					"Int8": {
						Type: schema.NewNamedType("Int8").Encode(),
					},
					"PtrArrayJSONPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("JSON")))).Encode(),
					},
					"PtrArrayUint16": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Int16"))).Encode(),
					},
					"PtrArrayFloat32": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNamedType("Float32"))).Encode(),
					},
					"ArrayFloat64Ptr": {
						Type: schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Float64"))).Encode(),
					},
					"PtrArrayUintPtr": {
						Type: schema.NewNullableType(schema.NewArrayType(schema.NewNullableType(schema.NewNamedType("Int32")))).Encode(),
					},
					"ArrayFloat32": {
						Type: schema.NewArrayType(schema.NewNamedType("Float32")).Encode(),
					},
					"TextPtr": {
						Type: schema.NewNullableType(schema.NewNamedType("String")).Encode(),
					},
				},
			},
			{
				Name:        "hello",
				Description: utils.ToPtr("sends a hello message"),
				ResultType:  schema.NewNullableType(schema.NewNamedType("HelloResult")).Encode(),
				Arguments:   map[string]schema.ArgumentInfo{},
			},
			{
				Name:        "getArticles",
				Description: utils.ToPtr("GetArticles"),
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
				Description: utils.ToPtr("CreateArticle"),
				ResultType:  schema.NewNullableType(schema.NewNamedType("CreateArticleResult")).Encode(),
				Arguments: map[string]schema.ArgumentInfo{
					"author": {
						Type: schema.NewNamedType("CreateArticleArgumentsAuthor").Encode(),
					},
				},
			},
			{
				Name:        "increase",
				Description: utils.ToPtr("Increase"),
				ResultType:  schema.NewNamedType("Int32").Encode(),
				Arguments:   map[string]schema.ArgumentInfo{},
			},
			{
				Name:        "createAuthor",
				Description: utils.ToPtr("creates an author"),
				ResultType:  schema.NewNullableType(schema.NewNamedType("CreateAuthorResult")).Encode(),
				Arguments: map[string]schema.ArgumentInfo{
					"name": {
						Type: schema.NewNamedType("String").Encode(),
					},
				},
			},
			{
				Name:        "createAuthors",
				Description: utils.ToPtr("creates a list of authors"),
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
			"Boolean": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationBoolean().Encode(),
			},
			"Bytes": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationBytes().Encode(),
			},
			"CommentString": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationString().Encode(),
			},
			"Date": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationDate().Encode(),
			},
			"Float32": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationFloat32().Encode(),
			},
			"Float64": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationFloat64().Encode(),
			},
			"Foo": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationJSON().Encode(),
			},
			"Int16": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationInt16().Encode(),
			},
			"Int32": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationInt32().Encode(),
			},
			"Int64": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationInt64().Encode(),
			},
			"Int8": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationInt8().Encode(),
			},
			"JSON": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationJSON().Encode(),
			},
			"RawJSON": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationJSON().Encode(),
			},
			"SomeEnum": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationEnum([]string{"foo", "bar"}).Encode(),
			},
			"String": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationString().Encode(),
			},
			"TimestampTZ": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationTimestampTZ().Encode(),
			},
			"URL": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationString().Encode(),
			},
			"UUID": schema.ScalarType{
				AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
				Representation:      schema.NewTypeRepresentationUUID().Encode(),
			},
		},
	}
}
