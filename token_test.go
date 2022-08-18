package uni_filter

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	exprs := []string{
		"Connections[].Find or (what=1 and aha=2 or aaa=3)",
		//"(what=1 and aha=2 or aaa=3)",
		//"what=1 and aha=2 or aaa=3)",
		"Connections[].Find (what=1 and aha=2 or aaa=3)",
		//"Connections[].Find or (or aaa=3)",
		//"Connections[].Find or (what =  1 and aha = 2 or aaa__like=3)",
		//"Connections[].Find or (what=1 and aha=2 or aaa=3)",
	}

	for i, expr := range exprs {
		ex, err := Parse(expr)
		if err != nil {
			t.Errorf("test case at index %d failed: %s\n", i, err)
		} else {
			fmt.Printf("### SUCCESS when Parse %s\n", expr)
			fmt.Println(ex.Strings())
		}
	}
}
