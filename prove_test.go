package prove

import "testing"

func TestParseArgs(t *testing.T) {
	p := NewProve()
	p.ParseArgs([]string{"-j", "100", "--exec", "foobar"})

	if p.Jobs != 100 {
		t.Errorf("want 100\ngot %d", p.Jobs)
	}
	if p.Exec != "foobar" {
		t.Errorf("want foobar\ngot %s", p.Exec)
	}
}
