package driver

import "io"

type WriterBuild func(tag string) (io.Writer, error)
