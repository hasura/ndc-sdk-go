package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	errEnvironmentValueRequired         = errors.New("require either value or env")
	errEnvironmentVariableRequired      = errors.New("the variable name of env is empty")
	errEnvironmentVariableValueRequired = errors.New("the variable value of env is empty")
)

// EnvString represents either a literal string or an environment reference
type EnvString struct {
	Value    *string `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvString creates an EnvString instance
func NewEnvString(env string, value string) EnvString {
	return EnvString{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvStringValue creates an EnvString with a literal value
func NewEnvStringValue(value string) EnvString {
	return EnvString{
		Value: &value,
	}
}

// NewEnvStringVariable creates an EnvString with a variable name
func NewEnvStringVariable(name string) EnvString {
	return EnvString{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvString) UnmarshalJSON(b []byte) error {
	type Plain EnvString
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentValue(rawValue.Value, rawValue.Variable); err != nil {
		return fmt.Errorf("EnvString: %w", err)
	}
	*ev = EnvString(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvString) Get() (string, error) {
	if err := validateEnvironmentValue(ev.Value, ev.Variable); err != nil {
		return "", err
	}
	if ev.Variable != nil {
		value := os.Getenv(*ev.Variable)
		if value != "" {
			return value, nil
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return "", errEnvironmentVariableValueRequired
}

// EnvInt represents either a literal integer or an environment reference
type EnvInt struct {
	Value    *int64  `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvInt creates an EnvInt instance.
func NewEnvInt(env string, value int64) EnvInt {
	return EnvInt{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvIntValue creates an EnvInt with a literal value
func NewEnvIntValue(value int64) EnvInt {
	return EnvInt{
		Value: &value,
	}
}

// NewEnvIntVariable creates an EnvInt with a variable name
func NewEnvIntVariable(name string) EnvInt {
	return EnvInt{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvInt) UnmarshalJSON(b []byte) error {
	type Plain EnvInt
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentValue(rawValue.Value, rawValue.Variable); err != nil {
		return fmt.Errorf("EnvInt: %w", err)
	}
	*ev = EnvInt(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvInt) Get() (int64, error) {
	if err := validateEnvironmentValue(ev.Value, ev.Variable); err != nil {
		return 0, err
	}
	if ev.Variable != nil {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseInt(rawValue, 10, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, errEnvironmentVariableValueRequired
}

// GetOrDefault gets literal value or from system environment.
// Returns the default value if the environment value is empty
func (ev EnvInt) GetOrDefault(defaultValue int64) (int64, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, errEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}
		return 0, err
	}
	return result, nil
}

// EnvBool represents either a literal boolean or an environment reference
type EnvBool struct {
	Value    *bool   `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvBool creates an EnvBool instance.
func NewEnvBool(env string, value bool) EnvBool {
	return EnvBool{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvBoolValue creates an EnvBool with a literal value
func NewEnvBoolValue(value bool) EnvBool {
	return EnvBool{
		Value: &value,
	}
}

// NewEnvBoolVariable creates an EnvBool with a variable name
func NewEnvBoolVariable(name string) EnvBool {
	return EnvBool{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvBool) UnmarshalJSON(b []byte) error {
	type Plain EnvBool
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentValue(rawValue.Value, rawValue.Variable); err != nil {
		return fmt.Errorf("EnvBool: %w", err)
	}
	*ev = EnvBool(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvBool) Get() (bool, error) {
	if err := validateEnvironmentValue(ev.Value, ev.Variable); err != nil {
		return false, err
	}
	if ev.Variable != nil {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseBool(rawValue)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return false, errEnvironmentVariableValueRequired
}

// GetOrDefault gets literal value or from system environment.
// Returns the default value if the environment value is empty
func (ev EnvBool) GetOrDefault(defaultValue bool) (bool, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, errEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}
		return false, err
	}
	return result, nil
}

// EnvFloat represents either a literal floating point number or an environment reference
type EnvFloat struct {
	Value    *float64 `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string  `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvFloat creates an EnvFloat instance.
func NewEnvFloat(env string, value float64) EnvFloat {
	return EnvFloat{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvFloatValue creates an EnvFloat with a literal value
func NewEnvFloatValue(value float64) EnvFloat {
	return EnvFloat{
		Value: &value,
	}
}

// NewEnvFloatVariable creates an EnvFloat with a variable name
func NewEnvFloatVariable(name string) EnvFloat {
	return EnvFloat{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvFloat) UnmarshalJSON(b []byte) error {
	type Plain EnvFloat
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentValue(rawValue.Value, rawValue.Variable); err != nil {
		return fmt.Errorf("EnvFloat: %w", err)
	}
	*ev = EnvFloat(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvFloat) Get() (float64, error) {
	if err := validateEnvironmentValue(ev.Value, ev.Variable); err != nil {
		return 0, err
	}
	if ev.Variable != nil {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return strconv.ParseFloat(rawValue, 64)
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	return 0, errEnvironmentVariableValueRequired
}

// GetOrDefault gets literal value or from system environment.
// Returns the default value if the environment value is empty
func (ev EnvFloat) GetOrDefault(defaultValue float64) (float64, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, errEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}
		return 0, err
	}
	return result, nil
}

func validateEnvironmentValue[T any](value *T, variable *string) error {
	if value == nil && variable == nil {
		return errEnvironmentValueRequired
	}
	if variable != nil && *variable == "" {
		return errEnvironmentVariableRequired
	}

	return nil
}

func validateEnvironmentMapValue(variable *string) error {
	if variable != nil && *variable == "" {
		return errEnvironmentVariableRequired
	}

	return nil
}

// EnvMapString represents either a literal string map or an environment reference
type EnvMapString struct {
	Value    map[string]string `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string           `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvMapString creates an EnvMapString instance.
func NewEnvMapString(env string, value map[string]string) EnvMapString {
	return EnvMapString{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapStringValue creates an EnvMapString with a literal value
func NewEnvMapStringValue(value map[string]string) EnvMapString {
	return EnvMapString{
		Value: value,
	}
}

// NewEnvMapStringVariable creates an EnvMapString with a variable name
func NewEnvMapStringVariable(name string) EnvMapString {
	return EnvMapString{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapString) UnmarshalJSON(b []byte) error {
	type Plain EnvMapString
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentMapValue(rawValue.Variable); err != nil {
		return fmt.Errorf("EnvMapString: %w", err)
	}
	*ev = EnvMapString(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvMapString) Get() (map[string]string, error) {
	if err := validateEnvironmentMapValue(ev.Variable); err != nil {
		return nil, err
	}
	if ev.Variable != nil {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseStringMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapInt represents either a literal int map or an environment reference
type EnvMapInt struct {
	Value    map[string]int64 `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string          `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvMapInt creates an EnvMapInt instance.
func NewEnvMapInt(env string, value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapIntValue creates an EnvMapInt with a literal value
func NewEnvMapIntValue(value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Value: value,
	}
}

// NewEnvMapIntVariable creates an EnvMapInt with a variable name
func NewEnvMapIntVariable(name string) EnvMapInt {
	return EnvMapInt{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapInt) UnmarshalJSON(b []byte) error {
	type Plain EnvMapInt
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentMapValue(rawValue.Variable); err != nil {
		return fmt.Errorf("EnvMapInt: %w", err)
	}
	*ev = EnvMapInt(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvMapInt) Get() (map[string]int64, error) {
	if err := validateEnvironmentMapValue(ev.Variable); err != nil {
		return nil, err
	}
	if ev.Variable != nil {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseIntegerMapFromString[int64](rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapFloat represents either a literal float map or an environment reference
type EnvMapFloat struct {
	Value    map[string]float64 `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string            `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvMapFloat creates an EnvMapFloat instance.
func NewEnvMapFloat(env string, value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapFloatValue creates an EnvMapFloat with a literal value
func NewEnvMapFloatValue(value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Value: value,
	}
}

// NewEnvMapFloatVariable creates an EnvMapFloat with a variable name
func NewEnvMapFloatVariable(name string) EnvMapFloat {
	return EnvMapFloat{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapFloat) UnmarshalJSON(b []byte) error {
	type Plain EnvMapFloat
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentMapValue(rawValue.Variable); err != nil {
		return fmt.Errorf("EnvMapFloat: %w", err)
	}
	*ev = EnvMapFloat(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvMapFloat) Get() (map[string]float64, error) {
	if err := validateEnvironmentMapValue(ev.Variable); err != nil {
		return nil, err
	}
	if ev.Variable != nil {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseFloatMapFromString[float64](rawValue)
		}
	}

	return ev.Value, nil
}

// EnvMapBool represents either a literal bool map or an environment reference
type EnvMapBool struct {
	Value    map[string]bool `json:"value,omitempty" yaml:"value,omitempty" jsonschema:"anyof_required=value"`
	Variable *string         `json:"env,omitempty" yaml:"env,omitempty" jsonschema:"anyof_required=env"`
}

// NewEnvMapBool creates an EnvMapBool instance.
func NewEnvMapBool(env string, value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapBoolValue creates an EnvMapBool with a literal value
func NewEnvMapBoolValue(value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Value: value,
	}
}

// NewEnvMapBoolVariable creates an EnvMapBool with a variable name
func NewEnvMapBoolVariable(name string) EnvMapBool {
	return EnvMapBool{
		Variable: &name,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (ev *EnvMapBool) UnmarshalJSON(b []byte) error {
	type Plain EnvMapBool
	var rawValue Plain
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}
	if err := validateEnvironmentMapValue(rawValue.Variable); err != nil {
		return fmt.Errorf("EnvMapBool: %w", err)
	}
	*ev = EnvMapBool(rawValue)
	return nil
}

// Get gets literal value or from system environment
func (ev EnvMapBool) Get() (map[string]bool, error) {
	if err := validateEnvironmentMapValue(ev.Variable); err != nil {
		return nil, err
	}
	if ev.Variable != nil {
		rawValue := os.Getenv(*ev.Variable)
		if rawValue != "" {
			return ParseBoolMapFromString(rawValue)
		}
	}

	return ev.Value, nil
}
