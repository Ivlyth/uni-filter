package uni_filter

import "testing"

func TestOPEnds(t *testing.T) {
	cases := []struct {
		opValue string
		v       string
		ret     bool
	}{
		{
			"at",
			"what",
			true,
		},
		{
			"ha",
			"ahaha",
			true,
		},
		{
			"si",
			"this",
			false,
		},
	}

	for i, tc := range cases {
		op, _ := NewOPEnds(tc.opValue)
		ret := op.check(tc.v, true)
		if ret != tc.ret {
			t.Errorf("test case at index %d failed\n", i)
		}
	}
}

func TestOPIEnds(t *testing.T) {
	cases := []struct {
		opValue string
		v       string
		ret     bool
	}{
		{
			"AT",
			"what",
			true,
		},
		{
			"Ha",
			"ahaha",
			true,
		},
		{
			"Si",
			"this",
			false,
		},
	}

	for i, tc := range cases {
		op, _ := NewOPIEnds(tc.opValue)
		ret := op.check(tc.v, true)
		if ret != tc.ret {
			t.Errorf("test case at index %d failed\n", i)
		}
	}
}
