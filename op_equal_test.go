package uni_filter

import "testing"

func TestOPEqual(t *testing.T) {
	cases := []struct {
		opValue string
		v       any
		ret     bool
	}{
		{
			"what",
			"what",
			true,
		},
		{
			"ahaha",
			"hahaa",
			false,
		},
		{
			"123456",
			123456,
			true,
		},
		{
			"123456.78",
			123456.78,
			true,
		},
		{
			"true",
			true,
			true,
		},
		{
			"false",
			false,
			true,
		},
	}

	for i, tc := range cases {
		op, _ := NewOPEqual(tc.opValue)
		ret := op.check(tc.v, true)
		if ret != tc.ret {
			t.Errorf("test case at index %d failed\n", i)
		}
	}
}
