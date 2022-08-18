package uni_filter

import "errors"

type OPEqual struct {
	s string
}

func (op *OPEqual) Name() string {
	return "equal"
}

func (op *OPEqual) parse() error {
	if op.s == "" {
		return errors.New("op param for equal can not be empty")
	}
	return nil
}

func (op *OPEqual) check(v any, exists bool) bool {
	if !exists {
		return false
	}
	if op.s == convert2string(v) { // FIXME
		return true
	}
	return false
}

func NewOPEqual(s string) (OP, error) {
	op := &OPEqual{s: s}
	return op, op.parse()
}

func init() {
	register("equal", NewOPEqual)
	register("eq", NewOPEqual)
}
