package flexconfig

import (
	"testing"
)

func Test_dts(t *testing.T) {
	s := dotsToSlash("")
	if len(s) > 0 {
		t.Errorf("Unexpected value: %s", s)
	}

	s = dotsToSlash("a.b.c")
	if s != "/a/b/c" {
		t.Errorf("Unexpected value: %s", s)
	}

	s = dotsToSlash(".a.b.c.")
	if s != "/a/b/c/" {
		t.Errorf("Unexpected value: %s", s)
	}
}

func Test_std(t *testing.T) {
	s := slashToDots("")
	if len(s) > 0 {
		t.Errorf("Unexpected value: %s", s)
	}

	s = slashToDots("a/b/c")
	if s != "a.b.c" {
		t.Errorf("Unexpected value: %s", s)
	}

	s = slashToDots("/a/b/c/")
	if s != ".a.b.c." {
		t.Errorf("Unexpected value: %s", s)
	}
}
