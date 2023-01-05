package tools

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

func selectTFName(name string) string {
	for _, f := range getFieldExceptions() {
		if f.apiName == name {
			name = f.tfName
		}
	}

	return name
}

func selectAPIName(name string) string {
	for _, f := range getFieldExceptions() {
		if f.tfName == name {
			name = f.apiName
		}
	}

	return name
}

func WriteStringField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.GetName())
	stringValue := fmt.Sprint(fieldOutput.GetValue())
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.StringValue(stringValue))
	field.Set(v)
}

func WriteBoolField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.GetName())
	boolValue, _ := fieldOutput.GetValue().(bool)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.BoolValue(boolValue))
	field.Set(v)
}

func WriteIntField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.GetName())
	intValue, _ := fieldOutput.GetValue().(float64)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.Int64Value(int64(intValue)))
	field.Set(v)
}

func WriteFloatField(fieldOutput *sonarr.Field, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.GetName())
	floatValue, _ := fieldOutput.GetValue().(float64)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.Float64Value(floatValue))
	field.Set(v)
}

func WriteStringSliceField(ctx context.Context, fieldOutput *sonarr.Field, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.GetName())
	sliceValue, _ := fieldOutput.GetValue().([]interface{})
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	setValue := types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)

	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(setValue)
	field.Set(v)
}

func WriteIntSliceField(ctx context.Context, fieldOutput *sonarr.Field, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.GetValue().([]interface{})
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	setValue := types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)

	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldOutput.GetName()) })
	v := reflect.ValueOf(setValue)
	field.Set(v)
}

func ReadStringField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	stringField := (*types.String)(field.Addr().UnsafePointer())

	if !stringField.IsNull() && !stringField.IsUnknown() {
		field := sonarr.NewField()
		field.SetName(fieldName)
		field.SetValue(stringField.ValueString())

		return field
	}

	return nil
}

func ReadBoolField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	boolField := (*types.Bool)(field.Addr().UnsafePointer())

	if !boolField.IsNull() && !boolField.IsUnknown() {
		field := sonarr.NewField()
		field.SetName(fieldName)
		field.SetValue(boolField.ValueBool())

		return field
	}

	return nil
}

func ReadIntField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	intField := (*types.Int64)(field.Addr().UnsafePointer())

	if !intField.IsNull() && !intField.IsUnknown() {
		field := sonarr.NewField()
		field.SetName(fieldName)
		field.SetValue(intField.ValueInt64())

		return field
	}

	return nil
}

func ReadFloatField(name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	floatField := (*types.Float64)(field.Addr().UnsafePointer())

	if !floatField.IsNull() && !floatField.IsUnknown() {
		field := sonarr.NewField()
		field.SetName(fieldName)
		field.SetValue(floatField.ValueFloat64())

		return field
	}

	return nil
}

func ReadStringSliceField(ctx context.Context, name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	sliceField := (*types.Set)(field.Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]string, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		field := sonarr.NewField()
		field.SetName(fieldName)
		field.SetValue(slice)

		return field
	}

	return nil
}

func ReadIntSliceField(ctx context.Context, name string, fieldCase interface{}) *sonarr.Field {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	sliceField := (*types.Set)(field.Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]int64, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		field := sonarr.NewField()
		field.SetName(fieldName)
		field.SetValue(slice)

		return field
	}

	return nil
}
