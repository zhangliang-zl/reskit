package component

type build struct {
	run   func() error
	close func() error
	obj   interface{}
}

func (d *build) Run() error {
	if d.run == nil {
		return nil
	}

	return d.run()
}

func (d *build) Close() error {
	if d.close == nil {
		return nil
	}

	return d.close()
}

func (d *build) Instance() interface{} {
	return d.obj
}

func Make(obj interface{}, run, close func() error) Interface {
	return &build{
		obj:   obj,
		run:   run,
		close: close,
	}
}
