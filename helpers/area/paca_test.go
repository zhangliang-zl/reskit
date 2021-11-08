package area

import (
	"testing"
)

func TestProvinces(t *testing.T) {
	v := Provinces()
	if len(v) == 0 {
		t.Error("province empty!")
	}
}

func TestCities(t *testing.T) {
	v := Cities("河北省")
	if len(v) == 0 {
		t.Error("city empty!")
	}
}

func TestAreas(t *testing.T) {
	v := Areas("河北省", "石家庄市")
	if len(v) == 0 {
		t.Error("area empty!")
	}
}

func TestStreets(t *testing.T) {
	v := Streets("河北省", "石家庄市", "正定县")
	if len(v) == 0 {
		t.Error("street empty!")
	}
}

func TestPCAS(t *testing.T) {
	if len(PACS()["河北省"]["石家庄市"]["正定县"]) == 0 {
		t.Error("pacs empty!")
	}
}
