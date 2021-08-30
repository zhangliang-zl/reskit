package area

import "testing"

func TestHkMoTwProvinces(t *testing.T) {
	v := HkMoTwProvinces()
	if len(v) == 0 {
		t.Error("province empty!")
	}
}

func TestHkMoTwCities(t *testing.T) {
	v := HkMoTwCities("香港特别行政区")
	if len(v) == 0 {
		t.Error("city empty!")
	}
}

func TestHkMoTwArea(t *testing.T) {
	v := HkMoTwAreas("香港特别行政区", "九龙")
	if len(v) == 0 {
		t.Error("area empty!")
	}
}

func TestHkMoTw(t *testing.T) {
	v := HkMoTw()
	if len(v) == 0 {
		t.Error("HkMoTw empty!")
	}
}
