package populator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func PopulatePluginConfig(configStore any, configName string, target any) error {
	cs, ok := configStore.(map[string]any)
	if !ok {
		return ErrWrongInputConfigStoreType(configStore, map[string]any{})
	}
	c, ok := cs[configName]
	if !ok {
		return ErrSpecifiedConfigIsMissing(configName)
	}

	return Populate(c, target)
}

func Populate(source, target any) error {
	in, ok := source.(map[string]any)
	if !ok {
		return ErrWrongInputDataType(source, map[string]any{})
	}

	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Ptr || tv.IsNil() {
		return ErrInvalidTargetType
	}

	tv = tv.Elem()
	if tv.Kind() != reflect.Struct {
		return ErrInvalidTargetType
	}

	tt := tv.Type()

	for i := range tv.NumField() {
		fieldValue := tv.Field(i)
		if fieldValue.Kind() == reflect.Struct {
			if !fieldValue.Addr().CanInterface() {
				continue
			}
			if err := Populate(source, fieldValue.Addr().Interface()); err != nil {
				return err
			}
		}

		fieldType := tt.Field(i)
		fieldTag := fieldType.Tag.Get("krakend")
		if fieldTag == "" {
			continue
		}

		if !fieldValue.CanSet() {
			return ErrFieldIsPrivate(fieldType.Name)
		}

		tag, err := parseTag(fieldTag)
		if err != nil {
			return err
		}

		var (
			val any
			ok  bool
		)

		for _, key := range tag.keys {
			val, ok = in[key]
			if ok {
				break
			}
		}

		if !ok {
			switch {
			case tag.def == "":
				val = tag.def
			case tag.req:
				return ErrRequiredValueIsMissing(tag.keys[0])
			default:
				continue
			}
		}

		if err := set(fieldType.Type, fieldValue, val); err != nil {
			return err
		}

		delete(in, fieldTag)
	}

	return nil
}

func set(t reflect.Type, f reflect.Value, value any) error {
	if !f.CanSet() {
		return errors.New("cannot set field")
	}

	return transform(value, f.Addr().Interface())
}

func transform(source any, target any) error {
	if target == nil {
		return errors.New("target cannot be nil")
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return errors.New("target must be a non-nil pointer")
	}
	targetElem := targetValue.Elem()
	switch v := source.(type) {
	case string:
		var result any
		var err error
		switch targetElem.Kind() {
		case reflect.Int:
			result, err = strconv.Atoi(v)
		case reflect.Int64:
			if targetElem.Type() == reflect.TypeOf(time.Duration(0)) {
				result, err = time.ParseDuration(v)
			} else {
				result, err = strconv.ParseInt(v, 10, 64)
			}
		case reflect.Float32:
			var f float64
			f, err = strconv.ParseFloat(v, 32)
			result = float32(f)
		case reflect.Float64:
			result, err = strconv.ParseFloat(v, 64)
		case reflect.Struct:
			if targetElem.Type() == reflect.TypeOf(time.Time{}) {
				result, err = time.Parse(time.RFC3339, v)
			} else {
				return fmt.Errorf("unsupported target struct type: %T", targetElem.Interface())
			}
		case reflect.String:
			result = v
		default:
			return fmt.Errorf("unsupported target type: %T", targetElem.Interface())
		}
		if err != nil {
			return err
		}
		targetElem.Set(reflect.ValueOf(result))
		return nil

	case int, int64, float32, float64:
		switch targetElem.Kind() {
		case reflect.Int64:
			targetElem.SetInt(v.(int64))
		case reflect.Int:
			targetElem.SetInt(int64(v.(int)))
		case reflect.Float32:
			targetElem.SetFloat(float64(v.(float32)))
		case reflect.Float64:
			targetElem.SetFloat(v.(float64))
		default:
			panic("unhandled default case")
		}
		return nil

	default:
		return fmt.Errorf("unsupported source type: %T", source)
	}

	return nil
}

type tag struct {
	keys []string
	def  string
	req  bool
}

func parseTag(tagString string) (tag, error) {
	var t tag
	keys := strings.Split(tagString, ",")
	for _, key := range keys {
		if !strings.Contains(key, "=") {
			t.keys = append(t.keys, key)
			continue
		}
		data := strings.SplitN(key, "=", 2)
		switch strings.ToLower(strings.TrimSpace(data[0])) {
		case "default":
			t.def = data[1]
		case "required":
			bv, _ := strconv.ParseBool(data[1])
			t.req = bv
		default:
			continue
		}
	}

	return t, nil
}
