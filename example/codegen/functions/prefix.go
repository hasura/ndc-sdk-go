package functions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/utils"
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
	TimePtr         *time.Time
	DurationPtr     *time.Duration
	CustomScalarPtr *CommentText

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
}

func FunctionGetTypes(ctx context.Context, state *types.State, arguments *GetTypesArguments) (*GetTypesArguments, error) {
	return &GetTypesArguments{
		UUID:     uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5a"),
		Bool:     true,
		String:   "hello",
		Int:      1,
		Int8:     2,
		Int16:    3,
		Int32:    4,
		Int64:    5,
		Uint:     6,
		Uint8:    7,
		Uint16:   8,
		Uint32:   9,
		Uint64:   10,
		Float32:  1.1,
		Float64:  2.2,
		Time:     time.Date(2023, 3, 5, 7, 0, 56, 0, time.UTC),
		Duration: 10 * time.Second,
		CustomScalar: CommentText{
			comment: "a comment",
		},
		UUIDPtr:     utils.ToPtr(uuid.MustParse("b085b0b9-007c-440e-9661-0d8f2de98a5b")),
		BoolPtr:     utils.ToPtr(true),
		IntPtr:      utils.ToPtr(11),
		Int8Ptr:     utils.ToPtr(int8(12)),
		Int16Ptr:    utils.ToPtr(int16(13)),
		Int32Ptr:    utils.ToPtr(int32(14)),
		Int64Ptr:    utils.ToPtr(int64(15)),
		UintPtr:     utils.ToPtr(uint(16)),
		Uint8Ptr:    utils.ToPtr(uint8(17)),
		Uint16Ptr:   utils.ToPtr(uint16(18)),
		Uint32Ptr:   utils.ToPtr(uint32(19)),
		Uint64Ptr:   utils.ToPtr(uint64(20)),
		Float32Ptr:  utils.ToPtr(float32(3.3)),
		Float64Ptr:  utils.ToPtr(float64(4.4)),
		TimePtr:     utils.ToPtr(time.Date(2023, 3, 5, 7, 0, 0, 0, time.UTC)),
		DurationPtr: utils.ToPtr(time.Minute),
		CustomScalarPtr: &CommentText{
			comment: "a comment pointer",
		},
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
		ArrayObjectPtr: &[]struct {
			Content string "json:\"content\""
		}{
			{
				Content: "a content pointer",
			},
		},
		NamedObject: Author{
			ID:        "1",
			Duration:  10 * time.Minute,
			CreatedAt: time.Date(2023, 3, 5, 5, 0, 0, 0, time.UTC),
		},
		NamedObjectPtr: &Author{
			ID:        "2",
			Duration:  11 * time.Minute,
			CreatedAt: time.Date(2023, 3, 5, 4, 0, 0, 0, time.UTC),
		},
		NamedArray: []Author{
			{
				ID:        "3",
				Duration:  12 * time.Minute,
				CreatedAt: time.Date(2023, 3, 5, 3, 0, 0, 0, time.UTC),
			},
		},
		NamedArrayPtr: nil,
	}, nil
}
