package uni_filter

import (
	"errors"
	"fmt"
	"strings"
)

type OPEnds struct {
	s          string
	ignoreCase bool

	is string
}

func (op *OPEnds) Name() string {
	if op.ignoreCase {
		return "iends"
	}
	return "ends"
}

func (op *OPEnds) parse() error {
	if op.s == "" {
		return errors.New(fmt.Sprintf("op param for %s can not be empty", op.Name()))
	}
	if op.ignoreCase {
		op.is = strings.ToLower(op.s)
	}
	return nil
}

func (op *OPEnds) check(v any, exists bool) bool {
	if !exists {
		return false
	}
	s := convert2string(v)
	if op.ignoreCase {
		s = strings.ToLower(s)
		return strings.HasSuffix(s, op.is)
	} else {
		return strings.HasSuffix(s, op.s)
	}
}

func NewOPEnds(s string) (OP, error) {
	op := &OPEnds{s: s}
	return op, op.parse()
}

func NewOPIEnds(s string) (OP, error) {
	op := &OPEnds{s: s, ignoreCase: true}
	return op, op.parse()
}

func init() {
	register("ends", NewOPEnds)
	register("iends", NewOPIEnds)
}
