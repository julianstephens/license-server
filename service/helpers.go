package service

import "reflect"

type Item struct {
	Key   string
	Value any
}

func StructItems[T any](obj T) []Item {
	values := reflect.ValueOf(obj)
	types := values.Type()

	var res []Item
	for i := 0; i < values.NumField(); i++ {
		k := types.Field(i).Name
		val := values.Field(i)
		res = append(res, Item{Key: k, Value: val})
	}

	return res
}

func SetProperty[T any](obj T, propName string, propValue any) *T {
	reflect.ValueOf(obj).Elem().FieldByName(propName).Set(reflect.ValueOf(propValue))
	return &obj
}

func Unpack[T any](obj T) interface{} {
	s := reflect.ValueOf(obj)
	ret := make([]interface{}, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		ret[i] = s.Field(i).Interface()
	}

	return ret
}
