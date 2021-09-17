package srv

import (
	"fmt"
	"testing"
)

func TestClient_UUID(t *testing.T) {
	c := NewClient("10.25.177.169:5001")

	i := 0
	for {
		item, err := c.UUID()
		fmt.Println(item, err)

		if i > 10 {
			break
		}
		i++
	}

}
