package cmd

import (
	"reflect"
	"strings"
)

func newRegister() *register {
	return &register{
		funcs: map[string]rf{},
	}
}

type register struct {
	funcs map[string]rf
}

type rf func() error

func (r *register) registerAPICompare(ac *apiCompare) error {
	rv := reflect.ValueOf(ac)
	rt := reflect.TypeOf(ac)

	for i := 0; i < rv.NumMethod(); i++ {
		name := rt.Method(i).Name
		m := rv.MethodByName(name)
		if !strings.HasPrefix(name, methodPrefix) {
			continue
		}

		name = strings.TrimPrefix(name, methodPrefix)
		r.funcs[name] = func() error {
			res := m.Call([]reflect.Value{})
			if res[0].Interface() == nil {
				return nil
			}

			return res[0].Interface().(error)
		}
	}

	return nil
}
