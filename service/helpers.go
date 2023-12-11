package service

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
)

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

func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func GetId(ctx *gin.Context) (string, error) {
	id := ctx.Param("id")
	if id == "" {
		return id, errors.New("no resource id provided")
	}
	return id, nil
}
