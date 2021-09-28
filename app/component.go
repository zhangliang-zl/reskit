package app

type Component interface {
	Object() interface{}
	Close() error
}

type componentIns struct {
	closeFunc func() error
	object    interface{}
}

func (c componentIns) Object() interface{} {
	return c.object
}

func (c componentIns) Close() error {
	return c.closeFunc()
}

func MakeComponent(object interface{}, closeFunc func() error) Component {
	return componentIns{
		object:    object,
		closeFunc: closeFunc,
	}
}
