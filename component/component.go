package component

import "fmt"

type Interface interface {
	Run() error
	Close() error
	Instance() interface{}
}

func BuildKey(cType string, id string) string {
	return fmt.Sprintf("%s%s", id, cType)
}
