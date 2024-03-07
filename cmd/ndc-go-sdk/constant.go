package main

const (
	connectorOutputFile   = "connector.generated.go"
	schemaOutputFile      = "schema.generated.json"
	typeMethodsOutputFile = "types.generated.go"
	googleUuidPackageName = "github.com/google/uuid"
)

const textBlockErrorCheck = `
    if err != nil {
		  return err
    }
`

const textBlockErrorCheck2 = `
    if err != nil {
      return nil, err
    }
`
const textBlockUUIDParsers = `
func decodeUUIDHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
		if to.PkgPath() != "github.com/google/uuid" || to.Name() != "UUID" {
			return data, nil
		}
		result, err := _parseNullableUUID(data)
		if err != nil || result == nil {
			return uuid.UUID{}, err
		}

		return *result, nil
	}
}

func _parseUUID(value any) (uuid.UUID, error) {
	result, err := _parseNullableUUID(value)
	if err != nil {
		return uuid.UUID{}, err
	}
	if result == nil {
		return uuid.UUID{}, errors.New("the uuid value must not be null")
	}
	return *result, nil
}

func _parseNullableUUID(value any) (*uuid.UUID, error) {
	if utils.IsNil(value) {
		return nil, nil
	}
	switch v := value.(type) {
	case string:
		result, err := uuid.Parse(v)
		if err != nil {
			return nil, err
		}
		return &result, nil
	case *string:
		if v == nil {
			return nil, nil
		}
		result, err := uuid.Parse(*v)
		if err != nil {
			return nil, err
		}
		return &result, nil
	case [16]byte:
		result := uuid.UUID(v)
		return &result, nil
	case *[16]byte:
		if v == nil {
			return nil, nil
		}
		result := uuid.UUID(*v)
		return &result, nil
	default:
		return nil, fmt.Errorf("failed to parse uuid, got: %+v", value)
	}
}

func _getObjectUUID(object map[string]any, key string) (uuid.UUID, error) {
	value, ok := utils.GetAny(object, key)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("field %s is required", key)
	}
	result, err := _parseUUID(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

func _getNullableObjectUUID(object map[string]any, key string) (*uuid.UUID, error) {
	value, ok := utils.GetAny(object, key)
	if !ok {
		return nil, nil
	}
	result, err := _parseNullableUUID(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}
`
