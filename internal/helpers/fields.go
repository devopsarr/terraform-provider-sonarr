package helpers

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr"
)

func WriteStringField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	stringValue, _ := fieldOutput.Value.(string)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(fieldOutput.Name)) })
	v := reflect.ValueOf(types.String{Value: stringValue})
	field.Set(v)
}

func WriteBoolField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	boolValue, _ := fieldOutput.Value.(bool)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(fieldOutput.Name)) })
	v := reflect.ValueOf(types.Bool{Value: boolValue})
	field.Set(v)
}

func WriteIntField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	intValue, _ := fieldOutput.Value.(float64)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(fieldOutput.Name)) })
	v := reflect.ValueOf(types.Int64{Value: int64(intValue)})
	field.Set(v)
}

func WriteFloatField(fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	floatValue, _ := fieldOutput.Value.(float64)
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(fieldOutput.Name)) })
	v := reflect.ValueOf(types.Float64{Value: floatValue})
	field.Set(v)
}

func WriteStringSliceField(ctx context.Context, fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.Value.([]interface{})
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	setValue := types.Set{ElemType: types.StringType}
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)

	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(fieldOutput.Name)) })
	v := reflect.ValueOf(setValue)
	field.Set(v)
}

func WriteIntSliceField(ctx context.Context, fieldOutput *starr.FieldOutput, fieldCase interface{}) {
	sliceValue, _ := fieldOutput.Value.([]interface{})
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	setValue := types.Set{ElemType: types.Int64Type}
	tfsdk.ValueFrom(ctx, sliceValue, setValue.Type(ctx), &setValue)

	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(fieldOutput.Name)) })
	v := reflect.ValueOf(setValue)
	field.Set(v)
}

func ReadStringField(name string, fieldCase interface{}) *starr.FieldInput {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(name)) })
	stringField := (*types.String)(field.Addr().UnsafePointer())

	if !stringField.IsNull() && !stringField.IsUnknown() {
		return &starr.FieldInput{
			Name:  name,
			Value: stringField.Value,
		}
	}

	return nil
}

func ReadBoolField(name string, fieldCase interface{}) *starr.FieldInput {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(name)) })
	boolField := (*types.Bool)(field.Addr().UnsafePointer())

	if !boolField.IsNull() && !boolField.IsUnknown() {
		return &starr.FieldInput{
			Name:  name,
			Value: boolField.Value,
		}
	}

	return nil
}

func ReadIntField(name string, fieldCase interface{}) *starr.FieldInput {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(name)) })
	intField := (*types.Int64)(field.Addr().UnsafePointer())

	if !intField.IsNull() && !intField.IsUnknown() {
		return &starr.FieldInput{
			Name:  name,
			Value: intField.Value,
		}
	}

	return nil
}

func ReadFloatField(name string, fieldCase interface{}) *starr.FieldInput {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(name)) })
	intField := (*types.Float64)(field.Addr().UnsafePointer())

	if !intField.IsNull() && !intField.IsUnknown() {
		return &starr.FieldInput{
			Name:  name,
			Value: intField.Value,
		}
	}

	return nil
}

func ReadStringSliceField(ctx context.Context, name string, fieldCase interface{}) *starr.FieldInput {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(name)) })
	sliceField := (*types.Set)(field.Addr().UnsafePointer())

	if len(sliceField.Elems) != 0 {
		slice := make([]string, len(sliceField.Elems))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return &starr.FieldInput{
			Name:  name,
			Value: slice,
		}
	}

	return nil
}

func ReadIntSliceField(ctx context.Context, name string, fieldCase interface{}) *starr.FieldInput {
	value := reflect.ValueOf(fieldCase)
	value = value.Elem()
	field := value.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, strings.ToLower(name)) })
	sliceField := (*types.Set)(field.Addr().UnsafePointer())

	if len(sliceField.Elems) != 0 {
		slice := make([]int64, len(sliceField.Elems))
		tfsdk.ValueAs(ctx, sliceField, &slice)

		return &starr.FieldInput{
			Name:  name,
			Value: slice,
		}
	}

	return nil
}
