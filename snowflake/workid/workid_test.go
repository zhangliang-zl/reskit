package workid

import (
	"testing"
)

func TestFactory_Get(t *testing.T) {
	fac := NewFactory("t1:111@tcp(localhost:3306)/t1?charset=utf8&parseTime=true&loc=Local")
	var firstGetID int64
	for i := 0; i < 10; i++ {
		id, err := fac.Get()
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			firstGetID = id
		}
		if firstGetID != id {
			t.Error("The IDs obtained are not equal")
		}
	}

}
