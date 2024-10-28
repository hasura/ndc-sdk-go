package arguments

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/scalar"
)

type GetTypesArguments struct {
	UUID         uuid.UUID
	Bool         bool
	String       string
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Uint         uint
	Uint8        uint8
	Uint16       uint16
	Uint32       uint32
	Uint64       uint64
	Float32      float32
	Float64      float64
	Time         time.Time
	Text         types.Text
	CustomScalar types.CommentText
	Enum         types.SomeEnum
	BigInt       scalar.BigInt
	Date         scalar.Date
	URL          scalar.URL
	Duration     scalar.Duration

	UUIDPtr         *uuid.UUID
	BoolPtr         *bool
	StringPtr       *string
	IntPtr          *int
	Int8Ptr         *int8
	Int16Ptr        *int16
	Int32Ptr        *int32
	Int64Ptr        *int64
	UintPtr         *uint
	Uint8Ptr        *uint8
	Uint16Ptr       *uint16
	Uint32Ptr       *uint32
	Uint64Ptr       *uint64
	Float32Ptr      *float32
	Float64Ptr      *float64
	TimePtr         *time.Time
	TextPtr         *types.Text
	CustomScalarPtr *types.CommentText
	EnumPtr         *types.SomeEnum
	BigIntPtr       *scalar.BigInt
	DatePtr         *scalar.Date

	ArrayBool       []bool
	ArrayString     []string
	ArrayInt        []int
	ArrayInt8       []int8
	ArrayInt16      []int16
	ArrayInt32      []int32
	ArrayInt64      []int64
	ArrayUint       []uint
	ArrayUint8      []uint8
	ArrayUint16     []uint16
	ArrayUint32     []uint32
	ArrayUint64     []uint64
	ArrayFloat32    []float32
	ArrayFloat64    []float64
	ArrayUUID       []uuid.UUID
	ArrayBoolPtr    []*bool
	ArrayStringPtr  []*string
	ArrayIntPtr     []*int
	ArrayInt8Ptr    []*int8
	ArrayInt16Ptr   []*int16
	ArrayInt32Ptr   []*int32
	ArrayInt64Ptr   []*int64
	ArrayUintPtr    []*uint
	ArrayUint8Ptr   []*uint8
	ArrayUint16Ptr  []*uint16
	ArrayUint32Ptr  []*uint32
	ArrayUint64Ptr  []*uint64
	ArrayFloat32Ptr []*float32
	ArrayFloat64Ptr []*float64
	ArrayUUIDPtr    []*uuid.UUID
	ArrayJSON       []any
	ArrayJSONPtr    []*interface{}
	ArrayRawJSON    []json.RawMessage
	ArrayRawJSONPtr []*json.RawMessage
	ArrayBigInt     []scalar.BigInt
	ArrayBigIntPtr  []*scalar.BigInt
	ArrayTime       []time.Time
	ArrayTimePtr    []*time.Time

	PtrArrayBool       *[]bool
	PtrArrayString     *[]string
	PtrArrayInt        *[]int
	PtrArrayInt8       *[]int8
	PtrArrayInt16      *[]int16
	PtrArrayInt32      *[]int32
	PtrArrayInt64      *[]int64
	PtrArrayUint       *[]uint
	PtrArrayUint8      *[]uint8
	PtrArrayUint16     *[]uint16
	PtrArrayUint32     *[]uint32
	PtrArrayUint64     *[]uint64
	PtrArrayFloat32    *[]float32
	PtrArrayFloat64    *[]float64
	PtrArrayUUID       *[]uuid.UUID
	PtrArrayBoolPtr    *[]*bool
	PtrArrayStringPtr  *[]*string
	PtrArrayIntPtr     *[]*int
	PtrArrayInt8Ptr    *[]*int8
	PtrArrayInt16Ptr   *[]*int16
	PtrArrayInt32Ptr   *[]*int32
	PtrArrayInt64Ptr   *[]*int64
	PtrArrayUintPtr    *[]*uint
	PtrArrayUint8Ptr   *[]*uint8
	PtrArrayUint16Ptr  *[]*uint16
	PtrArrayUint32Ptr  *[]*uint32
	PtrArrayUint64Ptr  *[]*uint64
	PtrArrayFloat32Ptr *[]*float32
	PtrArrayFloat64Ptr *[]*float64
	PtrArrayUUIDPtr    *[]*uuid.UUID
	PtrArrayJSON       *[]any
	PtrArrayJSONPtr    *[]*interface{}
	PtrArrayRawJSON    *[]json.RawMessage
	PtrArrayRawJSONPtr *[]*json.RawMessage
	PtrArrayBigInt     *[]scalar.BigInt
	PtrArrayBigIntPtr  *[]*scalar.BigInt
	PtrArrayTime       *[]time.Time
	PtrArrayTimePtr    *[]*time.Time

	Object struct {
		ID           uuid.UUID                               `json:"id"`
		CreatedAt    time.Time                               `json:"created_at"`
		GenericField types.CustomHeadersResult[types.Author] `json:"generic_field,omitempty"`
	}
	ObjectPtr *struct {
		Long int
		Lat  int
	}
	ArrayObject []struct {
		Content string `json:"content"`
	}
	ArrayObjectPtr *[]struct {
		Content string `json:"content"`
	}
	NamedObject    types.Author
	NamedObjectPtr *types.Author
	NamedArray     []types.Author
	NamedArrayPtr  *[]types.Author

	Map         map[string]any
	MapPtr      *map[string]any
	ArrayMap    []map[string]any
	ArrayMapPtr *[]map[string]any

	JSON       any
	JSONPtr    *interface{}
	RawJSON    json.RawMessage
	RawJSONPtr *json.RawMessage
	Bytes      scalar.Bytes
	BytesPtr   *scalar.Bytes

	UUIDEmpty         uuid.UUID         `json:"uuid_empty,omitempty"`
	BoolEmpty         bool              `json:"bool_empty,omitempty"`
	StringEmpty       string            `json:"string_empty,omitempty"`
	IntEmpty          int               `json:"int_empty,omitempty"`
	Int8Empty         int8              `json:"int8_empty,omitempty"`
	Int16Empty        int16             `json:"int16_empty,omitempty"`
	Int32Empty        int32             `json:"int32_empty,omitempty"`
	Int64Empty        int64             `json:"int64_empty,omitempty"`
	UintEmpty         uint              `json:"uint_empty,omitempty"`
	Uint8Empty        uint8             `json:"uint8_empty,omitempty"`
	Uint16Empty       uint16            `json:"uint16_empty,omitempty"`
	Uint32Empty       uint32            `json:"uint32_empty,omitempty"`
	Uint64Empty       uint64            `json:"uint64_empty,omitempty"`
	Float32Empty      float32           `json:"float32_empty,omitempty"`
	Float64Empty      float64           `json:"float64_empty,omitempty"`
	TimeEmpty         time.Time         `json:"time_empty,omitempty"`
	TextEmpty         types.Text        `json:"text_empty,omitempty"`
	CustomScalarEmpty types.CommentText `json:"custom_scalar_empty,omitempty"`
	EnumEmpty         types.SomeEnum    `json:"enum_empty,omitempty"`
	BigIntEmpty       scalar.BigInt     `json:"bigint_empty,omitempty"`
	DateEmpty         scalar.Date       `json:"date_empty,omitempty"`
	URLEmpty          scalar.URL        `json:"url_empty,omitempty"`
	MapEmpty          map[string]any    `json:"map_empty,omitempty"`
	ArrayMapEmpty     []map[string]any  `json:"array_map_empty,omitempty"`
	JSONEmpty         any               `json:"any_empty,omitempty"`
	RawJSONEmpty      json.RawMessage   `json:"raw_json_empty,omitempty"`

	ArrayBoolEmpty       []bool             `json:"array_bool_empty,omitempty"`
	ArrayStringEmpty     []string           `json:"array_string_empty,omitempty"`
	ArrayIntEmpty        []int              `json:"array_int_empty,omitempty"`
	ArrayInt8Empty       []int8             `json:"array_int8_empty,omitempty"`
	ArrayInt16Empty      []int16            `json:"array_int16_empty,omitempty"`
	ArrayInt32Empty      []int32            `json:"array_int32_empty,omitempty"`
	ArrayInt64Empty      []int64            `json:"array_int64_empty,omitempty"`
	ArrayUintEmpty       []uint             `json:"array_uint_empty,omitempty"`
	ArrayUint8Empty      []uint8            `json:"array_uint8_empty,omitempty"`
	ArrayUint16Empty     []uint16           `json:"array_uint16_empty,omitempty"`
	ArrayUint32Empty     []uint32           `json:"array_uint32_empty,omitempty"`
	ArrayUint64Empty     []uint64           `json:"array_uint64_empty,omitempty"`
	ArrayFloat32Empty    []float32          `json:"array_float32_empty,omitempty"`
	ArrayFloat64Empty    []float64          `json:"array_float64_empty,omitempty"`
	ArrayUUIDEmpty       []uuid.UUID        `json:"array_uuid_empty,omitempty"`
	ArrayBoolPtrEmpty    []*bool            `json:"array_bool_ptr_empty,omitempty"`
	ArrayStringPtrEmpty  []*string          `json:"array_string_ptr_empty,omitempty"`
	ArrayIntPtrEmpty     []*int             `json:"array_int_ptr_empty,omitempty"`
	ArrayInt8PtrEmpty    []*int8            `json:"array_int8_ptr_empty,omitempty"`
	ArrayInt16PtrEmpty   []*int16           `json:"array_int16_ptr_empty,omitempty"`
	ArrayInt32PtrEmpty   []*int32           `json:"array_int32_ptr_empty,omitempty"`
	ArrayInt64PtrEmpty   []*int64           `json:"array_int64_ptr_empty,omitempty"`
	ArrayUintPtrEmpty    []*uint            `json:"array_uint_ptr_empty,omitempty"`
	ArrayUint8PtrEmpty   []*uint8           `json:"array_uint8_ptr_empty,omitempty"`
	ArrayUint16PtrEmpty  []*uint16          `json:"array_uint16_ptr_empty,omitempty"`
	ArrayUint32PtrEmpty  []*uint32          `json:"array_uint32_ptr_empty,omitempty"`
	ArrayUint64PtrEmpty  []*uint64          `json:"array_uint64_ptr_empty,omitempty"`
	ArrayFloat32PtrEmpty []*float32         `json:"array_float32_ptr_empty,omitempty"`
	ArrayFloat64PtrEmpty []*float64         `json:"array_float64_ptr_empty,omitempty"`
	ArrayUUIDPtrEmpty    []*uuid.UUID       `json:"array_uuid_ptr_empty,omitempty"`
	ArrayJSONEmpty       []any              `json:"array_json_empty,omitempty"`
	ArrayJSONPtrEmpty    []*interface{}     `json:"array_json_ptr_empty,omitempty"`
	ArrayRawJSONEmpty    []json.RawMessage  `json:"array_raw_json_empty,omitempty"`
	ArrayRawJSONPtrEmpty []*json.RawMessage `json:"array_raw_json_ptr_empty,omitempty"`
	ArrayBigIntEmpty     []scalar.BigInt    `json:"array_bigint_empty,omitempty"`
	ArrayBigIntPtrEmpty  []*scalar.BigInt   `json:"array_bigint_ptr_empty,omitempty"`
	ArrayTimeEmpty       []time.Time        `json:"array_time_empty,omitempty"`
	ArrayTimePtrEmpty    []*time.Time       `json:"array_time_ptr_empty,omitempty"`

	privateField string
	someFunc     func()
}

type GetCustomHeadersInput struct {
	ID  uuid.UUID `json:"id"`
	Num int       `json:"num"`
}
