package tf

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type Set struct {
	values []interface{}
}

func (set *Set) List() []interface{} { return set.values }
func (set *Set) Len() int            { return len(set.values) }

type Data struct {
	typ    tftypes.Object
	values map[string]tftypes.Value
	source map[string]tftypes.Value
	shape  map[string]valueShape
	prior  map[string]tftypes.Value
	errors diag.Diagnostics
}

type valueShape struct {
	computed bool
	fields   map[string]valueShape
}

func NewData(raw tftypes.Value) (*Data, error) {
	typ, ok := raw.Type().(tftypes.Object)
	if !ok {
		return nil, fmt.Errorf("expected Terraform object, got %s", raw.Type())
	}
	data := &Data{typ: typ, values: map[string]tftypes.Value{}}
	if !raw.IsNull() && raw.IsKnown() {
		if err := raw.As(&data.values); err != nil {
			return nil, err
		}
	}
	for name, attrType := range typ.AttributeTypes {
		if _, exists := data.values[name]; !exists {
			data.values[name] = tftypes.NewValue(attrType, nil)
		}
	}
	data.source = make(map[string]tftypes.Value, len(data.values))
	for name, value := range data.values {
		data.source[name] = value
	}
	return data, nil
}

func (data *Data) Get(name string) interface{} {
	value, exists := data.values[name]
	if !exists {
		return nil
	}
	return decode(value)
}

func (data *Data) GetOk(name string) (interface{}, bool) {
	value := data.Get(name)
	if value == nil {
		return nil, false
	}
	if set, ok := value.(*Set); ok {
		return value, set.Len() > 0
	}
	reflected := reflect.ValueOf(value)
	if reflected.Kind() == reflect.Slice || reflected.Kind() == reflect.Map {
		return value, reflected.Len() > 0
	}
	return value, !reflected.IsZero()
}

func (data *Data) NewValueKnown(name string) bool {
	value, exists := data.values[name]
	return exists && value.IsKnown()
}

func (data *Data) Configured(name string) bool {
	value, exists := data.values[name]
	return exists && value.IsKnown() && !value.IsNull()
}

func (data *Data) GetOkExists(name string) (interface{}, bool) {
	return data.Get(name), data.Configured(name)
}

func (data *Data) HasChanges(names ...string) bool {
	for _, name := range names {
		previous, exists := data.prior[name]
		if !exists || !previous.Equal(data.values[name]) {
			return true
		}
	}
	return false
}

func (data *Data) Set(name string, value interface{}) error {
	typ, exists := data.typ.AttributeTypes[name]
	if !exists {
		err := fmt.Errorf("unknown Terraform attribute %q", name)
		data.errors.Append(FromErr(err)...)
		return err
	}
	encoded, err := encode(typ, value)
	if err != nil {
		err = fmt.Errorf("setting %s: %w", name, err)
		data.errors.Append(FromErr(err)...)
		return err
	}
	if shape, exists := data.shape[name]; exists {
		encoded = preserveEmpty(data.source[name], encoded, shape)
	}
	data.values[name] = encoded
	return nil
}

func (data *Data) SetId(id string) { _ = data.Set("id", id) }
func (data *Data) Id() string {
	value, _ := data.Get("id").(string)
	return value
}

func (data *Data) Raw() tftypes.Value {
	return tftypes.NewValue(data.typ, data.values)
}

func (data *Data) Complete() tftypes.Value {
	return knownOrNull(data.Raw())
}

func (data *Data) Diagnostics() diag.Diagnostics { return data.errors }

func FromErr(err error) diag.Diagnostics {
	if err == nil {
		return nil
	}
	return diag.Diagnostics{diag.NewErrorDiagnostic("Veeam Provider Error", err.Error())}
}

func Errorf(format string, args ...interface{}) diag.Diagnostics {
	return FromErr(fmt.Errorf(format, args...))
}

func decode(value tftypes.Value) interface{} {
	typ := value.Type()
	if typ.Is(tftypes.String) {
		var result string
		if value.IsKnown() && !value.IsNull() {
			_ = value.As(&result)
		}
		return result
	}
	if typ.Is(tftypes.Bool) {
		var result bool
		if value.IsKnown() && !value.IsNull() {
			_ = value.As(&result)
		}
		return result
	}
	if typ.Is(tftypes.Number) {
		var result big.Float
		if value.IsKnown() && !value.IsNull() {
			_ = value.As(&result)
		}
		integer, _ := result.Int64()
		return int(integer)
	}
	switch typ.(type) {
	case tftypes.List, tftypes.Set:
		var elements []tftypes.Value
		if value.IsKnown() && !value.IsNull() {
			_ = value.As(&elements)
		}
		result := make([]interface{}, len(elements))
		for index, element := range elements {
			result[index] = decode(element)
		}
		if _, isSet := typ.(tftypes.Set); isSet {
			return &Set{values: result}
		}
		return result
	case tftypes.Map, tftypes.Object:
		var elements map[string]tftypes.Value
		if value.IsKnown() && !value.IsNull() {
			_ = value.As(&elements)
		}
		result := map[string]interface{}{}
		for name, element := range elements {
			result[name] = decode(element)
		}
		return result
	}
	return nil
}

func encode(typ tftypes.Type, input interface{}) (tftypes.Value, error) {
	if input == nil {
		return tftypes.NewValue(typ, nil), nil
	}
	if set, ok := input.(*Set); ok {
		input = set.List()
	}
	value := reflect.ValueOf(input)
	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return tftypes.NewValue(typ, nil), nil
		}
		value = value.Elem()
	}
	if (value.Kind() == reflect.Slice || value.Kind() == reflect.Map) && value.IsNil() {
		return tftypes.NewValue(typ, nil), nil
	}
	if typ.Is(tftypes.String) && value.Kind() == reflect.String {
		return tftypes.NewValue(typ, value.String()), nil
	}
	if typ.Is(tftypes.Bool) && value.Kind() == reflect.Bool {
		return tftypes.NewValue(typ, value.Bool()), nil
	}
	if typ.Is(tftypes.Number) {
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return tftypes.NewValue(typ, value.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return tftypes.NewValue(typ, value.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return tftypes.NewValue(typ, value.Float()), nil
		}
	}
	var elementType tftypes.Type
	switch collection := typ.(type) {
	case tftypes.List:
		elementType = collection.ElementType
	case tftypes.Set:
		elementType = collection.ElementType
	case tftypes.Map:
		elementType = collection.ElementType
	}
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		if elementType == nil {
			return tftypes.Value{}, fmt.Errorf("cannot store collection in %s", typ)
		}
		elements := make([]tftypes.Value, value.Len())
		for index := range elements {
			var err error
			elements[index], err = encode(elementType, value.Index(index).Interface())
			if err != nil {
				return tftypes.Value{}, err
			}
		}
		return tftypes.NewValue(typ, elements), nil
	}
	if value.Kind() == reflect.Map {
		elements := map[string]tftypes.Value{}
		if object, ok := typ.(tftypes.Object); ok {
			for name, attrType := range object.AttributeTypes {
				entry := value.MapIndex(reflect.ValueOf(name))
				var item interface{}
				if entry.IsValid() {
					item = entry.Interface()
				}
				encoded, err := encode(attrType, item)
				if err != nil {
					return tftypes.Value{}, fmt.Errorf("%s: %w", name, err)
				}
				elements[name] = encoded
			}
		} else if elementType != nil {
			for _, key := range value.MapKeys() {
				encoded, err := encode(elementType, value.MapIndex(key).Interface())
				if err != nil {
					return tftypes.Value{}, err
				}
				elements[key.String()] = encoded
			}
		} else {
			return tftypes.Value{}, fmt.Errorf("cannot store map in %s", typ)
		}
		return tftypes.NewValue(typ, elements), nil
	}
	return tftypes.Value{}, fmt.Errorf("cannot store %T in %s", input, typ)
}

func knownOrNull(value tftypes.Value) tftypes.Value {
	if !value.IsKnown() {
		return tftypes.NewValue(value.Type(), nil)
	}
	if value.IsNull() {
		return value
	}
	switch value.Type().(type) {
	case tftypes.Object, tftypes.Map:
		var elements map[string]tftypes.Value
		_ = value.As(&elements)
		for name, element := range elements {
			elements[name] = knownOrNull(element)
		}
		return tftypes.NewValue(value.Type(), elements)
	case tftypes.List, tftypes.Set:
		var elements []tftypes.Value
		_ = value.As(&elements)
		for index, element := range elements {
			elements[index] = knownOrNull(element)
		}
		return tftypes.NewValue(value.Type(), elements)
	}
	return value
}

func emptyValue(value tftypes.Value) bool {
	if value.IsNull() {
		return true
	}
	if !value.IsKnown() {
		return false
	}
	decoded := decode(value)
	if set, ok := decoded.(*Set); ok {
		return set.Len() == 0
	}
	reflected := reflect.ValueOf(decoded)
	if !reflected.IsValid() {
		return true
	}
	if reflected.Kind() == reflect.Slice || reflected.Kind() == reflect.Map {
		return reflected.Len() == 0
	}
	return reflected.IsZero()
}

func preserveEmpty(previous, next tftypes.Value, shape valueShape) tftypes.Value {
	if shape.computed || !previous.IsKnown() {
		return next
	}
	if previous.IsNull() && emptyValue(next) {
		return previous
	}
	if next.IsNull() && emptyValue(previous) {
		return previous
	}
	if previous.IsNull() || next.IsNull() || !next.IsKnown() {
		return next
	}
	switch previous.Type().(type) {
	case tftypes.Object:
		var before, after map[string]tftypes.Value
		_ = previous.As(&before)
		_ = next.As(&after)
		for name, candidate := range after {
			after[name] = preserveEmpty(before[name], candidate, shape.fields[name])
		}
		return tftypes.NewValue(next.Type(), after)
	case tftypes.List:
		var before, after []tftypes.Value
		_ = previous.As(&before)
		_ = next.As(&after)
		for index := range after {
			if index < len(before) {
				after[index] = preserveEmpty(before[index], after[index], shape)
			}
		}
		return tftypes.NewValue(next.Type(), after)
	}
	return next
}
