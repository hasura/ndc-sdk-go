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
	errEnvironmentVariableRequired      = errors.New("the environment variable name is empty")
	errEnvironmentVariableValueRequired = errors.New("the environment variable value is empty")
)

// EnvString represents either a literal string or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvString struct {
	Value    *string `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvString creates an EnvString instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvString(env string, value string) EnvString {
	return EnvString{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvStringValue creates an EnvString with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvStringValue(value string) EnvString {
	return EnvString{
		Value: &value,
	}
}

// NewEnvStringVariable creates an EnvString with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func (ev EnvString) Get() (string, error) {
	if err := validateEnvironmentValue(ev.Value, ev.Variable); err != nil {
		return "", err
	}

	var value string

	var envExisted bool
	if ev.Variable != nil {
		value, envExisted = os.LookupEnv(*ev.Variable)
		if value != "" {
			return value, nil
		}
	}

	if ev.Value != nil {
		return *ev.Value, nil
	}

	if envExisted {
		return "", nil
	}

	return "", getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func (ev EnvString) GetOrDefault(defaultValue string) (string, error) {
	result, err := ev.Get()
	if err != nil {
		if errors.Is(err, errEnvironmentVariableValueRequired) {
			return defaultValue, nil
		}

		return "", err
	} else if result == "" {
		result = defaultValue
	}

	return result, nil
}

// EnvInt represents either a literal integer or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvInt struct {
	Value    *int64  `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvInt creates an EnvInt instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvInt(env string, value int64) EnvInt {
	return EnvInt{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvIntValue creates an EnvInt with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvIntValue(value int64) EnvInt {
	return EnvInt{
		Value: &value,
	}
}

// NewEnvIntVariable creates an EnvInt with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// EnvBool represents either a literal boolean or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvBool struct {
	Value    *bool   `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvBool creates an EnvBool instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvBool(env string, value bool) EnvBool {
	return EnvBool{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvBoolValue creates an EnvBool with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvBoolValue(value bool) EnvBool {
	return EnvBool{
		Value: &value,
	}
}

// NewEnvBoolVariable creates an EnvBool with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

	return false, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// EnvFloat represents either a literal floating point number or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvFloat struct {
	Value    *float64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string  `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvFloat creates an EnvFloat instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvFloat(env string, value float64) EnvFloat {
	return EnvFloat{
		Variable: &env,
		Value:    &value,
	}
}

// NewEnvFloatValue creates an EnvFloat with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvFloatValue(value float64) EnvFloat {
	return EnvFloat{
		Value: &value,
	}
}

// NewEnvFloatVariable creates an EnvFloat with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

	return 0, getEnvVariableValueRequiredError(ev.Variable)
}

// GetOrDefault returns the default value if the environment value is empty.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// EnvMapString represents either a literal string map or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvMapString struct {
	Value    map[string]string `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string           `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapString creates an EnvMapString instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapString(env string, value map[string]string) EnvMapString {
	return EnvMapString{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapStringValue creates an EnvMapString with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapStringValue(value map[string]string) EnvMapString {
	return EnvMapString{
		Value: value,
	}
}

// NewEnvMapStringVariable creates an EnvMapString with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// EnvMapInt represents either a literal int map or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvMapInt struct {
	Value    map[string]int64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string          `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapInt creates an EnvMapInt instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapInt(env string, value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapIntValue creates an EnvMapInt with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapIntValue(value map[string]int64) EnvMapInt {
	return EnvMapInt{
		Value: value,
	}
}

// NewEnvMapIntVariable creates an EnvMapInt with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// EnvMapFloat represents either a literal float map or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvMapFloat struct {
	Value    map[string]float64 `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string            `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapFloat creates an EnvMapFloat instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapFloat(env string, value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapFloatValue creates an EnvMapFloat with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapFloatValue(value map[string]float64) EnvMapFloat {
	return EnvMapFloat{
		Value: value,
	}
}

// NewEnvMapFloatVariable creates an EnvMapFloat with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// EnvMapBool represents either a literal bool map or an environment reference.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
type EnvMapBool struct {
	Value    map[string]bool `json:"value,omitempty" jsonschema:"anyof_required=value" mapstructure:"value" yaml:"value,omitempty"`
	Variable *string         `json:"env,omitempty"   jsonschema:"anyof_required=env"   mapstructure:"env"   yaml:"env,omitempty"`
}

// NewEnvMapBool creates an EnvMapBool instance.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapBool(env string, value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Variable: &env,
		Value:    value,
	}
}

// NewEnvMapBoolValue creates an EnvMapBool with a literal value.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
func NewEnvMapBoolValue(value map[string]bool) EnvMapBool {
	return EnvMapBool{
		Value: value,
	}
}

// NewEnvMapBoolVariable creates an EnvMapBool with a variable name.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

// Get gets literal value or from system environment.
//
// Deprecated: this module was moved to github.com/hasura/goenvconf.
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

func getEnvVariableValueRequiredError(envName *string) error {
	if envName != nil {
		return fmt.Errorf("%s: %w", *envName, errEnvironmentVariableValueRequired)
	}

	return errEnvironmentVariableValueRequired
}
