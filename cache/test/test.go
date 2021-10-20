package test

import (
	"context"
	"github.com/zhangliang-zl/reskit/cache"
	"testing"
	"time"

)

func AllCase(cache cache.Cache, t *testing.T) {
	builtinTest(cache, t)
	getOrSet(cache, t)
	delete(cache, t)
}

type structObj struct {
	ID   string
	Name string
}

func callbackFunc(id string) (structObj, error) {
	return structObj{
		id,
		"Bob",
	}, nil
}

func builtinTest(cache cache.Cache, t *testing.T) {
	ctx := context.TODO()
	boolVal(cache, ctx, t)
	stringVal(cache, ctx, t)
	intVal(cache, ctx, t)
	mapVal(cache, ctx, t)
	structObjVal(cache, ctx, t)
	complexVal(cache, ctx, t)
	sliceVal(cache, ctx, t)
}

func boolVal(cache cache.Cache, ctx context.Context, t *testing.T) {
	var outBool bool
	if err := cache.Set(ctx, "boolVal", true, 3*time.Second); err != nil {
		t.Error(err)
	}
	if _, err := cache.Get(ctx, "boolVal", &outBool); err != nil {
		t.Error(err)
	}
}

func stringVal(cache cache.Cache, ctx context.Context, t *testing.T) {
	var out string
	if err := cache.Set(ctx, "test1", "test1", 3*time.Second); err != nil {
		t.Error(err)
	}
	if _, err := cache.Get(ctx, "test1", &out); err != nil {
		t.Error(err)
	}
}

func intVal(cache cache.Cache, ctx context.Context, t *testing.T) {
	var intVal int64
	if err := cache.Set(ctx, "intVal", 3, 3*time.Second); err != nil {
		t.Error(err)
	}
	if _, err := cache.Get(ctx, "intVal", &intVal); err != nil {
		t.Error(err)
	}
}

func mapVal(cache cache.Cache, ctx context.Context, t *testing.T) {
	var mapVal map[string]int
	if err := cache.Set(ctx, "mapVal", map[string]int{"test1": 1}, 3*time.Second); err != nil {
		t.Error(err)
	}
	if _, err := cache.Get(ctx, "mapVal", &mapVal); err != nil {
		t.Error(err)
	}
	if v, _ := mapVal["test1"]; v != 1 {
		t.Errorf("mapVal error ")
	}
}

func structObjVal(cache cache.Cache, ctx context.Context, t *testing.T) {
	var structVal structObj
	if err := cache.Set(ctx, "structVal", structObj{"1", "Bob"}, 3*time.Second); err != nil {
		t.Error(err)
	}
	if _, err := cache.Get(ctx, "structVal", &structVal); err != nil {
		t.Error(err)
	}
	if structVal.Name != "Bob" {
		t.Errorf("structVal error ")
	}
}

func complexVal(cache cache.Cache, ctx context.Context, t *testing.T) {
	var complexVal map[string][]structObj
	if err := cache.Set(ctx, "complexVal", map[string][]structObj{
		"key": {structObj{"111", "Bob"}},
	}, 3*time.Second); err != nil {
		t.Error(err)
	}
	if _, err := cache.Get(ctx, "complexVal", &complexVal); err != nil {
		t.Error(err)
	}
	if complexVal["key"][0].Name != "Bob" {
		t.Errorf("complexVal error ")
	}
}

func sliceVal(cache cache.Cache, ctx context.Context, t *testing.T) {
	var sliceVal []string
	if err := cache.Set(ctx, "sliceVal", []string{"111", "Bob"}, 3*time.Second); err != nil {
		t.Error(err)
	}
	if _, err := cache.Get(ctx, "sliceVal", &sliceVal); err != nil {
		t.Error(err)
	}
	if len(sliceVal) == 2 && sliceVal[0] != "111" {
		t.Errorf("sliceVal error ")
	}
}

const (
	callbackFunc1Key = "callbackFunc:1"
	delete1          = "delete:1"
)

func getOrSet(cache cache.Cache, t *testing.T) {
	var id = "abc"
	var val1 structObj
	var val2 structObj

	callback := func() (interface{}, error) {
		return callbackFunc(id)
	}

	ctx := context.TODO()
	err := cache.GetOrSet(ctx, callbackFunc1Key, &val1, time.Second*10, callback)
	if err != nil {
		t.Error(err)
	}

	_, err = cache.Get(ctx, callbackFunc1Key, &val2)
	if err != nil {
		t.Error(err)
	}

	if val1.ID != id {
		t.Errorf("ID Error")
	}

	if val1.ID != val2.ID {
		t.Error()
	}
}

func delete(cache cache.Cache, t *testing.T) {
	ctx := context.TODO()
	err := cache.Set(ctx, delete1, "111", time.Second*10)
	if err != nil {
		t.Error(err)
	}

	var s string
	_, err = cache.Get(ctx, delete1, &s)
	if err != nil {
		t.Error(err)
	}

	if s != "111" {
		t.Error("Get error ")
	}

	if err := cache.Delete(ctx, delete1); err != nil {
		t.Error(err)
	}

	var ss string
	exist, err := cache.Get(ctx, delete1, &ss)
	if err != nil {
		t.Error(err)
	}

	if exist {
		t.Error("exist after Delete()")
	}
}
