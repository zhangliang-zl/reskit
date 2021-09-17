package srv

import (
	"fmt"
	"testing"
)

func TestClient_UUID(t *testing.T) {
	c := NewClient("10.25.177.169:5001")
	i, e := c.UUID()
	fmt.Println(i, e)
}
