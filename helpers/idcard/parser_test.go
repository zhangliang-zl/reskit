package idcard

import (
	"testing"
)

func TestParserValidate(t *testing.T) {

	male := []string{
		"110103200301019917",
		"110103200301017751",
		"53040120030101187X",
		"530401200301019513",
		"530401200301019337",
	}

	female := []string{
		"35052119740101080X",
		"350521197401019142",
		"350521197401012426",
		"35052119740101232x",
		"350521197401014747",
	}

	invalid := []string{
		"35052119740101003X",
		"350521197401019642",
		"350521197101012426",
		"35052119740101242x",
		"35021197401014747",
	}

	for _, v := range male {
		p := NewParser(v)
		if !p.Validate() {
			t.Errorf("male Validate() fail %s", v)
		}

		if p.Gender() != GenderMale {
			t.Errorf("male Gender() fail %s", v)
		}
	}

	for _, v := range female {
		p := NewParser(v)
		if !p.Validate() {
			t.Errorf("female Validate() fail %s", v)
		}

		if p.Gender() != GenderFemale {
			t.Errorf("female Gender() fail %s", v)
		}
	}

	for _, v := range invalid {
		p := NewParser(v)
		if p.Validate() {
			t.Errorf("invalid Validate() is pass ?  %s", v)
		}
	}
}
