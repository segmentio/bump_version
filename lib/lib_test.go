package bump_version

import (
	"reflect"
	"testing"
)

func TestChangeVersion(t *testing.T) {
	testCases := []struct {
		in    string
		vtype VersionType
		out   string
	}{
		{"0.4", Major, "1.0"},
		{"0.4.0", Major, "1.0.0"},
		{"1.0", Major, "2.0"},
		{"1", Major, "2"},
		{"1.0.1", Minor, "1.1.0"},
	}
	for _, tt := range testCases {
		v, err := changeVersion(tt.vtype, tt.in)
		if err != nil {
			t.Fatal(err)
		}
		if v.String() != tt.out {
			t.Errorf("changeVersion(%s, %s): got %s, want %s", tt.vtype, tt.in, v.String(), tt.out)
		}
	}
}

func TestVersionString(t *testing.T) {
	typ := reflect.TypeOf(VERSION)
	if typ.String() != "string" {
		t.Errorf("expected VERSION to be a string, got %#v (type %#v)", VERSION, typ.String())
	}
}

func TestString(t *testing.T) {
	v := Version{0, 0, 1}
	if v.String() != "0.0.1" {
		t.Errorf("wrong version string reported")
	}
	v = Version{1, 0, 1}
	if v.String() != "1.0.1" {
		t.Errorf("wrong version string reported")
	}
	v = Version{1, -1, -1}
	if want := "1"; v.String() != want {
		t.Errorf("wrong version string reported: got %q want %q", v.String(), want)
	}
}

var lessTests = []struct {
	i, j string
	want bool
}{
	{"1", "2", true},
	{"1.1", "2", true},
	{"1.3.7", "2", true},
	{"1.3.7", "0.1", false},
	{"1.3.7", "1.3.8", true},
	{"1.3.7", "1.3.6", false},
	{"1.3.7", "1.3.7", false},
	{"1.3.7", "1", false},
}

func TestLess(t *testing.T) {
	for _, tt := range lessTests {
		i, _ := Parse(tt.i)
		j, _ := Parse(tt.j)
		got := Less(i, j)
		if got != tt.want {
			t.Errorf("Less(%q, %q): got %t, want %t", i, j, got, tt.want)
		}
	}
}
