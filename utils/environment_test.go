package utils

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestEnvString(t *testing.T) {
	t.Setenv("SOME_FOO", "bar")
	testCases := []struct {
		Input    EnvString
		Expected string
		ErrorMsg string
	}{
		{
			Input:    NewEnvStringValue("foo"),
			Expected: "foo",
		},
		{
			Input:    NewEnvStringVariable("SOME_FOO"),
			Expected: "bar",
		},
		{
			Input:    EnvString{},
			ErrorMsg: errEnvironmentValueRequired.Error(),
		},
		{
			Input: EnvString{
				Value:    ToPtr("foo"),
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvString{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, result, tc.Expected)
			}
		})
	}
}

func TestEnvBool(t *testing.T) {
	t.Setenv("SOME_FOO", "true")
	testCases := []struct {
		Input    EnvBool
		Expected bool
		ErrorMsg string
	}{
		{
			Input:    NewEnvBoolValue(true),
			Expected: true,
		},
		{
			Input:    NewEnvBoolVariable("SOME_FOO"),
			Expected: true,
		},
		{
			Input:    NewEnvBoolVariable("SOME_FOO_2"),
			ErrorMsg: errEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    EnvBool{},
			ErrorMsg: errEnvironmentValueRequired.Error(),
		},
		{
			Input: EnvBool{
				Value:    ToPtr(true),
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvBool{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
				if tc.ErrorMsg == errEnvironmentVariableValueRequired.Error() {
					newValue, err := tc.Input.GetOrDefault(true)
					assert.NilError(t, err)
					assert.Equal(t, newValue, true)
				}
			} else {
				assert.NilError(t, err)
				assert.Equal(t, result, tc.Expected)

				newValue, err := tc.Input.GetOrDefault(true)
				assert.NilError(t, err)
				assert.Equal(t, newValue, tc.Expected)
			}
		})
	}
}

func TestEnvInt(t *testing.T) {
	t.Setenv("SOME_FOO", "10")
	testCases := []struct {
		Input    EnvInt
		Expected int64
		ErrorMsg string
	}{
		{
			Input:    NewEnvIntValue(1),
			Expected: 1,
		},
		{
			Input:    NewEnvIntVariable("SOME_FOO"),
			Expected: 10,
		},
		{
			Input:    NewEnvIntVariable("SOME_FOO_2"),
			ErrorMsg: errEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    EnvInt{},
			ErrorMsg: errEnvironmentValueRequired.Error(),
		},
		{
			Input: EnvInt{
				Value:    ToPtr[int64](10),
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvInt{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
				if tc.ErrorMsg == errEnvironmentVariableValueRequired.Error() {
					newValue, err := tc.Input.GetOrDefault(100)
					assert.NilError(t, err)
					assert.Equal(t, newValue, int64(100))
				}
			} else {
				assert.NilError(t, err)
				assert.Equal(t, result, tc.Expected)

				newValue, err := tc.Input.GetOrDefault(100)
				assert.NilError(t, err)
				assert.Equal(t, newValue, tc.Expected)
			}
		})
	}
}

func TestEnvFloat(t *testing.T) {
	t.Setenv("SOME_FOO", "10.5")
	testCases := []struct {
		Input    EnvFloat
		Expected float64
		ErrorMsg string
	}{
		{
			Input:    NewEnvFloatValue(1.1),
			Expected: 1.1,
		},
		{
			Input:    NewEnvFloatVariable("SOME_FOO"),
			Expected: 10.5,
		},
		{
			Input:    NewEnvFloatVariable("SOME_FOO_2"),
			ErrorMsg: errEnvironmentVariableValueRequired.Error(),
		},
		{
			Input:    EnvFloat{},
			ErrorMsg: errEnvironmentValueRequired.Error(),
		},
		{
			Input: EnvFloat{
				Value:    ToPtr[float64](10),
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvFloat{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
				if tc.ErrorMsg == errEnvironmentVariableValueRequired.Error() {
					newValue, err := tc.Input.GetOrDefault(100.5)
					assert.NilError(t, err)
					assert.Equal(t, newValue, float64(100.5))
				}
			} else {
				assert.NilError(t, err)
				assert.Equal(t, result, tc.Expected)

				newValue, err := tc.Input.GetOrDefault(100)
				assert.NilError(t, err)
				assert.Equal(t, newValue, tc.Expected)
			}
		})
	}
}

func TestEnvMapBool(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=true;bar=false")
	testCases := []struct {
		Input    EnvMapBool
		Expected map[string]bool
		ErrorMsg string
	}{
		{
			Input: NewEnvMapBoolValue(map[string]bool{
				"foo": true,
			}),
			Expected: map[string]bool{
				"foo": true,
			},
		},
		{
			Input: NewEnvMapBoolVariable("SOME_FOO"),
			Expected: map[string]bool{
				"foo": true,
				"bar": false,
			},
		},
		{
			Input:    EnvMapBool{},
			Expected: nil,
		},
		{
			Input: EnvMapBool{
				Value:    map[string]bool{},
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvMapBool{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, result, tc.Expected)
			}
		})
	}
}

func TestEnvMapInt(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=2;bar=3")
	testCases := []struct {
		Input    EnvMapInt
		Expected map[string]int64
		ErrorMsg string
	}{
		{
			Input: NewEnvMapIntValue(map[string]int64{
				"foo": 1,
			}),
			Expected: map[string]int64{
				"foo": 1,
			},
		},
		{
			Input: NewEnvMapIntVariable("SOME_FOO"),
			Expected: map[string]int64{
				"foo": 2,
				"bar": 3,
			},
		},
		{
			Input:    EnvMapInt{},
			Expected: nil,
		},
		{
			Input: EnvMapInt{
				Value:    map[string]int64{},
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvMapInt{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, result, tc.Expected)
			}
		})
	}
}

func TestEnvMapFloat(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=2.2;bar=3.3")
	testCases := []struct {
		Input    EnvMapFloat
		Expected map[string]float64
		ErrorMsg string
	}{
		{
			Input: NewEnvMapFloatValue(map[string]float64{
				"foo": 1.1,
			}),
			Expected: map[string]float64{
				"foo": 1.1,
			},
		},
		{
			Input: NewEnvMapFloatVariable("SOME_FOO"),
			Expected: map[string]float64{
				"foo": 2.2,
				"bar": 3.3,
			},
		},
		{
			Input:    EnvMapFloat{},
			Expected: nil,
		},
		{
			Input: EnvMapFloat{
				Value:    map[string]float64{},
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvMapFloat{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, result, tc.Expected)
				if tc.Input.Variable != nil {
					assert.DeepEqual(t, tc.Input.value, tc.Expected)
				}
			}
		})
	}
}

func TestEnvMapString(t *testing.T) {
	t.Setenv("SOME_FOO", "foo=2.2;bar=3.3")
	testCases := []struct {
		Input    EnvMapString
		Expected map[string]string
		ErrorMsg string
	}{
		{
			Input: NewEnvMapStringValue(map[string]string{
				"foo": "1.1",
			}),
			Expected: map[string]string{
				"foo": "1.1",
			},
		},
		{
			Input: NewEnvMapStringVariable("SOME_FOO"),
			Expected: map[string]string{
				"foo": "2.2",
				"bar": "3.3",
			},
		},
		{
			Input:    EnvMapString{},
			Expected: nil,
		},
		{
			Input: EnvMapString{
				Value:    map[string]string{},
				Variable: ToPtr("SOME_FOO"),
			},
			ErrorMsg: errEnvironmentEitherValueOrEnv.Error(),
		},
		{
			Input: EnvMapString{
				Variable: ToPtr(""),
			},
			ErrorMsg: errEnvironmentVariableRequired.Error(),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result, err := tc.Input.Get()
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, result, tc.Expected)
			}
		})
	}
}
