package snowflake

import (
	"testing"
)

func TestNewWorkerIDFactory(t *testing.T) {
	dsn := "user:password@tcp(localhost:3306)/db1"
	f := NewWorkerIDFactory(dsn)

	for i := 0; i < 10; i++ {
		num1, err := f.WorkID()
		if err != nil {
			t.Error(err)
		}
		num2, err := f.WorkID()
		if err != nil {
			t.Error(err)
		}

		if num1 != num2 {
			t.Error(err)
		}
	}
}
