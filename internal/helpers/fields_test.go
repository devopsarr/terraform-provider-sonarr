package helpers

import (
	"context"
	"testing"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Fl       types.Float64
	Set      types.Set
	Str      types.String
	In       types.Int64
	SeedTime types.Int64
	Boo      types.Bool
}

func TestWriteStringField(t *testing.T) {
	t.Parallel()

	value := "string"

	tests := map[string]struct {
		fieldOutput sonarr.Field
		value       *string
		written     Test
		expected    Test
	}{
		"working": {
			value:    &value,
			written:  Test{},
			expected: Test{Str: types.StringValue(value)},
		},
		"nil": {
			value:    nil,
			written:  Test{},
			expected: Test{Str: types.StringNull()},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := sonarr.NewField()
			field.SetName("str")
			if test.value != nil {
				field.SetValue(*test.value)
			}
			writeStringField(field, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteBoolField(t *testing.T) {
	t.Parallel()

	value := true

	tests := map[string]struct {
		value    *bool
		written  Test
		expected Test
	}{
		"working": {
			value:    &value,
			written:  Test{},
			expected: Test{Boo: types.BoolValue(value)},
		},
		"nil": {
			value:    nil,
			written:  Test{},
			expected: Test{Boo: types.BoolNull()},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := sonarr.NewField()
			field.SetName("boo")
			if test.value != nil {
				field.SetValue(*test.value)
			}
			writeBoolField(field, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteIntField(t *testing.T) {
	t.Parallel()

	value := float64(50)

	tests := map[string]struct {
		// use float to simulate unmarshal response
		value    *float64
		name     string
		written  Test
		expected Test
	}{
		"working": {
			name:     "in",
			value:    &value,
			written:  Test{},
			expected: Test{In: types.Int64Value(50)},
		},
		"seedtime": {
			name:     "seedCriteria.seedTime",
			value:    &value,
			written:  Test{},
			expected: Test{SeedTime: types.Int64Value(50)},
		},
		"nil": {
			name:     "in",
			value:    nil,
			written:  Test{},
			expected: Test{In: types.Int64Null()},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := sonarr.NewField()
			field.SetName(test.name)
			if test.value != nil {
				field.SetValue(*test.value)
			}
			writeIntField(field, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteFloatField(t *testing.T) {
	t.Parallel()

	value := float64(3.5)

	tests := map[string]struct {
		value    *float64
		written  Test
		expected Test
	}{
		"working": {
			value:    &value,
			written:  Test{},
			expected: Test{Fl: types.Float64Value(value)},
		},
		"nil": {},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := sonarr.NewField()
			field.SetName("fl")
			if test.value != nil {
				field.SetValue(*test.value)
			}
			writeFloatField(field, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteIntSliceField(t *testing.T) {
	t.Parallel()

	field := sonarr.NewField()
	field.SetName("set")
	// use interface to simulate unmarshal response
	field.SetValue(append(make([]interface{}, 0), 1, 2))

	tests := map[string]struct {
		fieldOutput sonarr.Field
		set         []int64
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: *field,
			written:     Test{},
			set:         []int64{1, 2},
			expected:    Test{Set: types.SetValueMust(types.Int64Type, nil)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.expected.Set.Type(context.Background()), &test.expected.Set)
			writeIntSliceField(context.Background(), &test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteStringSliceField(t *testing.T) {
	t.Parallel()

	field := sonarr.NewField()
	field.SetName("set")
	// use interface to simulate unmarshal response
	field.SetValue(append(make([]interface{}, 0), "test1", "test2"))

	tests := map[string]struct {
		fieldOutput sonarr.Field
		set         []string
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: *field,
			written:     Test{},
			set:         []string{"test1", "test2"},
			expected:    Test{Set: types.SetValueMust(types.StringType, nil)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.expected.Set.Type(context.Background()), &test.expected.Set)
			writeStringSliceField(context.Background(), &test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestReadStringField(t *testing.T) {
	t.Parallel()

	field := sonarr.NewField()
	field.SetName("str")
	field.SetValue("string")

	tests := map[string]struct {
		expected  *sonarr.Field
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Str: types.StringValue("string"),
			},
			name:     "str",
			expected: field,
		},
		"nil": {
			fieldCase: Test{},
			name:      "str",
			expected:  nil,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := readStringField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadIntField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name      string
		tfName    string
		fieldCase Test
		value     int
	}{
		"working": {
			fieldCase: Test{
				In: types.Int64Value(10),
			},
			name:   "in",
			tfName: "in",
			value:  10,
		},
		"nil": {
			fieldCase: Test{},
			name:      "in",
			tfName:    "in",
			value:     0,
		},
		"seedtime": {
			fieldCase: Test{
				SeedTime: types.Int64Value(10),
			},
			name:   "seedCriteria.seedTime",
			tfName: "seedTime",
			value:  10,
		},
	}
	for name, test := range tests {
		test := test

		expected := sonarr.NewField()
		expected.SetName(test.name)
		expected.SetValue(int64(test.value))

		if test.value == 0 {
			expected = nil
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := readIntField(test.tfName, &test.fieldCase)
			assert.Equal(t, expected, field)
		})
	}
}

func TestReadBoolField(t *testing.T) {
	t.Parallel()

	field := sonarr.NewField()
	field.SetName("boo")
	field.SetValue(true)

	tests := map[string]struct {
		expected  *sonarr.Field
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Boo: types.BoolValue(true),
			},
			name:     "boo",
			expected: field,
		},
		"nil": {
			fieldCase: Test{},
			name:      "boo",
			expected:  nil,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := readBoolField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadFloatField(t *testing.T) {
	t.Parallel()

	field := sonarr.NewField()
	field.SetName("fl")
	field.SetValue(3.5)

	tests := map[string]struct {
		expected  *sonarr.Field
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Fl: types.Float64Value(3.5),
			},
			name:     "fl",
			expected: field,
		},
		"nil": {
			fieldCase: Test{},
			name:      "fl",
			expected:  nil,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := readFloatField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadStringSliceField(t *testing.T) {
	t.Parallel()

	field := sonarr.NewField()
	field.SetName("set")
	field.SetValue([]string{"test1", "test2"})

	tests := map[string]struct {
		expected  *sonarr.Field
		name      string
		set       []string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Set: types.SetValueMust(types.StringType, nil),
			},
			name:     "set",
			expected: field,
			set:      []string{"test1", "test2"},
		},
		"nil": {
			fieldCase: Test{
				Set: types.SetValueMust(types.StringType, nil),
			},
			name:     "set",
			expected: nil,
			set:      []string{},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.fieldCase.Set.Type(context.Background()), &test.fieldCase.Set)
			field := readStringSliceField(context.Background(), test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadIntSliceField(t *testing.T) {
	t.Parallel()

	field := sonarr.NewField()
	field.SetName("set")
	field.SetValue([]int64{1, 2})

	tests := map[string]struct {
		expected  *sonarr.Field
		name      string
		set       []int64
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Set: types.SetValueMust(types.Int64Type, nil),
			},
			name:     "set",
			expected: field,
			set:      []int64{1, 2},
		},
		"nil": {
			fieldCase: Test{
				Set: types.SetValueMust(types.Int64Type, nil),
			},
			name:     "set",
			expected: nil,
			set:      []int64{},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.fieldCase.Set.Type(context.Background()), &test.fieldCase.Set)
			field := readIntSliceField(context.Background(), test.name, &test.fieldCase)
			assert.Equal(t, test.expected, field)
		})
	}
}

func TestReadFields(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldLists Fields
		name       string
		value      interface{}
		testData   Test
	}{
		"string": {
			fieldLists: Fields{Strings: []string{"str"}},
			name:       "str",
			value:      "String",
			testData:   Test{Str: types.StringValue("String")},
		},
		"int": {
			fieldLists: Fields{Ints: []string{"in"}},
			name:       "in",
			value:      int64(55),
			testData:   Test{In: types.Int64Value(55)},
		},
		"bool": {
			fieldLists: Fields{Bools: []string{"boo"}},
			name:       "boo",
			value:      true,
			testData:   Test{Boo: types.BoolValue(true)},
		},
		"float": {
			fieldLists: Fields{Floats: []string{"fl"}},
			name:       "fl",
			value:      5.5,
			testData:   Test{Fl: types.Float64Value(5.5)},
		},
		"stringSlice": {
			fieldLists: Fields{StringSlices: []string{"set"}},
			name:       "set",
			value:      []string{"test1", "test2"},
			testData:   Test{Set: types.SetValueMust(types.StringType, nil)},
		},
		"intSlice": {
			fieldLists: Fields{IntSlices: []string{"set"}},
			name:       "set",
			value:      []int64{1, 9},
			testData:   Test{Set: types.SetValueMust(types.Int64Type, nil)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if test.name == "set" {
				tfsdk.ValueFrom(context.Background(), test.value, test.testData.Set.Type(context.Background()), &test.testData.Set)
			}

			expectedFields := make([]*sonarr.Field, 1)
			expectedFields[0] = sonarr.NewField()
			expectedFields[0].SetName(test.name)
			expectedFields[0].SetValue(test.value)

			fields := ReadFields(context.Background(), &test.testData, test.fieldLists)
			assert.Equal(t, &expectedFields, &fields)
		})
	}
}

func TestWriteFields(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldLists     Fields
		name           string
		value          interface{}
		fieldContainer Test
	}{
		"string": {
			fieldLists:     Fields{Strings: []string{"str"}},
			name:           "str",
			value:          "String",
			fieldContainer: Test{Str: types.StringValue("String")},
		},
		"int": {
			fieldLists:     Fields{Ints: []string{"in"}},
			name:           "in",
			value:          float64(55),
			fieldContainer: Test{In: types.Int64Value(55)},
		},
		"bool": {
			fieldLists:     Fields{Bools: []string{"boo"}},
			name:           "boo",
			value:          true,
			fieldContainer: Test{Boo: types.BoolValue(true)},
		},
		"float": {
			fieldLists:     Fields{Floats: []string{"fl"}},
			name:           "fl",
			value:          5.5,
			fieldContainer: Test{Fl: types.Float64Value(5.5)},
		},
		"stringSlice": {
			fieldLists:     Fields{StringSlices: []string{"set"}},
			name:           "set",
			value:          append(make([]interface{}, 0), []string{"test1", "test2"}),
			fieldContainer: Test{Set: types.SetValueMust(types.StringType, nil)},
		},
		"intSlice": {
			fieldLists:     Fields{IntSlices: []string{"set"}},
			name:           "set",
			value:          append(make([]interface{}, 0), []float64{1, 9}),
			fieldContainer: Test{Set: types.SetValueMust(types.Int64Type, nil)},
		},
		"intSliceException": {
			fieldLists:     Fields{IntSlicesExceptions: []string{"set"}},
			name:           "set",
			value:          append(make([]interface{}, 0), []float64{1, 9}),
			fieldContainer: Test{Set: types.SetValueMust(types.Int64Type, nil)},
		},
		"stringSliceException": {
			fieldLists:     Fields{StringSlicesExceptions: []string{"set"}},
			name:           "set",
			value:          append(make([]interface{}, 0), []string{"test1", "test2"}),
			fieldContainer: Test{Set: types.SetValueMust(types.StringType, nil)},
		},
		"sensitive": {
			fieldLists:     Fields{Sensitive: []string{"str"}},
			name:           "str",
			value:          "String",
			fieldContainer: Test{Str: types.StringValue("String")},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if test.name == "set" {
				tfsdk.ValueFrom(context.Background(), test.value, test.fieldContainer.Set.Type(context.Background()), &test.fieldContainer.Set)
			}

			fields := make([]*sonarr.Field, 1)
			fields[0] = sonarr.NewField()
			fields[0].SetName(test.name)
			fields[0].SetValue(test.value)

			container := Test{}
			if len(test.fieldLists.Sensitive) != 0 {
				// emulate the sensitive behaviour
				container = Test{
					Str: types.StringValue("String"),
				}
			}

			WriteFields(context.TODO(), &container, fields, test.fieldLists)
			assert.Equal(t, &test.fieldContainer, &container)
		})
	}
}
