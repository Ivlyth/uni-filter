package uni_filter

import "fmt"

func Contains(arr []string, t string) bool {
	for _, item := range arr {
		if item == t {
			return true
		}
	}
	return false
}

func All(flags []bool) bool {
	for _, flag := range flags {
		if !flag {
			return false
		}
	}
	return true
}

func Any(flags []bool) bool {
	for _, flag := range flags {
		if flag {
			return true
		}
	}
	return false
}

func convert2string(v any) string {
	return fmt.Sprintf("%v", v)
}

var _ = All(nil)
var _ = Any(nil)
