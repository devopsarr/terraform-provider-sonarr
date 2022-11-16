package tools

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"golift.io/starr"
)

type Test struct {
	Set types.Set
	Str types.String
	In  types.Int64
	Fl  types.Float64
	Boo types.Bool
}

func TestWriteStringField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldOutput starr.FieldOutput
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: starr.FieldOutput{
				Name:  "str",
				Value: "string",
			},
			written:  Test{},
			expected: Test{Str: types.StringValue("string")},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteStringField(&test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteBoolField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldOutput starr.FieldOutput
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: starr.FieldOutput{
				Name:  "boo",
				Value: true,
			},
			written:  Test{},
			expected: Test{Boo: types.BoolValue(true)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteBoolField(&test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteIntField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldOutput starr.FieldOutput
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: starr.FieldOutput{
				Name: "in",
				// use float to match unmarshal response
				Value: float64(50),
			},
			written:  Test{},
			expected: Test{In: types.Int64Value(50)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteIntField(&test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteFloatField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldOutput starr.FieldOutput
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: starr.FieldOutput{
				Name:  "fl",
				Value: float64(3.5),
			},
			written:  Test{},
			expected: Test{Fl: types.Float64Value(3.5)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			WriteFloatField(&test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteIntSliceField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldOutput starr.FieldOutput
		set         []int64
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: starr.FieldOutput{
				Name: "set",
				// use interface to match unmarshal response
				Value: append(make([]interface{}, 0), 1, 2),
			},
			written:  Test{},
			set:      []int64{1, 2},
			expected: Test{Set: types.SetValueMust(types.Int64Type, nil)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.expected.Set.Type(context.Background()), &test.expected.Set)
			WriteIntSliceField(context.Background(), &test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestWriteStringSliceField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fieldOutput starr.FieldOutput
		set         []string
		written     Test
		expected    Test
	}{
		"working": {
			fieldOutput: starr.FieldOutput{
				Name: "set",
				// use interface to match unmarshal response
				Value: append(make([]interface{}, 0), "test1", "test2"),
			},
			written:  Test{},
			set:      []string{"test1", "test2"},
			expected: Test{Set: types.SetValueMust(types.StringType, nil)},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.expected.Set.Type(context.Background()), &test.expected.Set)
			WriteStringSliceField(context.Background(), &test.fieldOutput, &test.written)
			assert.Equal(t, test.expected, test.written)
		})
	}
}

func TestReadStringField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expected  starr.FieldInput
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Str: types.StringValue("string"),
			},
			name: "str",
			expected: starr.FieldInput{
				Name:  "str",
				Value: "string",
			},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadStringField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, *field)
		})
	}
}

func TestReadIntField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expected  starr.FieldInput
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				In: types.Int64Value(10),
			},
			name: "in",
			expected: starr.FieldInput{
				Name:  "in",
				Value: int64(10),
			},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadIntField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, *field)
		})
	}
}

func TestReadBoolField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expected  starr.FieldInput
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Boo: types.BoolValue(true),
			},
			name: "boo",
			expected: starr.FieldInput{
				Name:  "boo",
				Value: true,
			},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadBoolField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, *field)
		})
	}
}

func TestReadFloatField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expected  starr.FieldInput
		name      string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Fl: types.Float64Value(3.5),
			},
			name: "fl",
			expected: starr.FieldInput{
				Name:  "fl",
				Value: 3.5,
			},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			field := ReadFloatField(test.name, &test.fieldCase)
			assert.Equal(t, test.expected, *field)
		})
	}
}

func TestReadStringSliceField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expected  starr.FieldInput
		name      string
		set       []string
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Set: types.SetValueMust(types.StringType, nil),
			},
			name: "set",
			expected: starr.FieldInput{
				Name:  "set",
				Value: []string{"test1", "test2"},
			},
			set: []string{"test1", "test2"},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.fieldCase.Set.Type(context.Background()), &test.fieldCase.Set)
			field := ReadStringSliceField(context.Background(), test.name, &test.fieldCase)
			assert.Equal(t, test.expected, *field)
		})
	}
}

func TestReadIntSliceField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expected  starr.FieldInput
		name      string
		set       []int64
		fieldCase Test
	}{
		"working": {
			fieldCase: Test{
				Set: types.SetValueMust(types.Int64Type, nil),
			},
			name: "set",
			expected: starr.FieldInput{
				Name:  "set",
				Value: []int64{1, 2},
			},
			set: []int64{1, 2},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tfsdk.ValueFrom(context.Background(), test.set, test.fieldCase.Set.Type(context.Background()), &test.fieldCase.Set)
			field := ReadIntSliceField(context.Background(), test.name, &test.fieldCase)
			assert.Equal(t, test.expected, *field)
		})
	}
}
