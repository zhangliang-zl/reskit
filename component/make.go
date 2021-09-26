package component

type componentIns struct {
	run   func() error
	close func() error
	obj   interface{}
}

func (d *componentIns) Run() error {
	if d.run == nil {
		return nil
	}

	return d.run()
}

func (d *componentIns) Close() error {
	if d.close == nil {
		return nil
	}

	return d.close()
}

func (d *componentIns) Instance() interface{} {
	return d.obj
}

func Make(obj interface{}, run, close func() error) Interface {
	return &componentIns{
		obj:   obj,
		run:   run,
		close: close,
	}
}
