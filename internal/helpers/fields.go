package helpers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type fieldException struct {
	apiName string
	tfName  string
}

// getFieldExceptions identifies the fields resources in which TF and API name differs.
func getFieldExceptions() []fieldException {
	return []fieldException{
		{
			apiName: "tags",
			tfName:  "fieldTags",
		},
		{
			apiName: "seedCriteria.seedTime",
			tfName:  "seedTime",
		},
		{
			apiName: "seedCriteria.seedRatio",
			tfName:  "seedRatio",
		},
		{
			apiName: "seedCriteria.seasonPackSeedTime",
			tfName:  "seasonPackSeedTime",
		},
	}
}

// selectTFName identifies the TF name starting from API name.
func selectTFName(name string) string {
	for _, f := range getFieldExceptions() {
		if f.apiName == name {
			name = f.tfName
		}
	}

	return name
}

// selectAPIName identifies the API name starting from TF name.
func selectAPIName(name string) string {
	for _, f := range getFieldExceptions() {
		if f.tfName == name {
			name = f.apiName
		}
	}

	return name
}

// selectWriteField identifies which struct field should be written.
func selectWriteField(fieldOutput *sonarr.Field, fieldCase interface{}) reflect.Value {
	fieldName := selectTFName(fieldOutput.GetName())
	value := reflect.ValueOf(fieldCase).Elem()

	return value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
}

// selectReadField identifies which struct field should be read.
func selectReadField(name string, fieldCase interface{}) reflect.Value {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	return value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
}

// setField sets the sonarr field value.
func setField(name string, value interface{}) *sonarr.Field {
	field := sonarr.NewField()
	field.SetName(name)
	field.SetValue(value)

	return field
}

// WriteStringField writes a sonarr string field into struct field.
func WriteStringField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	stringValue := fmt.Sprint(fieldOutput.GetValue())
	v := reflect.ValueOf(types.StringValue(stringValue))
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// WriteBoolField writes a sonarr bool field into struct field.
func WriteBoolField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	boolValue, _ := fieldOutput.GetValue().(bool)
	v := reflect.ValueOf(types.BoolValue(boolValue))
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// WriteIntField writes a sonarr int field into struct field.
func WriteIntField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	intValue, _ := fieldOutput.GetValue().(float64)
	v := reflect.ValueOf(types.Int64Value(int64(intValue)))
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// WriteFloatField writes a sonarr float field into struct field.
func WriteFloatField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	floatValue, _ := fieldOutput.GetValue().(float64)
	v := reflect.ValueOf(types.Float64Value(floatValue))
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// WriteStringSliceField writes a sonarr string slice field into struct field.
func WriteStringSliceField(ctx context.Context, fieldOutput *sonarr.Field, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.GetValue().([]interface{})
	setValue := types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)
	v := reflect.ValueOf(setValue)
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// WriteIntSliceField writes a sonarr int slice field into struct field.
func WriteIntSliceField(ctx context.Context, fieldOutput *sonarr.Field, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.GetValue().([]interface{})
	setValue := types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)
	v := reflect.ValueOf(setValue)
	selectWriteField(fieldOutput, fieldCase).Set(v)
}

// ReadStringField reads from a string struct field and return a sonarr field.
func ReadStringField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	stringField := (*types.String)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !stringField.IsNull() && !stringField.IsUnknown() {
		return setField(fieldName, stringField.ValueString())
	}

	return nil
}

// ReadBoolField reads from a bool struct field and return a sonarr field.
func ReadBoolField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	boolField := (*types.Bool)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !boolField.IsNull() && !boolField.IsUnknown() {
		return setField(fieldName, boolField.ValueBool())
	}

	return nil
}

// ReadIntField reads from a int struct field and return a sonarr field.
func ReadIntField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	intField := (*types.Int64)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !intField.IsNull() && !intField.IsUnknown() {
		return setField(fieldName, intField.ValueInt64())
	}

	return nil
}

// ReadFloatField reads from a float struct field and return a sonarr field.
func ReadFloatField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	floatField := (*types.Float64)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if !floatField.IsNull() && !floatField.IsUnknown() {
		return setField(fieldName, floatField.ValueFloat64())
	}

	return nil
}

// ReadStringSliceField reads from a string slice struct field and return a sonarr field.
func ReadStringSliceField(ctx context.Context, name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	sliceField := (*types.Set)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]string, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return setField(fieldName, slice)
	}

	return nil
}

// ReadIntSliceField reads from a int slice struct field and return a sonarr field.
func ReadIntSliceField(ctx context.Context, name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	sliceField := (*types.Set)(selectReadField(name, fieldCase).Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]int64, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return setField(fieldName, slice)
	}

	return nil
}
