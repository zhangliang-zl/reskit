package component

type Interface interface {
	Instance() interface{}
	Close() error
}

type buildIns struct {
	close    func() error
	instance interface{}
}

func (c buildIns) Instance() interface{} {
	return c.instance
}

func (c buildIns) Close() error {
	return c.close()
}

func Make(instance interface{}, closeFunc func() error) Interface {
	return buildIns{
		instance: instance,
		close:    closeFunc,
	}
}
