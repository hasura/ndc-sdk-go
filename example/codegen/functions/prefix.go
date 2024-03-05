package functions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
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
	Complex64    complex64
	Complex128   complex128
	Time         time.Time
	Duration     time.Duration
	CustomScalar CommentText

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
	Complex64Ptr    *complex64
	Complex128Ptr   *complex128
	CustomScalarPtr *CommentText

	Object struct {
		ID        uuid.UUID  `json:"id"`
		Decimal   complex128 `json:"decimal"`
		CreatedAt time.Time  `json:"created_at"`
	} `json:"author"`
	ObjectPtr *struct {
		Long int
		Lat  int
	}
	ArrayObject []struct {
		Content string `json:"content"`
	}
	NamedObject    Author
	NamedObjectPtr *Author
	NamedArray     []Author
}

func FunctionGetTypes(ctx context.Context, state *types.State, arguments *GetTypesArguments) (*GetTypesArguments, error) {
	return &GetTypesArguments{}, nil
}
