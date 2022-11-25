package tools

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr"
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

func WriteStringField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.Name)
	stringValue := fmt.Sprint(fieldOutput.Value)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.StringValue(stringValue))
	field.Set(v)
}

func WriteBoolField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.Name)
	boolValue, _ := fieldOutput.Value.(bool)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.BoolValue(boolValue))
	field.Set(v)
}

func WriteIntField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.Name)
	intValue, _ := fieldOutput.Value.(float64)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.Int64Value(int64(intValue)))
	field.Set(v)
}

func WriteFloatField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.Name)
	floatValue, _ := fieldOutput.Value.(float64)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(types.Float64Value(floatValue))
	field.Set(v)
}

func WriteStringSliceField(ctx context.Context, fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	fieldName := selectTFName(fieldOutput.Name)
	sliceValue, _ := fieldOutput.Value.([]interface{})
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	setValue := types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)

	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldName) })
	v := reflect.ValueOf(setValue)
	field.Set(v)
}

func WriteIntSliceField(ctx context.Context, fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.Value.([]interface{})
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	setValue := types.SetValueMust(types.Int64Type, nil)
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)

	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, fieldOutput.Name) })
	v := reflect.ValueOf(setValue)
	field.Set(v)
}

func ReadStringField(name string, fieldCase interface{}) *starr.FieldInput {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	stringField := (*types.String)(field.Addr().UnsafePointer())

	if !stringField.IsNull() && !stringField.IsUnknown() {
		return &starr.FieldInput{
			Name:  fieldName,
			Value: stringField.ValueString(),
		}
	}

	return nil
}

func ReadBoolField(name string, fieldCase interface{}) *starr.FieldInput {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	boolField := (*types.Bool)(field.Addr().UnsafePointer())

	if !boolField.IsNull() && !boolField.IsUnknown() {
		return &starr.FieldInput{
			Name:  fieldName,
			Value: boolField.ValueBool(),
		}
	}

	return nil
}

func ReadIntField(name string, fieldCase interface{}) *starr.FieldInput {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	intField := (*types.Int64)(field.Addr().UnsafePointer())

	if !intField.IsNull() && !intField.IsUnknown() {
		return &starr.FieldInput{
			Name:  fieldName,
			Value: intField.ValueInt64(),
		}
	}

	return nil
}

func ReadFloatField(name string, fieldCase interface{}) *starr.FieldInput {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	intField := (*types.Float64)(field.Addr().UnsafePointer())

	if !intField.IsNull() && !intField.IsUnknown() {
		return &starr.FieldInput{
			Name:  fieldName,
			Value: intField.ValueFloat64(),
		}
	}

	return nil
}

func ReadStringSliceField(ctx context.Context, name string, fieldCase interface{}) *starr.FieldInput {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	sliceField := (*types.Set)(field.Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]string, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return &starr.FieldInput{
			Name:  fieldName,
			Value: slice,
		}
	}

	return nil
}

func ReadIntSliceField(ctx context.Context, name string, fieldCase interface{}) *starr.FieldInput {
	fieldName := selectAPIName(name)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
	sliceField := (*types.Set)(field.Addr().UnsafePointer())

	if len(sliceField.Elements()) != 0 {
		slice := make([]int64, len(sliceField.Elements()))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return &starr.FieldInput{
			Name:  fieldName,
			Value: slice,
		}
	}

	return nil
}
