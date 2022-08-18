package uni_filter

type OPFunc func(*Filter, interface{}) bool

type OPNewFunc func(string) (OP, error)

type OP interface {
	parse() error // used for parse op params
	check(v any, exists bool) bool
}

var ops = make(map[string]OPNewFunc, 30)

func register(name string, newFunc OPNewFunc) {
	ops[name] = newFunc // FIXME duplicate detect
}
