package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-test/types"
	"github.com/hasura/ndc-sdk-go/scalar"
)

type Text string

// A foo scalar
type ScalarFoo struct {
	bar string
}

// FromValue decodes the scalar from an unknown value
func (s *ScalarFoo) FromValue(value any) error {
	s.bar = fmt.Sprint(value)
	return nil
}

// A hello result
type HelloResult struct {
	ID   uuid.UUID `json:"id"`
	Num  int       `json:"num"`
	Text Text      `json:"text"`
	Foo  ScalarFoo `json:"foo"`
}

// FunctionHello sends a hello message
func FunctionHello(ctx context.Context, state *types.State) (*HelloResult, error) {
	return &HelloResult{
		ID:   uuid.New(),
		Num:  1,
		Text: "world",
	}, nil
}

// A create author argument
type CreateAuthorArguments struct {
	Name string `json:"name"`
}

// A create authors argument
type CreateAuthorsArguments struct {
	Names []string `json:"names"`
}

// A create author result
type CreateAuthorResult struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ProcedureCreateAuthor creates an author
func ProcedureCreateAuthor(ctx context.Context, state *types.State, arguments *CreateAuthorArguments) (*CreateAuthorResult, error) {
	return &CreateAuthorResult{
		ID:   1,
		Name: arguments.Name,
	}, nil
}

// ProcedureCreateAuthors creates a list of authors
func ProcedureCreateAuthors(ctx context.Context, state *types.State, arguments *CreateAuthorsArguments) ([]CreateAuthorResult, error) {
	return []CreateAuthorResult{
		{
			ID:   1,
			Name: strings.Join(arguments.Names, ","),
		},
	}, nil
}

// FunctionGetBool return an scalar boolean
func FunctionGetBool(ctx context.Context, state *types.State) (bool, error) {
	return true, nil
}

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
	Text         Text
	CustomScalar CommentText
	Enum         SomeEnum
	BigInt       scalar.BigInt
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
	TextPtr         *Text
	CustomScalarPtr *CommentText
	EnumPtr         *SomeEnum
	BigIntPtr       *scalar.BigInt

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
	NamedObject    Author
	NamedObjectPtr *Author
	NamedArray     []Author
	NamedArrayPtr  *[]Author
	UUIDArray      []uuid.UUID

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

func FunctionGetTypes(ctx context.Context, state *types.State, arguments *GetTypesArguments) (*GetTypesArguments, error) {
	return arguments, nil
}
