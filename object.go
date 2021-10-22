package reskit

import "context"

type Object interface {
	Destroy(ctx context.Context) error
	Get() interface{}
}

type object struct {
	obj     interface{}
	destroy func(ctx context.Context) error
}

func (s *object) Destroy(ctx context.Context) error {
	return s.destroy(ctx)
}

func (s *object) Get() interface{} {
	return s.obj
}

func BuildObject(obj interface{}, destroy func(ctx context.Context) error) Object {
	return &object{
		obj:     obj,
		destroy: destroy,
	}
}
