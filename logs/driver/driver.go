package driver

type Writer interface {
	Record(m string)
}

type WriterBuild func(tag string) (Writer, error)
