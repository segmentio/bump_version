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
