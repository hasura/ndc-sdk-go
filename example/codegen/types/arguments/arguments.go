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
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
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
}

type GetCustomHeadersInput struct {
	ID  uuid.UUID `json:"id"`
	Num int       `json:"num"`
}
