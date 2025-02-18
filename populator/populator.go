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
	tt := tv.Type()
	if tv.Kind() != reflect.Ptr || tv.IsNil() {
		return ErrInvalidTargetType(tt.Name())
	}

	tv = tv.Elem()
	if tv.Kind() != reflect.Struct {
		return ErrInvalidTargetType(tt.Name())
	}

	//tt := tv.Type()

	for i := 0; i < tv.NumField(); i++ {
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

		if ok {
			if str, isStr := val.(string); isStr && str == "" && tag.def != "" {
				val = tag.def
			}
		} else {
			if tag.def != "" {
				val = tag.def
			} else if tag.req {
				return ErrRequiredValueIsMissing(tag.keys[0])
			} else {
				continue
			}
		}

		if err := set(fieldValue, val); err != nil {
			return fmt.Errorf("failed to set value for field %q: %w", fieldType.Name, err)
		}

		delete(in, fieldTag)
	}

	return nil
}

func set(f reflect.Value, value any) error {
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
		var err error
		switch targetElem.Kind() {
		case reflect.Int:
			var i int
			i, err = strconv.Atoi(v)
			if err != nil {
				return err
			}
			targetElem.SetInt(int64(i))
		case reflect.Int64:
			if targetElem.Type() == reflect.TypeOf(time.Duration(0)) {
				var d time.Duration
				d, err = time.ParseDuration(v)
				if err != nil {
					return err
				}
				targetElem.SetInt(int64(d))
			} else {
				var i int64
				i, err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				targetElem.SetInt(i)
			}
		case reflect.Float32:
			var f float64
			f, err = strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}
			targetElem.SetFloat(f)
		case reflect.Float64:
			var f float64
			f, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			targetElem.SetFloat(f)
		case reflect.Bool:
			var b bool
			b, err = strconv.ParseBool(v)
			if err != nil {
				return err
			}
			targetElem.SetBool(b)
		case reflect.Struct:
			if targetElem.Type() == reflect.TypeOf(time.Time{}) {
				var t time.Time
				t, err = time.Parse(time.RFC3339, v)
				if err != nil {
					return err
				}
				targetElem.Set(reflect.ValueOf(t))
			} else {
				return fmt.Errorf("unsupported target struct type: %T", targetElem.Interface())
			}
		case reflect.String:
			targetElem.SetString(v)
		default:
			return fmt.Errorf("unsupported target type: %T", targetElem.Interface())
		}
		return nil

	case bool:
		if targetElem.Kind() == reflect.Bool {
			targetElem.SetBool(v)
			return nil
		}
		return fmt.Errorf("unsupported target type: %T", targetElem.Interface())

	case int, int64, float32, float64:
		if targetElem.Type() == reflect.TypeOf(time.Duration(0)) {
			var d time.Duration
			switch num := v.(type) {
			case int:
				d = time.Duration(num)
			case int64:
				d = time.Duration(num)
			case float32:
				d = time.Duration(num)
			case float64:
				d = time.Duration(num)
			default:
				return fmt.Errorf("unsupported numeric type for time.Duration: %T", v)
			}
			targetElem.SetInt(int64(d))
			return nil
		}
		switch targetElem.Kind() {
		case reflect.Int:
			switch val := v.(type) {
			case int:
				targetElem.SetInt(int64(val))
			case int64:
				targetElem.SetInt(val)
			case float32:
				targetElem.SetInt(int64(val))
			case float64:
				targetElem.SetInt(int64(val))
			}
		case reflect.Int64:
			switch val := v.(type) {
			case int:
				targetElem.SetInt(int64(val))
			case int64:
				targetElem.SetInt(val)
			case float32:
				targetElem.SetInt(int64(val))
			case float64:
				targetElem.SetInt(int64(val))
			}
		case reflect.Float32:
			switch val := v.(type) {
			case int:
				targetElem.SetFloat(float64(val))
			case int64:
				targetElem.SetFloat(float64(val))
			case float32:
				targetElem.SetFloat(float64(val))
			case float64:
				targetElem.SetFloat(val)
			}
		case reflect.Float64:
			switch val := v.(type) {
			case int:
				targetElem.SetFloat(float64(val))
			case int64:
				targetElem.SetFloat(float64(val))
			case float32:
				targetElem.SetFloat(float64(val))
			case float64:
				targetElem.SetFloat(val)
			}
		default:
			return fmt.Errorf("unsupported target type: %T", targetElem.Interface())
		}
		return nil

	default:
		return fmt.Errorf("unsupported source type: %T", source)
	}
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
