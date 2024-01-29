package goregexpopulate

import (
	"fmt"
	"reflect"

	regen "github.com/zach-klippenstein/goregen"
)

type populate struct {
	genMap      map[string]regen.Generator
	entropyFunc func(bool) bool
}

func New(entropyFunc func(bool) bool) *populate {
	newGenCache()

	p := &populate{
		genMap: make(map[string]regen.Generator),
	}

	p.entropyFunc = entropyDefault
	if entropyFunc != nil {
		p.entropyFunc = entropyFunc
	}

	return p
}

func (p *populate) Populate(a any) error {
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)

	if v.Kind() != reflect.Pointer {
		return ErrNonPointer
	}

	if v.IsNil() {
		return ErrNilPointer
	}

	for i := 0; i < v.Elem().NumField(); i++ {
		structField, _ := t.Elem().FieldByName(v.Elem().Type().Field(i).Name)
		required, requiredFound := structField.Tag.Lookup("required")
		if !requiredFound {
			continue
		}

		if !p.entropyFunc(required == "true") {
			continue
		}

		field := v.Elem().Field(i)
		switch field.Kind() {
		case reflect.String:
			pattern, patternFound := structField.Tag.Lookup("pattern")
			if !patternFound {
				continue
			}

			generator, err := cache.register(pattern)
			if err != nil {
				return fmt.Errorf(
					"unable to register generator for pattern %q from struct field %q with cache instance: %w",
					pattern, structField.Name, err,
				)
			}

			if field.IsZero() && field.CanSet() {
				field.SetString(generator.Generate())
			}

		case reflect.Struct:
			if err := p.Populate(field.Addr().Interface()); err != nil {
				return fmt.Errorf("unable to populate struct container %q: %w", structField.Name, err)
			}

		case reflect.Pointer:
			if field.IsNil() {
				fieldPointerNil := field.Type().Elem()
				fieldTypeInstance := reflect.New(fieldPointerNil).Elem()
				fieldPointerInstance := reflect.New(reflect.Zero(fieldTypeInstance.Type()).Type())
				field.Set(fieldPointerInstance)
			}

			if err := p.Populate(field.Interface()); err != nil {
				return fmt.Errorf("unable to populate pointer container %q: %w", structField.Name, err)
			}
		}
	}

	return nil
}
