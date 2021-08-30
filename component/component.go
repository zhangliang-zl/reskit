package component

import "fmt"

type Interface interface {
	Run() error
	Close() error
	Instance() interface{}
}

const DefaultID = "default"

func BuildKey(cType string, id string) string {
	var trueID string
	if id == "" {
		trueID = DefaultID
	} else {
		trueID = id
	}

	return fmt.Sprintf("%s%s", trueID, cType)
}
