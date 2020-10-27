package mod

import "testing"

func TestFromGoMod(t *testing.T) {
	v, err := FromGoMod("../go.mod")
	if err != nil {
		t.Error(err)
	}
	if v != "github.com/gotrix/cli" {
		t.Errorf("expected %s, got %s", "github.com/gotrix/cli", v)
	}
}
